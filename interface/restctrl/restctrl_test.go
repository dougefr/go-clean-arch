// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package restctrl

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dougefr/go-clean-arch/core/usecase/businesserr"
	"github.com/stretchr/testify/assert"
)

func TestRespondError(t *testing.T) {
	t.Run("should results StatusNotFound when receive ErrCreateUserNotFound", func(t *testing.T) {
		res := respondError(businesserr.ErrCreateUserNotFound)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("should results StatusBadRequest when receive a business error", func(t *testing.T) {
		res := respondError(businesserr.ErrCreateUserErrEmptyEmail)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("should results StatusInternalServerError when receive an unknown error", func(t *testing.T) {
		res := respondError(errors.New("fake error"))
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
