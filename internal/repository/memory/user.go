package memory

import (
	"fmt"
	"sync"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
	"github.com/Dr0wsyC0der/async-task-service/internal/repository"
)

type UserStorage struct {
	users map[string]domain.User
	mu    sync.RWMutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[string]domain.User),
	}
}

func (u *UserStorage) Create(user domain.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	_, exists := u.users[user.Login]
	if exists {
		return domain.ErrUserAlreadyExists
	}
	u.users[user.Login] = user
	return nil
}

func (u *UserStorage) GetByLogin(login string) (domain.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	user, exists := u.users[login]
	if exists {
		return user, nil
	}
	return domain.User{}, fmt.Errorf("user not found")
}

var _ repository.UserRepository = (*UserStorage)(nil)
