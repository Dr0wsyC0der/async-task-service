package repository

import (
	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
)

type TaskRepository interface {
	Create(task domain.Task) error
	GetByID(id string) (domain.Task, error)
	Update(task domain.Task) error
}
