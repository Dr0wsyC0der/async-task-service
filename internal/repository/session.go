package repository

import (
	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
)

type SessionRepository interface {
	Create(session domain.Session) error
	GetByID(id string) (domain.Session, error)
}
