package handlers

import (
	"fmt"
	"net/http"

	"github.com/chau-t-tran/ws-relay/utils"
	"github.com/labstack/echo/v4"
)

func RootHandler(c echo.Context) error {
	sessionKey := utils.RandomKey()
	c.Redirect(http.StatusFound, fmt.Sprintf("/%s", sessionKey))
	return nil
}

func SessionHandler(c echo.Context) error {
	sessionKey := c.Param("sessionKey")
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"sessionKey": sessionKey,
		"host":       c.Request().Host,
	})
}
