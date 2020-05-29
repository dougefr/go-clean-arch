package usecase

import (
	"errors"
	"github.com/dougefr/go-clean-code/entity"
)

// UserRepo ...
type UserRepo interface {
	FindByEmail(email string) (entity.User, error)
	CreateUser(user entity.User) (entity.User, error)
}

// BusinessError ...
type BusinessError interface {
	error
}

var (
	// ErrCreateUserNotFound ...
	ErrCreateUserNotFound BusinessError = errors.New("not found")
	// ErrCreateUserErrEmptyName ...
	ErrCreateUserErrEmptyName BusinessError = errors.New("user name cannot be empty")
	// ErrCreateUserErrEmptyEmail ...
	ErrCreateUserErrEmptyEmail BusinessError = errors.New("user email cannot be empty")
)
