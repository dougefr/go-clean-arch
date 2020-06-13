// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package businesserr

type (
	// BusinessError ...
	BusinessError interface {
		error
		Code() string
	}

	businessError struct {
		error  string
		code string
	}
)

func newBusinessError(code, error string) BusinessError {
	return businessError{
		error:  error,
		code: code,
	}
}

// Error ...
func (b businessError) Error() string {
	return b.error
}

// Code ...
func (b businessError) Code() string {
	return b.code
}

// Business errors that use cases interactor can results
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
