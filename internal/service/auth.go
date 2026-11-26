package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
	"github.com/Dr0wsyC0der/async-task-service/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidRegistration = errors.New("invalid registration")

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) Register(login, password string) error {
	login = strings.TrimSpace(login)
	if login == "" || password == "" || len(password) > 72 {
		return ErrInvalidRegistration
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Login:        login,
		PasswordHash: string(hash),
		ID:           uuid.New().String(),
	}

	return s.userRepo.Create(*user)
}

func (s *AuthService) Login(login, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(strings.TrimSpace(login))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if password == "" || len(password) > 72 {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token := rand.Text()

	session := domain.Session{
		ID:     token,
		UserID: user.ID,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return "", err
	}

	return token, nil

}
