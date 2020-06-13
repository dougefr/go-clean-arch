// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package restctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dougefr/go-clean-arch/core/usecase/interactor"
	"github.com/dougefr/go-clean-arch/interface/iinfra"
	"github.com/google/uuid"
)

// User ...
type (
	User interface {
		Create(req RestRequest) RestResponse
		Search(req RestRequest) RestResponse
	}

	user struct {
		ucCreateUser interactor.CreateUser
		ucSearchUser interactor.SearchUser
		session      iinfra.Session
		logger       iinfra.LogProvider
	}

	// create user request body
	createReqBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// create user response body
	createResBody struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// search user response body
	searchResBody struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

// NewUser ...
func NewUser(ucCreateUser interactor.CreateUser,
	ucSearchUser interactor.SearchUser,
	session iinfra.Session,
	logger iinfra.LogProvider) User {
	return user{
		ucCreateUser: ucCreateUser,
		ucSearchUser: ucSearchUser,
		session:      session,
		logger:       logger,
	}
}

// Create ...
func (u user) Create(req RestRequest) (res RestResponse) {
	startTime := time.Now()
	ctx := context.WithValue(context.Background(), iinfra.ContextKeyGlobalLogAttrs, iinfra.LogAttrs{
		"request-id": uuid.New(),
	})
	u.logger.Debug(ctx, "starting create user")

	var reqBody createReqBody
	if err := json.Unmarshal(req.Body, &reqBody); err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when unmarshalling request body: %v", err))
		return respondError(err)
	}

	ucReqModel := interactor.CreateUserRequestModel{
		Name:  reqBody.Name,
		Email: reqBody.Email,
	}

	tx, err := u.session.BeginTx()
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when starting tx: %v", err))
		return respondError(err)
	}

	ctx = context.WithValue(ctx, iinfra.ContextKeyTx, tx) // add Tx to context to be use in gateways
	ucResModel, err := u.ucCreateUser.Execute(ctx, ucReqModel)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when executing core: %v", err))
		_ = u.session.RollbackTx(tx)
		return respondError(err)
	}

	var resBody createResBody
	resBody.ID = strconv.FormatInt(ucResModel.ID, 10) // format to string because int64 can be too big to JS
	resBody.Name = ucResModel.Name
	resBody.Email = ucResModel.Email

	res.Body, _ = json.Marshal(resBody)
	res.StatusCode = http.StatusCreated // 201

	if err = u.session.CommitTx(tx); err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when commiting tx: %v", err))
		return respondError(err)
	}

	u.logger.Debug(ctx, "ending create user method", iinfra.LogAttrs{
		"duration": time.Since(startTime),
	})

	return
}

// Search ...
func (u user) Search(req RestRequest) (res RestResponse) {
	startTime := time.Now()
	ctx := context.WithValue(context.Background(), iinfra.ContextKeyGlobalLogAttrs, iinfra.LogAttrs{
		"request-id": uuid.New(),
	})
	u.logger.Debug(ctx, "starting create user")

	// get filters from query param
	var filter interactor.SearchUserRequestModel
	filter.Email = req.GetQueryParam("email")

	ucResModel, err := u.ucSearchUser.Execute(ctx, filter)
	if err != nil {
		u.logger.Error(ctx, fmt.Sprintf("error when executing core: %v", err))
		return respondError(err)
	}

	var resBody []searchResBody
	for _, modelUser := range ucResModel.Users {
		resBody = append(resBody, searchResBody{
			ID:    strconv.FormatInt(modelUser.ID, 10),
			Name:  modelUser.Name,
			Email: modelUser.Email,
		})
	}

	res.Body, _ = json.Marshal(resBody)
	res.StatusCode = http.StatusOK

	u.logger.Debug(ctx, "ending create user method", iinfra.LogAttrs{
		"duration": time.Since(startTime),
	})

	return
}
