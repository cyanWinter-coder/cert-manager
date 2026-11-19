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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
)

func TestGetCertificateFailedIssuanceAttempts(t *testing.T) {
	certificateWithAttempts := func(attempts int) *cmapi.Certificate {
		return &cmapi.Certificate{
			Status: cmapi.CertificateStatus{FailedIssuanceAttempts: &attempts},
		}
	}

	tests := map[string]struct {
		crt  *cmapi.Certificate
		want int
	}{
		"nil certificate returns 0": {
			crt:  nil,
			want: 0,
		},
		"unset failedIssuanceAttempts returns 0": {
			crt:  &cmapi.Certificate{},
			want: 0,
		},
		"zero failedIssuanceAttempts returns 0": {
			crt:  certificateWithAttempts(0),
			want: 0,
		},
		"positive failedIssuanceAttempts is returned as-is": {
			crt:  certificateWithAttempts(3),
			want: 3,
		},
		"negative failedIssuanceAttempts (manual patch) is normalized to 0": {
			crt:  certificateWithAttempts(-1),
			want: 0,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, GetCertificateFailedIssuanceAttempts(test.crt))
		})
	}
}
