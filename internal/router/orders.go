package router

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/like2foxes/nirlir/pg"
)

func getOrders(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		orders, err := q.GetWorkOrders(ctx)
		if err != nil {
			log.Println("orders: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		if len(orders) < 1 {
			a := []struct {
				ID       int32
				Customer string
				Type     string
				Status   string
			}{
				{
					ID:       1,
					Customer: "koko",
					Type:     "installation",
					Status:   "open",
				},
			}
			log.Println(a)
			return c.Render(http.StatusOK, "orders", a)
		}
		return c.Render(http.StatusOK, "orders", orders)
	}
}

func getAddOrder(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		customers, err := q.ListCustomers(ctx)
		if err != nil {
			log.Println("addOrders: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		return c.Render(http.StatusOK, "addOrder", OrdersData{Customers: customers})
	}
}

type OrdersData struct {
	Customers []pg.ListCustomersRow
	Order     pg.GetWorkOrderRow
}

func postOrder(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		log.Println(c.FormParams())
		log.Println(c.FormValue("customer"))
		customer, err := strconv.Atoi(c.FormValue("customer"))
		if err != nil {
			log.Println("postOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		t := c.FormValue("type")
		status := c.FormValue("status")

		neworder, err := q.CreateWorkOrder(ctx, pg.CreateWorkOrderParams{
			CustomerID: int32(customer),
			Type:       t,
			Status:     status,
		})
		if err != nil {
			log.Println("postOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		cust, err := q.GetCustomer(ctx, int32(customer))
		if err != nil {
			log.Println("postOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		row := struct {
			ID       int32
			Customer string
			Type     string
			Status   string
		}{
			ID:       neworder.ID,
			Customer: cust.Name,
			Type:     t,
			Status:   status,
		}

		return c.Render(http.StatusOK, "postOrder", row)
	}
}

func deleteOrder(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println("deleteOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		err = q.DeleteWorkOrder(ctx, int32(id))
		if err != nil {
			log.Println("deleteOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.NoContent(http.StatusOK)
	}
}

func getEditOrder(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println("getEditOrder: id", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		order, err := q.GetWorkOrder(ctx, int32(id))
		if err != nil {
			log.Println("getEditOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		customers, err := q.ListCustomers(ctx)
		if err != nil {
			log.Println("getEditOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.Render(http.StatusOK, "editOrder", OrdersData{Customers: customers, Order: order})
	}
}

func putOrder(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println("putOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		customer, err := strconv.Atoi(c.FormValue("customer"))
		if err != nil {
			log.Println("putOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		t := c.FormValue("type")
		status := c.FormValue("status")

		order, err := q.UpdateWorkOrder(ctx, pg.UpdateWorkOrderParams{
			ID:         int32(id),
			CustomerID: int32(customer),
			Type:       t,
			Status:     status,
		})
		if err != nil {
			log.Println("putOrder: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		row := struct {
			ID       int32
			Customer string
			Type     string
			Status   string
		}{
			ID:       order.ID,
			Customer: c.FormValue("customer"),
			Type:     order.Type,
			Status:   order.Status,
		}

		return c.Render(http.StatusOK, "putOrder", row)
	}
}

func filterOrders(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		params, _ := c.FormParams()
		log.Println("all :", params)
		log.Println("param: ", c.FormValue("cancelled"))
		woparams := pg.GetFilteredWorkOrdersParams{}
		if c.FormValue("cancelled") != "" {
			woparams.Status = "cancelled"
		}
		if c.FormValue("open") != "" {
			woparams.Status_2 = "open"
		}
		if c.FormValue("closed") != "" {
			woparams.Status_3 = "close"
		}
		if c.FormValue("onsite") != "" {
			woparams.Status_4 = "onsite"
		}
		if c.FormValue("travel") != "" {
			woparams.Status_5 = "travel"
		}
		if c.FormValue("assigned") != "" {
			woparams.Status_6 = "assigned"
		}
		log.Println("woparams: ", woparams)

		orders, err := q.GetFilteredWorkOrders(ctx, woparams)
		log.Println(orders)

		if err != nil {
			log.Println("filterOrders: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.Render(http.StatusOK, "filterOrders", orders)
	}
}
