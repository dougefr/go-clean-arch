package businesserr

import (
	"errors"
)

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
	// ErrCreateUserAlreadyExists ...
	ErrCreateUserAlreadyExists BusinessError = errors.New("user already exists")
)
