package handlers

import (
	"net/http"

	"github.com/chau-t-tran/ws-relay/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

func RootHandler(c echo.Context) error {
	sessionKey := utils.RandomKey()
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"sessionKey": sessionKey,
		"host":       c.Request().Host,
	})
}
