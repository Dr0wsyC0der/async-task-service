package main

import (
	"log"
	"net/http"

	_ "github.com/Dr0wsyC0der/async-task-service/docs"
	"github.com/Dr0wsyC0der/async-task-service/internal/handler"
	"github.com/Dr0wsyC0der/async-task-service/internal/repository/memory"
	"github.com/Dr0wsyC0der/async-task-service/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Code Execution Service API
// @version         1.0
// @description     API for asynchronous code execution tasks.
// @host             localhost:8000
// @BasePath         /
func main() {
	memoryRepo := memory.NewTaskStorage()
	taskService := service.NewTaskService(memoryRepo)
	taskHandler := handler.NewTaskHandler(taskService)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /task", taskHandler.CreateTask)
	mux.HandleFunc("GET /status/{task_id}", taskHandler.GetTaskStatus)
	mux.HandleFunc("GET /result/{task_id}", taskHandler.GetTaskResult)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(":8000", mux))
}
