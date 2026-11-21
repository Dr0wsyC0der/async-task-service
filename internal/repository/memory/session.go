package memory

import (
	"fmt"
	"sync"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
	"github.com/Dr0wsyC0der/async-task-service/internal/repository"
)

type SessionStorage struct {
	sessions map[string]domain.Session
	mux      sync.RWMutex
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		sessions: make(map[string]domain.Session),
	}
}

func (s *SessionStorage) Create(session domain.Session) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, exists := s.sessions[session.ID]
	if exists {
		return fmt.Errorf("session already exists")
	}
	s.sessions[session.ID] = session
	return nil
}

func (s *SessionStorage) GetByID(id string) (domain.Session, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if session, exists := s.sessions[id]; exists {
		return session, nil
	}
	return domain.Session{}, fmt.Errorf("session not found")
}

var _ repository.SessionRepository = (*SessionStorage)(nil)
