package ws_manager

import (
	"log"
	"net/http"

	"github.com/chau-t-tran/ws-to-me/utils"
	"github.com/labstack/echo/v4"
)

type (
	SessionRegResponse struct {
		SessionKey string `json:"sessionKey"`
	}
)

// two different ways to register a session: via landing page or via API

func (sm *SessionManager) RegisterHandler(c echo.Context) error {
	sessionKey := utils.RandomKey()
	sm.RegisterSession(sessionKey)
	response := SessionRegResponse{
		SessionKey: sessionKey,
	}
	return c.JSON(http.StatusCreated, response)
}

func (sm *SessionManager) RootHandler(c echo.Context) error {
	sessionKey := utils.RandomKey()
	sm.RegisterSession(sessionKey)
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"sessionKey": sessionKey,
		"host":       c.Request().Host,
	})
}

// main websocket handler

func (sm *SessionManager) EchoHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	sessionKey := c.Param("sessionKey")
	if err != nil {
		log.Println("upgrade error:", err)
		return err
	}
	defer conn.Close()
	sm.AddConnection(sessionKey, conn)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		err = sm.Broadcast(sessionKey, conn.RemoteAddr().String(), message)
		if err != nil {
			return err
		}
	}
	return nil
}
