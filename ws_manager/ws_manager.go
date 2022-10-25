package ws_manager

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
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
	sessions         map[string][]*websocket.Conn
	lastUsed         map[string]time.Time
	currentTime      time.Time
	sessionManagerMu sync.Mutex

	cronScheduler *gocron.Scheduler
	maxAliveTime  time.Duration
}

func CreateSessionManager(sessionKeys []string) SessionManager {
	sm := SessionManager{
		sessions: map[string][]*websocket.Conn{},
		lastUsed: map[string]time.Time{},
	}
	for _, key := range sessionKeys {
		sm.sessions[key] = []*websocket.Conn{}
		sm.lastUsed[key] = time.Now()
	}

	sm.maxAliveTime = 24 * time.Hour
	sm.currentTime = time.Now()
	sm.cronScheduler = gocron.NewScheduler(time.UTC)
	_, _ = sm.cronScheduler.
		Every(1).
		Day().
		Do(sm.GarbageCollectDaily)
	sm.cronScheduler.StartAsync()
	return sm
}

func (sm *SessionManager) GarbageCollectDaily() {
	sm.sessionManagerMu.Lock()
	defer sm.sessionManagerMu.Unlock()
	sm.currentTime = sm.currentTime.Add(24 * time.Hour)
	for key, _ := range sm.sessions {
		lastTime, err := sm.GetLastUsedTime(key)
		if err != nil {
			panic(err)
		}
		aliveTime := sm.currentTime.Sub(lastTime)
		if aliveTime > sm.maxAliveTime {
			// in-loop deletion safe in go
			delete(sm.sessions, key)
			delete(sm.lastUsed, key)
		}
	}
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
	sm.sessionManagerMu.Lock()
	defer sm.sessionManagerMu.Unlock()
	_, err := sm.GetSession(sessionKey)
	if err == nil {
		return errors.New(
			fmt.Sprintf("Session %s already exists", sessionKey),
		)
	}
	sm.sessions[sessionKey] = []*websocket.Conn{}

	_, err = sm.GetLastUsedTime(sessionKey)
	if err == nil {
		return errors.New(
			fmt.Sprintf("Session %s last used time already exists", sessionKey),
		)
	}
	sm.lastUsed[sessionKey] = time.Now()
	return nil
}

func (sm *SessionManager) AddConnection(sessionKey string, ws *websocket.Conn) error {
	sm.sessionManagerMu.Lock()
	defer sm.sessionManagerMu.Unlock()
	if session, ok := sm.sessions[sessionKey]; ok {
		sm.sessions[sessionKey] = append(session, ws)
		return nil
	} else {
		return errors.New(
			fmt.Sprintf("Session %s does not exist", sessionKey),
		)
	}
}

func (sm *SessionManager) Broadcast(sessionKey string, senderAddr string, message []byte) error {
	sm.sessionManagerMu.Lock()
	defer sm.sessionManagerMu.Unlock()
	session, err := sm.GetSession(sessionKey)
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

	sm.lastUsed[sessionKey] = time.Now()
	return nil
}
