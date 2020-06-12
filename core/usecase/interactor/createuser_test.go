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

func TestCreateUserExecute(t *testing.T) {
	const fakeEmail = "fake@email.com"
	const fakeName = "fake name"

	t.Run("should return an error ErrCreateUserErrEmptyName when user name is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)

		uc := NewCreateUser(userGateway)
		_, err := uc.Execute(context.Background(), CreateUserRequestModel{
			Email: fakeEmail,
		})

		assert.EqualError(t, err, businesserr.ErrCreateUserErrEmptyName.Error())
	})

	t.Run("should return an error ErrCreateUserErrEmptyEmail when user email is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)

		uc := NewCreateUser(userGateway)
		_, err := uc.Execute(context.Background(), CreateUserRequestModel{
			Name: fakeName,
		})

		assert.EqualError(t, err, businesserr.ErrCreateUserErrEmptyEmail.Error())
	})

	t.Run("should return an unknown error when occur an error when checking if there is no user if the same email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		expectedErr := errors.New("fake-error")
		userGateway.EXPECT().FindByEmail(context.Background(), fakeEmail).Return(entity.User{}, expectedErr)

		uc := NewCreateUser(userGateway)
		_, err := uc.Execute(context.Background(), CreateUserRequestModel{
			Name:  fakeName,
			Email: fakeEmail,
		})

		assert.True(t, errors.Is(err, expectedErr))
	})

	t.Run("should return an error ErrCreateUserAlreadyExists when there is another user using the email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		userGateway.EXPECT().FindByEmail(context.Background(), fakeEmail).Return(entity.User{}, nil)

		uc := NewCreateUser(userGateway)
		_, err := uc.Execute(context.Background(), CreateUserRequestModel{
			Name:  fakeName,
			Email: fakeEmail,
		})

		assert.EqualError(t, err, businesserr.ErrCreateUserAlreadyExists.Error())
	})

	t.Run("should return an unknown error when occur an error when creating an user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		expectedErr := errors.New("fake-error")
		userGateway.EXPECT().FindByEmail(context.Background(), fakeEmail).Return(entity.User{}, businesserr.ErrCreateUserNotFound)
		userGateway.EXPECT().Create(context.Background(), entity.User{
			Name:  fakeName,
			Email: fakeEmail,
		}).Return(entity.User{}, expectedErr)

		uc := NewCreateUser(userGateway)
		_, err := uc.Execute(context.Background(), CreateUserRequestModel{
			Name:  fakeName,
			Email: fakeEmail,
		})

		assert.True(t, errors.Is(err, expectedErr))
	})

	t.Run("should return the user created data when everything goes fine", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userGateway := mock_igateway.NewMockUser(ctrl)
		userGateway.EXPECT().FindByEmail(context.Background(), fakeEmail).Return(entity.User{}, businesserr.ErrCreateUserNotFound)
		userGateway.EXPECT().Create(context.Background(), entity.User{
			Name:  fakeName,
			Email: fakeEmail,
		}).Return(entity.User{
			ID:    1,
			Name:  fakeName,
			Email: fakeEmail,
		}, nil)

		uc := NewCreateUser(userGateway)
		responseModel, _ := uc.Execute(context.Background(), CreateUserRequestModel{
			Name:  fakeName,
			Email: fakeEmail,
		})

		assert.Equal(t, CreateUserResponseModel{
			ID:    1,
			Name:  fakeName,
			Email: fakeEmail,
		}, responseModel)
	})
}
