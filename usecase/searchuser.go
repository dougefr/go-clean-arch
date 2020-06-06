package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/dougefr/go-clean-arch/entity"
	"github.com/dougefr/go-clean-arch/usecase/gateway"
)

// SearchUserRequestModel ...
type SearchUserRequestModel struct {
	Email string
}

// SearchUserResponseModel ...
type SearchUserResponseModel struct {
	Users []SearchUserResponseModelUser
}

// SearchUserResponseModelUser ...
type SearchUserResponseModelUser struct {
	ID    uint
	Name  string
	Email string
}

// SearchUser ...
type SearchUser interface {
	Execute(ctx context.Context, filter SearchUserRequestModel) (SearchUserResponseModel, error)
}

type searchUser struct {
	userGateway gateway.User
}

// NewSearchUser ...
func NewSearchUser(userGateway gateway.User) SearchUser {
	return searchUser{
		userGateway: userGateway,
	}
}

func (c searchUser) Execute(ctx context.Context,
	filter SearchUserRequestModel) (response SearchUserResponseModel, err error) {

	var users []entity.User
	if filter.Email == "" {
		users, err = c.userGateway.FindAll(ctx)
		if err != nil {
			err = fmt.Errorf("find all: %w", err)
			return
		}
	} else {
		var user entity.User
		user, err = c.userGateway.FindByEmail(ctx, filter.Email)
		if errors.Is(err, ErrCreateUserNotFound) {
			err = nil
			response.Users = make([]SearchUserResponseModelUser, 0)
			return
		}
		if err != nil {
			err = fmt.Errorf("find by email: %w", err)
			return
		}

		users = []entity.User{user}
	}

	for _, user := range users {
		response.Users = append(response.Users, SearchUserResponseModelUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return
}
