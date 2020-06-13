// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package businesserr

import (
	"errors"
)

type (
	// BusinessError ...
	BusinessError interface {
		error
		Code() string
	}

	businessError struct {
		err  error
		code string
	}
)

func newBusinessError(code, error string) BusinessError {
	return businessError{
		err:  errors.New(error),
		code: code,
	}
}

func (b businessError) Error() string {
	return b.err.Error()
}

func (b businessError) Code() string {
	return b.code
}

var (
	// ErrCreateUserNotFound ...
	ErrCreateUserNotFound = newBusinessError("ErrCreateUserNotFound", "not found")
	// ErrCreateUserErrEmptyName ...
	ErrCreateUserErrEmptyName = newBusinessError("ErrCreateUserErrEmptyName", "user name cannot be empty")
	// ErrCreateUserErrEmptyEmail ...
	ErrCreateUserErrEmptyEmail = newBusinessError("ErrCreateUserErrEmptyEmail", "user email cannot be empty")
	// ErrCreateUserAlreadyExists ...
	ErrCreateUserAlreadyExists = newBusinessError("ErrCreateUserAlreadyExists", "user already exists")
)
