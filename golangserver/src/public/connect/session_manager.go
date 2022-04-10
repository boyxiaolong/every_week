package connect

import (
	"sync"
	"public/command"
	"public/common"
)

var GSessionManager* SessionManager

func init() {
	GSessionManager = NewSessionManager()
	command.GCommand.RegCommand("break", DoBreak, "break")
}

func DoBreak(str *common.StringParse) (err error) {
	//GSessionManager.Break()
	return nil
}

type SessionManager struct {
	mutex         sync.RWMutex
	SessionMap map[uint64]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		SessionMap: make(map[uint64]*Session),
	}
}

func (manager *SessionManager) Len() int {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	return len(manager.SessionMap)
}

func (manager *SessionManager) Get(key uint64) *Session {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	session, _ := manager.SessionMap[key]
	return session
}

func (manager *SessionManager) Put(key uint64, session *Session) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	if session, exists := manager.SessionMap[key]; exists {
		manager.remove(key, session)
	}

	manager.SessionMap[key] = session
}

func (manager *SessionManager) remove(key uint64, session *Session) {
	delete(manager.SessionMap, key)
}

func (manager *SessionManager) Remove(key uint64) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	session, exists := manager.SessionMap[key]
	if exists {
		manager.remove(key, session)
	}

	return exists
}

func (manager *SessionManager) Clear()  {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.SessionMap = make(map[uint64]*Session)

}
