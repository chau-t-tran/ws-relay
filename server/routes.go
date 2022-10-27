package server

import (
	"fmt"

	"github.com/chau-t-tran/ws-to-me/ws_manager"
	"github.com/labstack/echo/v4"
)

func SetupWSRoutes(e *echo.Echo) {
	sm := ws_manager.CreateSessionManager([]string{})
	e.GET("/", sm.RootHandler)
	e.GET("/:sessionKey", sm.EchoHandler)
	e.POST("/register", sm.RegisterHandler)
	fmt.Println("WS routes setup!")
}
