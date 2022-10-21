package handlers

import (
	"net/http"

	"github.com/chau-t-tran/ws-relay/utils"
	"github.com/labstack/echo/v4"
)

type (
	RoomRegResponse struct {
		SessionKey string `json:"sessionKey"`
	}
)

func APIHandler(c echo.Context) error {
	response := RoomRegResponse{
		SessionKey: utils.RandomKey(),
	}
	return c.JSON(http.StatusCreated, response)
}
