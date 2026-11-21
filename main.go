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
	userRepo := memory.NewUserStorage()
	sessionRepo := memory.NewSessionStorage()
	authService := service.NewAuthService(userRepo, sessionRepo)
	taskService := service.NewTaskService(memoryRepo)
	taskHandler := handler.NewTaskHandler(taskService)
	authHandler := handler.NewAuthHandler(authService)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /task", taskHandler.CreateTask)
	mux.HandleFunc("GET /status/{task_id}", taskHandler.GetTaskStatus)
	mux.HandleFunc("GET /result/{task_id}", taskHandler.GetTaskResult)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(":8000", mux))
}
