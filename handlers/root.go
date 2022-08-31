package handlers

import (
	"fmt"
	"net/http"

	"github.com/chau-t-tran/ws-relay/utils"
	"github.com/labstack/echo/v4"
)

func RootHandler(c echo.Context) error {
	key := utils.RandomKey()
	fmt.Println(key)
	c.Redirect(http.StatusFound, fmt.Sprintf("/%s", key))
	return nil
}
