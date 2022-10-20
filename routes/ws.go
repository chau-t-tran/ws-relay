package routes

import (
	"fmt"

	"github.com/chau-t-tran/ws-relay/ws_manager"
	"github.com/labstack/echo/v4"
)

func SetupWSRoutes(e *echo.Echo) {
	sm := ws_manager.CreateSessionManager()
	e.GET("/:sessionKey", sm.EchoHandler)
	fmt.Println("WS routes setup!")
}
