package controller

import (
	"encoding/json"
	"github.com/dougefr/go-clean-code/usecase"
	"net/http"
)

// User ...
type User interface {
	CreateUser(req RestRequest) RestResponse
}

type user struct {
	ucCreateUser usecase.CreateUser
}

// NewUser ...
func NewUser(ucCreateUser usecase.CreateUser) User {
	return user{
		ucCreateUser: ucCreateUser,
	}
}

// CreateUser ...
func (u user) CreateUser(req RestRequest) (res RestResponse) {
	var reqBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var err error
	if err = json.Unmarshal(req.Body, &reqBody); err != nil {
		return respondError(err)
	}

	ucReqModel := usecase.CreateUserRequestModel{
		Name:  reqBody.Name,
		Email: reqBody.Email,
	}

	var ucResModel usecase.CreateUserResponseModel
	if ucResModel, err = u.ucCreateUser.Execute(ucReqModel); err != nil {
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
		return respondError(err)
	}

	res.StatusCode = http.StatusCreated

	return
}

func respondError(err error) (res RestResponse) {
	if be, ok := err.(usecase.BusinessError); ok {
		res.Body = []byte(be.Error())
		switch be {
		case usecase.ErrCreateUserNotFound:
			res.StatusCode = http.StatusNotFound
		default:
			res.StatusCode = http.StatusBadRequest
		}

		return
	}

	res.Body = []byte("internal server error")
	res.StatusCode = http.StatusInternalServerError

	return
}
