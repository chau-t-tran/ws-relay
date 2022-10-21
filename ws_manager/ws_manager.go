package ws_manager

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func CheckOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOrigin,
}

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

func (sm *SessionManager) Broadcast(key string, senderAddr string, message []byte) error {
	for _, c := range sm.GetConnections(key) {
		receiverAddr := c.RemoteAddr().String()
		if receiverAddr == senderAddr {
			continue
		}
		err := c.WriteMessage(1, message)
		if err != nil {
			return err
		}
	}
	return nil
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
	log.Println("Added to", sessionKey)
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
