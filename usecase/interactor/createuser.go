// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package interactor

import (
	"context"
	"errors"
	"fmt"

	"github.com/dougefr/go-clean-arch/entity"
	"github.com/dougefr/go-clean-arch/usecase/businesserr"
	"github.com/dougefr/go-clean-arch/usecase/igateway"
)

type (
	// CreateUserRequestModel ...
	CreateUserRequestModel struct {
		Name  string
		Email string
	}

	// CreateUserResponseModel ...
	CreateUserResponseModel struct {
		ID    int64
		Name  string
		Email string
	}

	// CreateUser ...
	CreateUser interface {
		Execute(ctx context.Context, user CreateUserRequestModel) (CreateUserResponseModel, error)
	}

	createUser struct {
		userGateway igateway.User
	}
)

// NewCreateUser ...
func NewCreateUser(userGateway igateway.User) CreateUser {
	return createUser{
		userGateway: userGateway,
	}
}

// Execute ...
func (c createUser) Execute(ctx context.Context,
	user CreateUserRequestModel) (response CreateUserResponseModel, err error) {
	// Static validations
	if user.Name == "" {
		err = businesserr.ErrCreateUserErrEmptyName
		return
	}
	if user.Email == "" {
		err = businesserr.ErrCreateUserErrEmptyEmail
		return
	}

	// Check if an user exists with the same email
	if _, err = c.userGateway.FindByEmail(ctx, user.Email); err != nil &&
		!errors.Is(err, businesserr.ErrCreateUserNotFound) {
		err = fmt.Errorf("find by email: %w", err)
		return
	}
	if err == nil {
		err = businesserr.ErrCreateUserAlreadyExists
		return
	}

	// Create the user
	userCreated, err := c.userGateway.Create(ctx, entity.User{
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		err = fmt.Errorf("create user: %w", err)
		return
	}

	response.ID = userCreated.ID
	response.Name = userCreated.Name
	response.Email = userCreated.Email

	return
}
