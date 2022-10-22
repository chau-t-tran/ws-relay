package ws_manager

import (
	"errors"
	"fmt"
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

func CreateSessionManager(sessionKeys []string) SessionManager {
	sm := SessionManager{
		sessions: map[string][]*websocket.Conn{},
	}
	for _, key := range sessionKeys {
		sm.sessions[key] = []*websocket.Conn{}
	}
	return sm
}

func (s *SessionManager) GetSession(sessionKey string) ([]*websocket.Conn, error) {
	session, ok := s.sessions[sessionKey]
	if !ok {
		return session, errors.New(
			fmt.Sprintf("Session %s not found", sessionKey),
		)
	}
	return session, nil
}

func (s *SessionManager) RegisterSession(sessionKey string) error {
	_, err := s.GetSession(sessionKey)
	if err == nil {
		return errors.New(
			fmt.Sprintf("Session %s already exists", sessionKey),
		)
	}

	s.sessions[sessionKey] = []*websocket.Conn{}
	return nil
}

func (s *SessionManager) AddConnection(sessionKey string, ws *websocket.Conn) error {
	if session, ok := s.sessions[sessionKey]; ok {
		s.sessions[sessionKey] = append(session, ws)
		return nil
	} else {
		return errors.New(
			fmt.Sprintf("Session %s does not exist", sessionKey),
		)
	}
}

func (sm *SessionManager) Broadcast(key string, senderAddr string, message []byte) error {
	session, err := sm.GetSession(key)
	if err != nil {
		return err
	}

	for _, c := range session {
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
