// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package restctrl

import (
	"net/http"

	"github.com/dougefr/go-clean-arch/usecase/businesserr"
)

// RestRequest ...
type (
	RestRequest struct {
		GetQueryParam func(key string) string
		Body          []byte
	}

	// RestResponse ...
	RestResponse struct {
		Body       []byte
		StatusCode int
	}
)

func respondError(err error) (res RestResponse) {
	if be, ok := err.(businesserr.BusinessError); ok {
		res.Body = []byte(be.Error())
		switch be {
		case businesserr.ErrCreateUserNotFound:
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
