package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/dougefr/go-clean-arch/core/entity"
	"github.com/dougefr/go-clean-arch/core/usecase/igateway"
)

type (
	// SearchUserRequestModel ...
	SearchUserRequestModel struct {
		Email string
	}

	// SearchUserResponseModel ...
	SearchUserResponseModel struct {
		Users []SearchUserResponseModelUser
	}

	// SearchUserResponseModelUser ...
	SearchUserResponseModelUser struct {
		ID    int64
		Name  string
		Email string
	}

	// SearchUser ...
	SearchUser interface {
		Execute(ctx context.Context, filter SearchUserRequestModel) (SearchUserResponseModel, error)
	}

	searchUser struct {
		userGateway igateway.User
	}
)

// NewSearchUser ...
func NewSearchUser(userGateway igateway.User) SearchUser {
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
