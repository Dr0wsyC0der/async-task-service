package domain

import "errors"

var ErrUserAlreadyExists = errors.New("user already exist")

type User struct {
	ID           string
	Login        string
	PasswordHash string `json:"-"`
}
