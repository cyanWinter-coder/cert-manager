/*
Copyright 2026 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package issuer

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRetryableError(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		assert.NoError(t, NewRetryableError(nil))
	})

	t.Run("wrapped error preserves message", func(t *testing.T) {
		err := NewRetryableError(errors.New("issuer unavailable"))
		assert.EqualError(t, err, "issuer unavailable")
	})
}

func TestIsRetryableError(t *testing.T) {
	sentinel := errors.New("connection refused")

	tests := map[string]struct {
		err  error
		want bool
	}{
		"nil error is not retryable": {
			err:  nil,
			want: false,
		},
		"plain error is not retryable": {
			err:  errors.New("boom"),
			want: false,
		},
		"wrapped error is retryable": {
			err:  NewRetryableError(sentinel),
			want: true,
		},
		"retryable error wrapped again is still retryable": {
			err:  fmt.Errorf("signing request: %w", NewRetryableError(sentinel)),
			want: true,
		},
		"wrapped sentinel is still matchable with errors.Is": {
			err:  NewRetryableError(sentinel),
			want: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, IsRetryableError(test.err))
		})
	}

	t.Run("errors.Is unwraps to the original error", func(t *testing.T) {
		assert.ErrorIs(t, NewRetryableError(sentinel), sentinel)
	})
}
