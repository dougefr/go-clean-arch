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

func TestUserCreate(t *testing.T) {
	const fakeJSON = `{"name":"fake name","email":"fake@email.com"}`
	const fakeName = "fake name"
	const fakeEmail = "fake@email.com"
	fakeError := errors.New("fake-error")

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
		session.EXPECT().BeginTx().Return(nil, fakeError)

		c := NewUser(nil, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(fakeJSON),
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
		ucCreateUser.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(interactor.CreateUserResponseModel{}, businesserr.ErrCreateUserErrEmptyEmail)

		c := NewUser(ucCreateUser, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(fakeJSON),
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
		ucCreateUser.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(interactor.CreateUserResponseModel{}, fakeError)

		c := NewUser(ucCreateUser, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(fakeJSON),
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
		session.EXPECT().CommitTx(gomock.Any()).Return(fakeError)

		ucCreateUser := mock_interactor.NewMockCreateUser(ctrl)
		ucCreateUser.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(interactor.CreateUserResponseModel{}, nil)

		c := NewUser(ucCreateUser, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(fakeJSON),
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
			Name:  fakeName,
			Email: fakeEmail,
		}).
			Return(interactor.CreateUserResponseModel{
				ID:    1,
				Name:  fakeName,
				Email: fakeEmail,
			}, nil)

		c := NewUser(ucCreateUser, nil, session, logger)
		res := c.Create(RestRequest{
			Body: []byte(fakeJSON),
		})

		var resBody createResBody
		json.Unmarshal(res.Body, &resBody)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, createResBody{
			ID:    "1",
			Name:  fakeName,
			Email: fakeEmail,
		}, resBody)
	})
}

func TestUserSearch(t *testing.T) {
	fakeError := errors.New("fake-error")
	const fakeName = "fake name"
	const fakeEmail = "fake@email.com"

	t.Run("should results in StatusInternalServerError if usecase interactor return any unknown error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any())
		logger.EXPECT().Error(gomock.Any(), gomock.Any())

		ucSearchUser := mock_interactor.NewMockSearchUser(ctrl)
		ucSearchUser.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(interactor.SearchUserResponseModel{}, fakeError)

		c := NewUser(nil, ucSearchUser, nil, logger)
		res := c.Search(RestRequest{
			GetQueryParam: func(key string) string {
				return fakeEmail
			},
		})

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should results in StatusOK if usecase interactor when everything goes fine", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mock_iinfra.NewMockLogProvider(ctrl)
		logger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

		ucSearchUser := mock_interactor.NewMockSearchUser(ctrl)
		ucSearchUser.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(interactor.SearchUserResponseModel{
			Users: []interactor.SearchUserResponseModelUser{
				{
					ID:    1,
					Name:  fakeName,
					Email: fakeEmail,
				},
				{
					ID:    2,
					Name:  fakeName,
					Email: fakeEmail,
				},
			},
		}, nil)

		c := NewUser(nil, ucSearchUser, nil, logger)
		res := c.Search(RestRequest{
			GetQueryParam: func(key string) string {
				return fakeEmail
			},
		})

		var resBody []searchResBody
		json.Unmarshal(res.Body, &resBody)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, []searchResBody{
			{
				ID:    "1",
				Name:  fakeName,
				Email: fakeEmail,
			},
			{
				ID:    "2",
				Name:  fakeName,
				Email: fakeEmail,
			},
		}, resBody)
	})
}
