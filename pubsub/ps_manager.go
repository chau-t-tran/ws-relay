package utils

import (
	"log"

	"github.com/gorilla/websocket"
)

type Session struct {
	Connections []*websocket.Conn
}

type SessionManager struct {
	Sessions map[string]Session
}

func (s *SessionManager) GetConnections(sessionKey string) []*websocket.Conn {
	return s.Sessions[sessionKey].Connections
}

func (s *SessionManager) AddConnection(sessionKey string, ws *websocket.Conn) {
	if session, ok := s.Sessions[sessionKey]; ok {
		session.Connections = append(session.Connections, ws)
	} else {
		s.Sessions[sessionKey] = Session{Connections: []*websocket.Conn{ws}}
	}
}

func (sm *SessionManager) Broadcast(key string, conn *websocket.Conn, message string) {
	for _, c := range sm.GetConnections(key) {
		if c == conn {
			continue
		}
		err := c.WriteMessage(1, []byte(message))
		if err != nil {
			log.Println("Error sending message")
			continue
		}
	}
}
