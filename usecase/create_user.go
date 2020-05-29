package usecase

import (
	"fmt"
	"github.com/dougefr/go-clean-code/entity"
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
	Execute(user CreateUserRequestModel) (CreateUserResponseModel, error)
}

type createUser struct {
	userRepo UserRepo
}

// NewCreateUser ...
func NewCreateUser(userRepo UserRepo) CreateUser {
	return createUser{
		userRepo: userRepo,
	}
}

func (c createUser) Execute(user CreateUserRequestModel) (response CreateUserResponseModel, err error) {
	if user.Name == "" {
		err = ErrCreateUserErrEmptyName
		return
	}
	if user.Email == "" {
		err = ErrCreateUserErrEmptyEmail
		return
	}

	/*if _, err = c.userRepo.FindByEmail(user.Email); !errors.Is(err, ErrCreateUserNotFound) {
		err = fmt.Errorf("find by email: %w", err)
		return
	}*/

	userCreated, err := c.userRepo.CreateUser(entity.User{
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
