package session

import (
	"sync"
)

type Session interface {
	GetId() int64
	OnOpen()
	OnReceive(msgByte []byte)
	OnClose()
	OnError(err error)
	Send(msgByte []byte) error
	Close() error
}

func newSessionManager() (sm *SessionManager) {
	sm = &SessionManager{}
	sm.sessions = make(map[int64]Session)
	sm.rwm = &sync.RWMutex{}
	return
}

type SessionManager struct {
	rwm      *sync.RWMutex
	sessions map[int64]Session
}

func (sm *SessionManager) PutSession(s Session) bool {
	sm.rwm.Lock()
	defer sm.rwm.Unlock()
	_, ok := sm.sessions[s.GetId()]
	if ok {
		return false
	}

	sm.sessions[s.GetId()] = s
	return true
}

func (sm *SessionManager) RemoveSession(s Session) {
	sm.rwm.Lock()
	defer sm.rwm.Unlock()
	_, ok := sm.sessions[s.GetId()]
	if !ok {
		return
	}
	delete(sm.sessions, s.GetId())
}

func (sm *SessionManager) SessionById(id int64) Session {
	sm.rwm.Lock()
	defer sm.rwm.Unlock()
	s, _ := sm.sessions[id]
	return s
}
