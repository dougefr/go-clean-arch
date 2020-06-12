package restctrl

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/dougefr/go-clean-arch/core/usecase/businesserr"
	"github.com/dougefr/go-clean-arch/core/usecase/interactor"
	"github.com/dougefr/go-clean-arch/core/usecase/interactor/mock_interactor"

	"github.com/dougefr/go-clean-arch/interface/iinfra/mock_iinfra"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUser_Create(t *testing.T) {
	t.Run("should results in StatusInternalServerError if the request body is an invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any())
		logger.EXPECT().Error(gomock.Any(), gomock.Any())
		c := NewUser(nil, nil, nil, logger)
		res := c.Create(RestRequest{
			Body: []byte("I'm an invalid JSON"),
		})

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should results in StatusInternalServerError if open a new Tx on database results in an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any())
		logger.EXPECT().Error(gomock.Any(), gomock.Any())

		session := mock_iinfra.NewMockSession(ctrl)
		expectedErr := errors.New("fake-error")
		session.EXPECT().BeginTx().Return(nil, expectedErr)

		c := NewUser(nil, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(`{"name":"fake name","email":"fake@email.com"}`),
		})

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should results in StatusBadRequest if usecase interactor return any business error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any())
		logger.EXPECT().Error(gomock.Any(), gomock.Any())

		session := mock_iinfra.NewMockSession(ctrl)
		session.EXPECT().BeginTx().Return(mock_iinfra.NewMockTx(ctrl), nil)
		session.EXPECT().RollbackTx(gomock.Any()).Return(nil)

		ucCreateUser := mock_interactor.NewMockCreateUser(ctrl)
		expectedErr := businesserr.ErrCreateUserAlreadyExists
		ucCreateUser.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(interactor.CreateUserResponseModel{}, expectedErr)

		c := NewUser(ucCreateUser, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(`{"name":"fake name","email":"fake@email.com"}`),
		})

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("should results in StatusInternalServerError if usecase interactor return any unknown error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any())
		logger.EXPECT().Error(gomock.Any(), gomock.Any())

		session := mock_iinfra.NewMockSession(ctrl)
		session.EXPECT().BeginTx().Return(mock_iinfra.NewMockTx(ctrl), nil)
		session.EXPECT().RollbackTx(gomock.Any()).Return(nil)

		ucCreateUser := mock_interactor.NewMockCreateUser(ctrl)
		expectedErr := errors.New("fake-error")
		ucCreateUser.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(interactor.CreateUserResponseModel{}, expectedErr)

		c := NewUser(ucCreateUser, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(`{"name":"fake name","email":"fake@email.com"}`),
		})

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should results in StatusInternalServerError when commiting Tx results in an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any())
		logger.EXPECT().Error(gomock.Any(), gomock.Any())

		session := mock_iinfra.NewMockSession(ctrl)
		session.EXPECT().BeginTx().Return(mock_iinfra.NewMockTx(ctrl), nil)
		expectedErr := errors.New("fake-error")
		session.EXPECT().CommitTx(gomock.Any()).Return(expectedErr)

		ucCreateUser := mock_interactor.NewMockCreateUser(ctrl)
		ucCreateUser.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(interactor.CreateUserResponseModel{}, nil)

		c := NewUser(ucCreateUser, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(`{"name":"fake name","email":"fake@email.com"}`),
		})

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should results in StatusCreated and returns the created user when everything goes fine", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

		session := mock_iinfra.NewMockSession(ctrl)
		session.EXPECT().BeginTx().Return(mock_iinfra.NewMockTx(ctrl), nil)
		session.EXPECT().CommitTx(gomock.Any()).Return(nil)

		ucCreateUser := mock_interactor.NewMockCreateUser(ctrl)
		ucCreateUser.EXPECT().Execute(gomock.Any(), interactor.CreateUserRequestModel{
			Name:  "fake name",
			Email: "fake@email.com",
		}).
			Return(interactor.CreateUserResponseModel{
				ID:    1,
				Name:  "fake name",
				Email: "fake@email.com",
			}, nil)

		c := NewUser(ucCreateUser, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(`{"name":"fake name","email":"fake@email.com"}`),
		})

		var resBody createResBody
		json.Unmarshal(res.Body, &resBody)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, createResBody{
			ID:    "1",
			Name:  "fake name",
			Email: "fake@email.com",
		}, resBody)
	})
}
