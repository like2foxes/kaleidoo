package router

import (
	"context"
	"html/template"
	"io"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/like2foxes/nirlir/pg"
)

func Start(jwtSecret string) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://nirlir:nirlir@localhost:5432/nirlir?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	err = conn.Ping(ctx)
	if err != nil {
		panic(err)
	}

	queries := pg.New(conn)
	ts := NewTemplates()
	ts.AddShared("index").AddFiles("index", "templates/index.html")
	ts.AddShared("login").AddFiles("login", "templates/users/login.html")
	ts.AddShared("register").AddFiles("register", "templates/users/register.html")
	ts.AddShared("customers").AddFiles("customers", "templates/customers/customers.html", "templates/customers/customer_row.html")
	ts.AddShared("orders").AddFiles("orders", "templates/orders/orders.html", "templates/orders/order_row.html")
	ts.AddFiles("addCustomer", "templates/customers/customer_modal.html", "templates/customers/form/customer_add.html")
	ts.AddFiles("postCustomer", "templates/customers/customer_row.html", "templates/shared/delete_svg.html", "templates/shared/edit_svg.html")
	ts.AddFiles("editCustomer", "templates/customers/customer_modal.html", "templates/customers/form/customer_edit.html")
	ts.AddFiles("putCustomer", "templates/customers/customer_row.html", "templates/shared/delete_svg.html", "templates/shared/edit_svg.html")
	ts.AddFiles("addOrder", "templates/orders/order_modal.html", "templates/orders/form/orders_add.html")
	ts.AddFiles("postOrder", "templates/orders/order_row.html", "templates/shared/delete_svg.html", "templates/shared/edit_svg.html")
	ts.AddFiles("editOrder", "templates/orders/order_modal.html", "templates/orders/form/orders_edit.html")
	ts.AddFiles("putOrder", "templates/orders/order_row.html", "templates/shared/delete_svg.html", "templates/shared/edit_svg.html")
	ts.AddFiles("filterOrders", "templates/orders/orders_table.html", "templates/orders/order_row.html", "templates/shared/delete_svg.html", "templates/shared/edit_svg.html")
	ts.AddShared("forecast").AddFiles("forecast", "templates/forecast/forecast.html", "templates/forecast/forecast_row.html")
	ts.AddFiles("logo", "templates/customers/detected_name.html")

	t := &Template{
		templates: ts,
	}

	e := echo.New()
	e.Renderer = t

	e.Static("/", "static")
	e.File("/favicon.ico", "static/assets/favicon.ico")

	e.GET("/", getIndex)
	e.GET("/login", getLogin)
	e.POST("/login", postLogin(queries, ctx, jwtSecret))
	e.GET("/register", getRegister)
	e.POST("/register", postRegister(queries, ctx))

	auth := e.Group("")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "cookie:token",
	}
	auth.Use(echojwt.WithConfig(config))
	auth.GET("", autherize)

	auth.GET("/customers", getCustomers(queries, ctx))
	auth.POST("/customers", postCustomer(queries, ctx))
	auth.GET("/customers/add", getAddCustomer)
	auth.GET("/customers/:id", getEditCustomer(queries, ctx))
	auth.DELETE("/customers/:id", deleteCustomer(queries, ctx))
	auth.PUT("/customers/:id", putCustomer(queries, ctx))

	auth.GET("/orders", getOrders(queries, ctx))
	auth.POST("/orders", postOrder(queries, ctx))
	auth.GET("/orders/add", getAddOrder(queries, ctx))
	auth.POST("/orders/filter", filterOrders(queries, ctx))
	auth.GET("/orders/:id", getEditOrder(queries, ctx))
	auth.DELETE("/orders/:id", deleteOrder(queries, ctx))
	auth.PUT("/orders/:id", putOrder(queries, ctx))

	auth.GET("/forecast", getForecast(queries, ctx))

	auth.POST("/logo", postLogo)

	e.Logger.Fatal(e.Start(":3000"))
}

func getIndex(c echo.Context) error {
	return c.Render(200, "index", nil)
}

func getForecast(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		forecast, err := q.ListForecast(ctx)
		if err != nil {
			log.Println("err", err)
			return c.String(500, "internal error")
		}
		log.Printf("Forecast: %v", forecast)
		return c.Render(200, "forecast", forecast)
	}
}

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	n := t.templates[name]
	log.Printf("Rendering template %s: %v", n.Name(), n.DefinedTemplates())
	if n == nil {
		return echo.ErrNotFound
	}
	if strings.Contains(n.DefinedTemplates(), "layout") {
		return n.ExecuteTemplate(w, "layout", data)
	}
	if strings.Contains(n.DefinedTemplates(), "modal") {
		log.Printf("Executing modal template %s", n.Name())
		return n.ExecuteTemplate(w, "modal", data)
	}
	if strings.Contains(n.DefinedTemplates(), "orders_table") {
		log.Printf("Executing orders_table template %s", n.Name())
		return n.ExecuteTemplate(w, "orders_table", data)
	}
	if strings.Contains(n.DefinedTemplates(), "row") {
		return n.ExecuteTemplate(w, "row", data)
	}
	log.Printf("Executing template %s", n.Name())
	return n.Execute(w, data)

}

type Templates map[string]*template.Template

func NewTemplates() Templates {
	return make(Templates)
}

func (t Templates) AddGlob(name, path string) Templates {
	if t[name] == nil {
		tpl := template.New(name)
		_, err := tpl.ParseGlob(path)
		if err != nil {
			log.Panicf("Error parsing glob for template %s with path %s: %v", name, path, err)
		}
		t[name] = tpl
		return t
	}
	_, err := t[name].ParseGlob(path)
	if err != nil {
		log.Panicf("Error parsing glob for template %s with path %s: %v", name, path, err)
	}

	return t
}

func (t Templates) AddFiles(name string, files ...string) Templates {
	if t[name] == nil {

		tpl := template.New(name)
		_, err := tpl.ParseFiles(files...)
		if err != nil {
			log.Panicf("Error parsing template %s from files %v: %v", name, files, err)
		}
		t[name] = tpl
		return t
	}
	_, err := t[name].ParseFiles(files...)
	if err != nil {
		log.Panicf("Error parsing template %s from files %v: %v", name, files, err)
	}

	return t
}

func (t Templates) AddShared(name string) Templates {
	t.AddFiles(name, "templates/shared/layout.html")
	return t.AddGlob(name, "templates/shared/*.html")
}

func autherize(c echo.Context) error {
	return c.Render(200, "index", nil)
}
