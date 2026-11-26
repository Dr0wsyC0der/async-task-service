package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
	"github.com/Dr0wsyC0der/async-task-service/internal/service"
)

type CreateTaskRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type CreateTaskResponse struct {
	ID string `json:"task_id"`
}

type TaskStatusResponse struct {
	Status domain.Status `json:"status"`
}

type TaskResultResponse struct {
	Result string `json:"result"`
}

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask godoc
// @Summary      Create a new task
// @Description  Submit code for asynchronous processing
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateTaskRequest true "Task data"
// @Success      201 {object} CreateTaskResponse
// @Failure      400 {string} string "Invalid request"
// @Failure      401 {string} string "Missing or invalid session token"
// @Failure      500 {string} string "Internal server error"
// @Router       /task [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(req.Code, req.Language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateTaskResponse{ID: task.ID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetTaskStatus godoc
// @Summary      Get task status
// @Description  Returns the current processing status of a task
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        task_id path string true "Task ID"
// @Success      200 {object} TaskStatusResponse
// @Failure      401 {string} string "Missing or invalid session token"
// @Failure      404 {string} string "Task not found"
// @Router       /status/{task_id} [get]
func (h *TaskHandler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("task_id")
	status, err := h.service.GetTaskStatus(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resp := TaskStatusResponse{Status: status}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetTaskResult godoc
// @Summary      Get task result
// @Description  Returns the result of a completed task
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        task_id path string true "Task ID"
// @Success      200 {object} TaskResultResponse
// @Success      202 "Task is still processing"
// @Failure      401 {string} string "Missing or invalid session token"
// @Failure      404 {string} string "Task not found"
// @Router       /result/{task_id} [get]
func (h *TaskHandler) GetTaskResult(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("task_id")
	result, err := h.service.GetTaskResult(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if result == "" {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	resp := TaskResultResponse{Result: result}
	json.NewEncoder(w).Encode(resp)
}
