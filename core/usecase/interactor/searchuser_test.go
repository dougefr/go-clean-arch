// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/dougefr/go-clean-arch/core/entity"
	"github.com/dougefr/go-clean-arch/core/usecase/businesserr"
	"github.com/dougefr/go-clean-arch/core/usecase/igateway/mock_igateway"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSearchUserExecute(t *testing.T) {
	const fakeEmail = "fake@email.com"

	t.Run("should return an error if an error was returned when finding all users", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		expectedErr := errors.New("fake-error")
		userGateway.EXPECT().FindAll(context.Background()).Return(nil, expectedErr)

		uc := NewSearchUser(userGateway)
		_, err := uc.Execute(context.Background(), SearchUserRequestModel{})

		assert.True(t, errors.Is(err, expectedErr))
	})

	t.Run("should return an error if an error was returned when finding users by email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		expectedErr := errors.New("fake-error")
		userGateway.EXPECT().FindByEmail(context.Background(), fakeEmail).Return(entity.User{}, expectedErr)

		uc := NewSearchUser(userGateway)
		_, err := uc.Execute(context.Background(), SearchUserRequestModel{Email: fakeEmail})

		assert.True(t, errors.Is(err, expectedErr))
	})

	t.Run("should return empty result if there is no user with the email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		userGateway.EXPECT().FindByEmail(context.Background(), fakeEmail).Return(entity.User{}, businesserr.ErrCreateUserNotFound)

		uc := NewSearchUser(userGateway)
		result, _ := uc.Execute(context.Background(), SearchUserRequestModel{Email: fakeEmail})

		assert.Empty(t, result.Users)
	})

	t.Run("should return all users when no email was used as filter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		userGateway.EXPECT().FindAll(context.Background()).Return([]entity.User{
			{
				ID:    1,
				Name:  "fake name 1",
				Email: "fake1@email.com",
			},
			{
				ID:    2,
				Name:  "fake name 2",
				Email: "fake2@email.com",
			},
		}, nil)

		uc := NewSearchUser(userGateway)
		result, _ := uc.Execute(context.Background(), SearchUserRequestModel{})

		assert.Equal(t, SearchUserResponseModel{
			Users: []SearchUserResponseModelUser{
				{
					ID:    1,
					Name:  "fake name 1",
					Email: "fake1@email.com",
				},
				{
					ID:    2,
					Name:  "fake name 2",
					Email: "fake2@email.com",
				},
			},
		}, result)
	})

	t.Run("should return one user when an user exists with the email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		userGateway.EXPECT().FindByEmail(context.Background(), fakeEmail).Return(entity.User{
			ID:    1,
			Name:  "fake name",
			Email: fakeEmail,
		}, nil)

		uc := NewSearchUser(userGateway)
		result, _ := uc.Execute(context.Background(), SearchUserRequestModel{Email: fakeEmail})

		assert.Equal(t, SearchUserResponseModel{
			Users: []SearchUserResponseModelUser{
				{
					ID:    1,
					Name:  "fake name",
					Email: fakeEmail,
				},
			},
		}, result)
	})
}
