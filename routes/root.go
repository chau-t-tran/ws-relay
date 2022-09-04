package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/chau-t-tran/ws-relay/handlers"
)

func SetupRootRoutes(e *echo.Echo) {
	e.GET("/", handlers.RootHandler)
	e.GET("/:sessionKey", handlers.SessionHandler)
	fmt.Println("Index routes setup!")
}
