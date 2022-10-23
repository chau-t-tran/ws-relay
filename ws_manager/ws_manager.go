package ws_manager

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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
	lastUsed map[string]time.Time
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

func (sm *SessionManager) GetSession(sessionKey string) ([]*websocket.Conn, error) {
	session, ok := sm.sessions[sessionKey]
	if !ok {
		return session, errors.New(
			fmt.Sprintf("Session %s not found", sessionKey),
		)
	}
	return session, nil
}

func (sm *SessionManager) GetLastUsedTime(sessionKey string) (time.Time, error) {
	time, ok := sm.lastUsed[sessionKey]
	if !ok {
		return time, errors.New(
			fmt.Sprintf("Session %s last used time not found", sessionKey),
		)
	}
	return time, nil
}

func (sm *SessionManager) RegisterSession(sessionKey string) error {
	_, err := sm.GetSession(sessionKey)
	if err == nil {
		return errors.New(
			fmt.Sprintf("Session %s already exists", sessionKey),
		)
	}

	sm.sessions[sessionKey] = []*websocket.Conn{}
	return nil
}

func (sm *SessionManager) AddConnection(sessionKey string, ws *websocket.Conn) error {
	if session, ok := sm.sessions[sessionKey]; ok {
		sm.sessions[sessionKey] = append(session, ws)
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
