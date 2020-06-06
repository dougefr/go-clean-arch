package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/dougefr/go-clean-arch/entity"
	"github.com/dougefr/go-clean-arch/usecase/gateway"
)

// CreateUserRequestModel ...
type CreateUserRequestModel struct {
	Name  string
	Email string
}

// CreateUserResponseModel ...
type CreateUserResponseModel struct {
	ID    int
	Name  string
	Email string
}

// CreateUser ...
type CreateUser interface {
	Execute(ctx context.Context, user CreateUserRequestModel) (CreateUserResponseModel, error)
}

type createUser struct {
	userGateway gateway.User
}

// NewCreateUser ...
func NewCreateUser(userGateway gateway.User) CreateUser {
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

	userCreated, err := c.userGateway.CreateUser(ctx, entity.User{
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		err = fmt.Errorf("create user: %w", err)
		return
	}

	response.Name = userCreated.Name
	response.Email = userCreated.Email

	return
}
