package memory

import (
	"fmt"
	"sync"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
	"github.com/Dr0wsyC0der/async-task-service/internal/repository"
)

type UserStorage struct {
	users map[string]domain.User
	mux   sync.RWMutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[string]domain.User),
	}
}

func (u *UserStorage) Create(user domain.User) error {
	u.mux.Lock()
	defer u.mux.Unlock()
	_, exists := u.users[user.Login]
	if exists {
		return fmt.Errorf("user already exists")
	}
	u.users[user.Login] = user
	return nil
}

func (u *UserStorage) GetByLogin(login string) (domain.User, error) {
	u.mux.RLock()
	defer u.mux.RUnlock()
	user, exists := u.users[login]
	if exists {
		return user, nil
	}
	return domain.User{}, fmt.Errorf("user not found")
}

var _ repository.UserRepository = (*UserStorage)(nil)
