package router

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	logorec "github.com/like2foxes/nirlir/logo_rec"
	"github.com/like2foxes/nirlir/pg"
)

func getCustomers(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		customers, err := q.ListCustomers(ctx)

		if err != nil {
			log.Println("err", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		log.Printf("customers %v\n", customers)

		return c.Render(http.StatusOK, "customers", customers)
	}
}

func getEditCustomer(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, "invalid id")
		}

		customer, err := q.GetCustomer(ctx, int32(id))
		if err != nil {
			log.Println("customer error: ", err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.Render(http.StatusOK, "editCustomer", customer)
	}
}

func getAddCustomer(c echo.Context) error {
	log.Println("add customer")
	return c.Render(http.StatusOK, "addCustomer", nil)
}

func postCustomer(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		logoFile, err := c.FormFile("logo")

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		src, err := logoFile.Open()
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		defer src.Close()

		path := "static/images/" + name + "." + logoFile.Filename
		dst, err := os.Create(path)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		p := pgtype.Text{String: strings.TrimPrefix(path, "static/"), Valid: true}
		newCustomer, err := q.CreateCustomer(ctx, pg.CreateCustomerParams{Name: name, Logo: p})
		if err != nil {
			log.Println(err)
			os.Remove(path)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.Render(http.StatusOK, "postCustomer", newCustomer)
	}
}

func deleteCustomer(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, "invalid id")
		}

		err = q.DeleteCustomerWorkOrders(ctx, int32(id))
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		logo, err := q.GetCustomerLogo(ctx, int32(id))
		if err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, "no user")
		}

		path := fmt.Sprintf("static/%s", logo.String)
		if err = os.Remove(path); err != nil {
			log.Println(err)
		}

		err = q.DeleteCustomer(ctx, int32(id))
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.NoContent(http.StatusOK)
	}
}

func putCustomer(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			return c.String(http.StatusBadRequest, "invalid id")
		}

		name := c.FormValue("name")
		logoFile, err := c.FormFile("logo")
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		src, err := logoFile.Open()
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		defer src.Close()

		path := "static/images/" + name + "." + logoFile.Filename
		dst, err := os.Create(path)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		p := pgtype.Text{String: strings.TrimPrefix(path, "static/"), Valid: true}
		cust, err := q.UpdateCustomer(ctx, pg.UpdateCustomerParams{ID: int32(id), Name: name, Logo: p})
		if err != nil {
			log.Println(err)
			os.Remove(path)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.Render(http.StatusOK, "putCustomer", cust)
	}
}

func postLogo(c echo.Context) error {
	f, err := c.FormFile("logo")
	if err != nil {
		log.Println("form file ", err)
		return c.String(http.StatusInternalServerError, "internal error")
	}

	src, err := f.Open()
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "internal error")
	}
	defer src.Close()

	results, err := logorec.RecognizeLogo(bufio.NewReader(src))

	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "internal error")
	}

	return c.Render(http.StatusOK, "logo", results[1])
}
