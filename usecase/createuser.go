package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/dougefr/go-clean-arch/entity"
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

func (c createUser) Execute(ctx context.Context,
	user CreateUserRequestModel) (response CreateUserResponseModel, err error) {
	if user.Name == "" {
		err = ErrCreateUserErrEmptyName
		return
	}
	if user.Email == "" {
		err = ErrCreateUserErrEmptyEmail
		return
	}

	if _, err = c.userGateway.FindByEmail(ctx, user.Email); err != nil && !errors.Is(err, ErrCreateUserNotFound) {
		err = fmt.Errorf("find by email: %w", err)
		return
	}
	if err == nil {
		err = ErrCreateUserAlreadyExists
		return
	}

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
