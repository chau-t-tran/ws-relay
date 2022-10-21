package handlers

import (
	"net/http"

	"github.com/chau-t-tran/ws-relay/utils"
	"github.com/labstack/echo/v4"
)

func RootHandler(c echo.Context) error {
	sessionKey := utils.RandomKey()
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"sessionKey": sessionKey,
		"host":       c.Request().Host,
	})
}
