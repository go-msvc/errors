package errors

import (
	"errors"
	"fmt"
	"time"
)

func Retry(err error, wait time.Duration) RetryableError {
	if err == nil {
		return nil
	}
	return retryableError{
		baseError: baseError{
			wrapped: err,
			source:  GetCaller(2),
		},
		at: time.Now().Add(wait),
	}
}

func Retryf(wait time.Duration, format string, args ...interface{}) RetryableError {
	return retryableError{
		baseError: baseError{
			wrapped: fmt.Errorf(format, args...),
			source:  GetCaller(2),
		},
		at: time.Now().Add(wait),
	}
}

func IsRetryable(err error) bool {
	var re RetryableError
	return errors.As(err, &re)
}

func RetryableAt(err error) (at time.Time, ok bool) {
	var re RetryableError
	if errors.As(err, &re) {
		return re.CanRetryAt(), true
	}
	return at, false
}

type RetryableError interface {
	BaseError
	Retryable
}

type Retryable interface {
	CanRetryAt() time.Time //do not retry before this time
}

var _ RetryableError = (*retryableError)(nil)

type retryableError struct {
	baseError
	at time.Time
}

func (err retryableError) CanRetryAt() time.Time {
	return err.at
}

func (err retryableError) String() string {
	return err.at.String()
}
