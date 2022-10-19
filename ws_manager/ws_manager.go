package ws_manager

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type SessionManager struct {
	sessions map[string][]*websocket.Conn
}

func CreateSessionManager() SessionManager {
	sm := SessionManager{
		sessions: map[string][]*websocket.Conn{},
	}
	return sm
}

func (s *SessionManager) GetConnections(sessionKey string) []*websocket.Conn {
	return s.sessions[sessionKey]
}

func (s *SessionManager) AddConnection(sessionKey string, ws *websocket.Conn) {
	if session, ok := s.sessions[sessionKey]; ok {
		s.sessions[sessionKey] = append(session, ws)
	} else {
		s.sessions[sessionKey] = []*websocket.Conn{ws}
	}
}

func (sm *SessionManager) Broadcast(key string, senderAddr string, message string) {
	for _, c := range sm.GetConnections(key) {
		receiverAddr := c.RemoteAddr().String()
		if receiverAddr == senderAddr {
			continue
		}
		err := c.WriteMessage(1, []byte(message))
		if err != nil {
			continue
		}
	}
}

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

		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		err = conn.WriteMessage(mt, []byte(message))
		if err != nil {
			return err
		}
	}
	return nil
}
