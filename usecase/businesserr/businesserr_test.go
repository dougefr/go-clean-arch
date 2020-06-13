// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

package businesserr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBusinessError(t *testing.T) {
	const fakeError = "fakeError"
	const fakeCode = "fakeCode"

	t.Run("should return the correct error string and error code", func(t *testing.T) {
		err := newBusinessError(fakeCode, fakeError)
		assert.Equal(t, fakeError, err.Error())
		assert.Equal(t, fakeCode, err.Code())
	})
}
