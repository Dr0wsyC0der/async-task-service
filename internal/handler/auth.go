package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
	"github.com/Dr0wsyC0der/async-task-service/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

type CredentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req CredentialsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req CredentialsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Register(req.Username, req.Password); err != nil {
		if errors.Is(err, service.ErrInvalidRegistration) {
			http.Error(w, "invalid registration", http.StatusBadRequest)
			return
		} else if errors.Is(err, domain.ErrUserAlreadyExists) {
			http.Error(w, "user already exist", http.StatusConflict)
			return
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}
