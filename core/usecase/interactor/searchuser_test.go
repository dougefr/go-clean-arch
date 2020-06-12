package interactor

import (
	"context"
	"errors"
	"github.com/dougefr/go-clean-arch/core/entity"
	"github.com/dougefr/go-clean-arch/core/usecase/businesserr"
	"github.com/dougefr/go-clean-arch/core/usecase/igateway/mock_igateway"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchUser_Execute(t *testing.T) {
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
		userGateway.EXPECT().FindByEmail(context.Background(), "fake@email.com").Return(entity.User{}, expectedErr)

		uc := NewSearchUser(userGateway)
		_, err := uc.Execute(context.Background(), SearchUserRequestModel{Email: "fake@email.com"})

		assert.True(t, errors.Is(err, expectedErr))
	})

	t.Run("should return empty result if there is no user with the email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		userGateway.EXPECT().FindByEmail(context.Background(), "fake@email.com").Return(entity.User{}, businesserr.ErrCreateUserNotFound)

		uc := NewSearchUser(userGateway)
		result, _ := uc.Execute(context.Background(), SearchUserRequestModel{Email: "fake@email.com"})

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
		userGateway.EXPECT().FindByEmail(context.Background(), "fake@email.com").Return(entity.User{
			ID:    1,
			Name:  "fake name",
			Email: "fake@email.com",
		}, nil)

		uc := NewSearchUser(userGateway)
		result, _ := uc.Execute(context.Background(), SearchUserRequestModel{Email: "fake@email.com"})

		assert.Equal(t, SearchUserResponseModel{
			Users: []SearchUserResponseModelUser{
				{
					ID:    1,
					Name:  "fake name",
					Email: "fake@email.com",
				},
			},
		}, result)
	})
}
