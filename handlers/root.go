package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")
}
