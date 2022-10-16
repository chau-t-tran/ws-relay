package utils

import (
	"github.com/gorilla/websocket"
)

type SessionManager struct {
	sessions map[string][]*websocket.Conn
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
