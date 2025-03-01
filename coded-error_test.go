package errors

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCodedError(t *testing.T) {
	//default error
	err := errors.New("broken")
	assert.False(t, HasCode(err))
	if _, ok := GetCode(err); ok {
		t.Fatal("normal error seems to have a code")
	}

	//coded error
	const someCode = 12345
	err = Code(err, someCode)
	assert.True(t, HasCode(err))
	if code, ok := GetCode(err); !ok {
		t.Fatal("coded has no code")
	} else {
		assert.Equal(t, someCode, code)
	}
	if xerr := errors.Unwrap(err); xerr == nil {
		t.Fatal("not unwrapped")
	} else {
		if !strings.Contains(xerr.Error(), "broken") {
			t.Fatalf("unwrapped != broken: %+v", xerr)
		}
	}

	//wrap again to be retryable and coded
	err = Retry(Wrapf(err, "another issue"), time.Hour*2) //later than original

	//still has the code
	assert.True(t, HasCode(err))
	if code, ok := GetCode(err); !ok {
		t.Fatal("coded has no code")
	} else {
		assert.Equal(t, someCode, code)
	}
	if xerr := errors.Unwrap(err); xerr == nil {
		t.Fatal("not unwrapped")
	} else {
		if !strings.Contains(err.Error(), "broken") {
			t.Fatalf("unwrapped != broken: %+v", xerr)
		}
	}
}
