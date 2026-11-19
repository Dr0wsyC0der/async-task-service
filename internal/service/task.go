package service

import (
	"log"
	"time"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
	"github.com/Dr0wsyC0der/async-task-service/internal/repository"
	"github.com/google/uuid"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(code, language string) (domain.Task, error) {
	task := domain.Task{
		Code:     code,
		Language: language,
		Status:   domain.StatusInProgress,
		ID:       uuid.New().String(),
	}
	err := s.repo.Create(task)
	if err != nil {
		return domain.Task{}, err
	}
	go s.processTask(task)
	return task, nil
}

func (s *TaskService) processTask(task domain.Task) {
	time.Sleep(2 * time.Second)
	task.Status = domain.StatusReady
	task.Result = "Processed result for task ID: " + task.ID
	err := s.repo.Update(task)
	if err != nil {
		log.Printf("failed to update task %s: %v", task.ID, err)
	}
}

func (s *TaskService) GetTaskStatus(id string) (domain.Status, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	return task.Status, nil
}

func (s *TaskService) GetTaskResult(id string) (string, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	if task.Status != domain.StatusReady {
		return "", nil
	}
	return task.Result, nil
}
