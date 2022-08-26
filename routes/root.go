package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRootRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})
	fmt.Println("Index routes setup!")
}
