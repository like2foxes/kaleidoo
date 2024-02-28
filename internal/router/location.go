package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getLocation(c echo.Context) error {
	return c.Render(http.StatusOK, "location", nil)
}
