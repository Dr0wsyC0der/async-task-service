package repository

import (
	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
)

type UserRepository interface {
	Create(user domain.User) error
	GetByLogin(login string) (domain.User, error)
}
