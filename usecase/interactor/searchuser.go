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
	if filter.Email == "" { // if email filter was not informed, find all users
		return c.findAll(ctx)
	}

	return c.findByEmail(ctx, filter.Email)
}

func (c searchUser) findAll(ctx context.Context) (response SearchUserResponseModel, err error) {
	users, err := c.userGateway.FindAll(ctx)
	if err != nil {
		err = fmt.Errorf("find all: %w", err)
		return
	}

	response = userToResponseModel(users)
	return
}

func (c searchUser) findByEmail(ctx context.Context, email string) (response SearchUserResponseModel, err error) {
	var user entity.User
	user, err = c.userGateway.FindByEmail(ctx, email)
	if errors.Is(err, businesserr.ErrCreateUserNotFound) {
		err = nil
		response.Users = make([]SearchUserResponseModelUser, 0)
		return
	}
	if err != nil {
		err = fmt.Errorf("find by email: %w", err)
		return
	}

	response = userToResponseModel([]entity.User{user})
	return
}

func userToResponseModel(users []entity.User) (response SearchUserResponseModel) {
	for _, user := range users {
		response.Users = append(response.Users, SearchUserResponseModelUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return
}
