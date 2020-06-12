package restctrl

import (
	usecase2 "github.com/dougefr/go-clean-arch/core/usecase"
	"net/http"
)

// RestRequest ...
type RestRequest struct {
	GetQueryParam func(key string) string
	Body          []byte
}

// RestResponse ...
type RestResponse struct {
	Body       []byte
	StatusCode int
}

func respondError(err error) (res RestResponse) {
	if be, ok := err.(usecase2.BusinessError); ok {
		res.Body = []byte(be.Error())
		switch be {
		case usecase2.ErrCreateUserNotFound:
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
