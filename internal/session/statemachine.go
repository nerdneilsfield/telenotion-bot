package session

import "sync"

type StateMachine struct {
	mu       sync.RWMutex
	sessions map[int64]*Session
}

func NewStateMachine() *StateMachine {
	return &StateMachine{sessions: make(map[int64]*Session)}
}

func (sm *StateMachine) StartSession(chatID int64) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.sessions[chatID]; exists {
		return false
	}

	sm.sessions[chatID] = &Session{ChatID: chatID, Blocks: []Block{}}
	return true
}

func (sm *StateMachine) ClearSession(chatID int64) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[chatID]
	if !exists {
		return false
	}

	session.Blocks = []Block{}
	return true
}

func (sm *StateMachine) DiscardSession(chatID int64) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.sessions[chatID]; !exists {
		return false
	}

	delete(sm.sessions, chatID)
	return true
}

func (sm *StateMachine) EndSession(chatID int64) (*Session, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[chatID]
	if !exists {
		return nil, false
	}

	delete(sm.sessions, chatID)
	return session, true
}

func (sm *StateMachine) GetSession(chatID int64) *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.sessions[chatID]
}

func (sm *StateMachine) IsActive(chatID int64) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	_, exists := sm.sessions[chatID]
	return exists
}

func (sm *StateMachine) AppendBlock(chatID int64, block Block) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[chatID]
	if !exists {
		return
	}

	session.Blocks = append(session.Blocks, block)
}
