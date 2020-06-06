package restctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dougefr/go-clean-arch/interface/iinfra"
	"github.com/dougefr/go-clean-arch/usecase"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// User ...
type User interface {
	CreateUser(req RestRequest) RestResponse
	SearchUser(req RestRequest) RestResponse
}

type user struct {
	ucCreateUser usecase.CreateUser
	ucSearchUser usecase.SearchUser
	session      iinfra.Session
	logger       iinfra.LogProvider
}

// NewUser ...
func NewUser(ucCreateUser usecase.CreateUser,
	ucSearchUser usecase.SearchUser,
	session iinfra.Session,
	logger iinfra.LogProvider) User {
	return user{
		ucCreateUser: ucCreateUser,
		ucSearchUser: ucSearchUser,
		session:      session,
		logger:       logger,
	}
}

// CreateUser ...
func (u user) CreateUser(req RestRequest) (res RestResponse) {
	startTime := time.Now()
	ctx := context.WithValue(context.Background(), iinfra.ContextKeyGlobalLogAttrs, iinfra.LogAttrs{
		"request-id": uuid.New(),
	})
	u.logger.Debug(ctx, "starting create user")

	var reqBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var err error
	if err = json.Unmarshal(req.Body, &reqBody); err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when unmarshalling request body: %v", err))
		return respondError(err)
	}

	ucReqModel := usecase.CreateUserRequestModel{
		Name:  reqBody.Name,
		Email: reqBody.Email,
	}

	tx, err := u.session.BeginTx()
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when starting tx: %v", err))
		return respondError(err)
	}

	ctx = context.WithValue(ctx, iinfra.ContextKeyTx, tx)
	ucResModel, err := u.ucCreateUser.Execute(ctx, ucReqModel)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when executing usecase: %v", err))
		if err2 := u.session.RollbackTx(tx); err2 != nil {
			u.logger.Error(ctx, fmt.Sprintf("error when rollbacking tx: %v", err2))
			return respondError(err2)
		}
		return respondError(err)
	}

	var resBody struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	resBody.ID = ucResModel.ID
	resBody.Name = ucResModel.Name
	resBody.Email = ucResModel.Email

	if res.Body, err = json.Marshal(resBody); err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when marshalling response body: %v", err))
		if err2 := u.session.RollbackTx(tx); err2 != nil {
			u.logger.Error(ctx, fmt.Sprintf("error when rollbacking tx: %v", err2))
			return respondError(err2)
		}
		return respondError(err)
	}

	res.StatusCode = http.StatusCreated

	if err = u.session.CommitTx(tx); err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when commiting tx: %v", err))
		return respondError(err)
	}

	u.logger.Debug(ctx, "ending create user method", iinfra.LogAttrs{
		"duration": time.Since(startTime),
	})

	return
}

// SearchUser ...
func (u user) SearchUser(req RestRequest) (res RestResponse) {
	startTime := time.Now()
	ctx := context.WithValue(context.Background(), iinfra.ContextKeyGlobalLogAttrs, iinfra.LogAttrs{
		"request-id": uuid.New(),
	})
	u.logger.Debug(ctx, "starting create user")

	var filter usecase.SearchUserRequestModel
	filter.Email = req.GetQueryParam("email")

	ucResModel, err := u.ucSearchUser.Execute(ctx, filter)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when executing usecase: %v", err))
		return respondError(err)
	}

	type resBodyType struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var resBody []resBodyType
	for _, modelUser := range ucResModel.Users {
		resBody = append(resBody, resBodyType{
			ID:    modelUser.ID,
			Name:  modelUser.Name,
			Email: modelUser.Email,
		})
	}

	if res.Body, err = json.Marshal(resBody); err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when marshalling response body: %v", err))
		return respondError(err)
	}

	u.logger.Debug(ctx, "ending create user method", iinfra.LogAttrs{
		"duration": time.Since(startTime),
	})

	return
}
