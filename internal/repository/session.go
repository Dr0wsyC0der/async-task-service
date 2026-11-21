package repository

import (
	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
)

type SessionRepository interface {
	Create(session domain.Session) error
	GetByLogin(login string) (domain.Session, error)
}
