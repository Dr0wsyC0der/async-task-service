package memory

import (
	"fmt"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
	"github.com/Dr0wsyC0der/async-task-service/internal/repository"
)

type TaskStorage struct {
	tasks map[string]domain.Task
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks: make(map[string]domain.Task),
	}
}

func (s *TaskStorage) Create(task domain.Task) error {
	s.tasks[task.ID] = task
	return nil
}

func (s *TaskStorage) GetByID(id string) (domain.Task, error) {
	task, exists := s.tasks[id]
	if !exists {
		return domain.Task{}, fmt.Errorf("task with ID %s not found", id)
	}
	return task, nil
}

func (s *TaskStorage) Update(task domain.Task) error {
	_, exists := s.tasks[task.ID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}
	s.tasks[task.ID] = task
	return nil
}

var _ repository.TaskRepository = (*TaskStorage)(nil)
