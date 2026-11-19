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

import "errors"

// RetryableError marks an error returned by an issuer implementation as
// transient: the operation that failed is expected to succeed if it is
// retried later, for example when the issuer's backing service is
// temporarily unavailable.
//
// This type is the coordination point between two retry mechanisms:
//
//   - Errors wrapped with NewRetryableError are retried by the
//     CertificateRequest controller via its rate-limited workqueue. The
//     CertificateRequest is marked Pending (not Failed), so the
//     Certificate's failure state (status.lastFailureTime and
//     status.failedIssuanceAttempts) is left untouched and the exponential
//     issuance backoff in the certificates-trigger controller is not
//     consumed by transient errors.
//   - Terminal errors must instead be surfaced by the issuer setting the
//     CertificateRequest's Ready condition to False with reason Failed
//     (and returning a nil response). The certificates-issuing controller
//     then records the failure on the Certificate and the
//     certificates-trigger controller retries with exponential backoff.
//
// Issuer implementations are not required to use this type: an unwrapped
// error keeps the historical behavior of being retried by the workqueue
// without any condition being set. Using RetryableError only adds
// visibility (a Pending condition and an event) and documents intent.
type RetryableError struct {
	err error
}

// NewRetryableError wraps err to mark it as a transient failure that the
// controller should retry. It returns nil if err is nil.
func NewRetryableError(err error) error {
	if err == nil {
		return nil
	}
	return &RetryableError{err: err}
}

func (e *RetryableError) Error() string {
	return e.err.Error()
}

// Unwrap allows errors.Is and errors.As to inspect the wrapped error.
func (e *RetryableError) Unwrap() error {
	return e.err
}

// IsRetryableError returns true if err (or any error it wraps) is a
// RetryableError, i.e. the issuer has signaled that the failure is
// transient and the controller should retry rather than treat the request
// as terminally failed.
func IsRetryableError(err error) bool {
	var retryableErr *RetryableError
	return errors.As(err, &retryableErr)
}
