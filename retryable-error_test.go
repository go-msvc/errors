package errors

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryableError(t *testing.T) {
	//default error
	const msg = "broken"
	orgErr := errors.New(msg)
	assert.Equal(t, msg, orgErr.Error())
	assert.False(t, IsRetryable(orgErr))
	if _, ok := RetryableAt(orgErr); ok {
		t.Fatal("normal error seems retryable")
	}

	//wrap to make a retryable error
	var err error
	err = Retry(orgErr, time.Hour)
	//the message did not change, just retry was added
	assert.Equal(t, msg, err.Error())
	//it is now retryable
	assert.True(t, IsRetryable(err))
	if at, ok := RetryableAt(err); !ok {
		t.Fatal("retryable not retryable")
	} else {
		wait := time.Until(at)
		if wait > time.Hour || wait < time.Minute*59 { //expect just less than an hour
			t.Fatalf("wait %v not 59..60min", wait)
		}
	}
	//and retry can be unwrapped to get the original error which is not retryable
	if err := errors.Unwrap(err); err == nil {
		t.Fatal("not unwrapped")
	} else {
		assert.Equal(t, msg, err.Error())
		assert.False(t, IsRetryable(err))
	}

	//wrapping with another errors, it remains retryable
	err = Wrap(err, "something else failed")
	assert.True(t, IsRetryable(err))
	assert.Equal(t, "something else failed because broken", err.Error())

	//wrapping with a formatted error, still retryable
	const someValue = 555
	err = Wrapf(err, "this is wrong(%d)", someValue)
	expMsg := fmt.Sprintf("this is wrong(%d) because something else failed because broken", someValue)
	assert.True(t, IsRetryable(err))
	assert.Equal(t, expMsg, err.Error())

	//wrapping with a code, still retryable
	const someCode = 995432
	err = Code(err, someCode)
	assert.True(t, IsRetryable(err))
	assert.Equal(t, expMsg, err.Error())

	//unwind to the formatter message, then the previous then on retryable
	err = Unwrap(err)
	assert.True(t, IsRetryable(err))
	err = Unwrap(err)
	assert.True(t, IsRetryable(err))
	err = Unwrap(err)
	assert.True(t, IsRetryable(err))
	err = Unwrap(err)
	assert.False(t, IsRetryable(err))
}
