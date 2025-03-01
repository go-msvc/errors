package errors

import (
	"errors"
	"fmt"
	"strconv"
)

func Code(err error, code int) CodedError {
	if err == nil {
		return nil
	}
	return codedError{
		baseError: baseError{
			wrapped: err,
			source:  GetCaller(2),
		},
		code: code,
	}
}

func Codef(code int, format string, args ...interface{}) CodedError {
	return codedError{
		baseError: baseError{
			wrapped: fmt.Errorf(format, args...),
			source:  GetCaller(2),
		},
		code: code,
	}
}

func HasCode(err error) bool {
	var re CodedError
	return errors.As(err, &re)
}

func GetCode(err error) (code int, ok bool) {
	var ce CodedError
	if errors.As(err, &ce) {
		return ce.Code(), true
	}
	return code, false
}

type CodedError interface {
	BaseError
	Coded
}

type Coded interface {
	Code() int
}

var _ CodedError = (*codedError)(nil)

type codedError struct {
	baseError
	code int
}

func (err codedError) Code() int {
	return err.code
}

func (err codedError) String() string {
	return strconv.FormatInt(int64(err.code), 10)
}
