package utils

import (
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

/*
func (sm *SessionManager) Broadcast(message string, conn websocket.Conn) {
	for _, c := range sm.Sessions {
	}
}
*/
