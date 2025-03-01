package errors

import (
	"errors"
	"fmt"
	"io"
)

// Error(msg) is identical to New(msg)
// they both create a new error that records the source reference where this is done
func New(msg string) BaseError {
	return &msgError{
		baseError: baseError{
			source: GetCaller(2),
		},
		msg: msg,
	}
}

// Error(msg) is identical to New(msg)
// they both create a new error that records the source reference where this is done
func Error(msg string) BaseError {
	return &msgError{
		baseError: baseError{
			source: GetCaller(2),
		},
		msg: msg,
	}
}

// Errorf() is the same as Error() and New(), but does message formatting
func Errorf(format string, args ...interface{}) BaseError {
	return &msgError{
		baseError: baseError{
			source: GetCaller(2),
		},
		msg: fmt.Sprintf(format, args...),
	}
}

// Wrap() an existing error with a message and capturing the source where you wrapped
func Wrap(err error, msg string) BaseError {
	if err == nil {
		return nil
	}
	return &msgError{
		baseError: baseError{
			wrapped: err,
			source:  GetCaller(2),
		},
		msg: msg,
	}
}

// Wrapf() is like Wrap, but does message formatting
func Wrapf(err error, format string, args ...interface{}) BaseError {
	if err == nil {
		return nil
	}
	return &msgError{
		baseError: baseError{
			wrapped: err,
			source:  GetCaller(2),
		},
		msg: fmt.Sprintf(format, args...),
	}
}

// wrappers around go's default package so you do not need to directly import that too
// which will clutter the "errors" namespace in your packages.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

type msgError struct {
	baseError
	msg string
}

// return the error string, not recursing into wrapped errors
func (err msgError) String() string {
	return err.msg
}

// return the error string, and recurse into wrapped errors
// to construct a string like abc because def because xyz
func (err msgError) Error() string {
	if err.wrapped != nil {
		return err.String() + " because " + err.wrapped.Error()
	}
	return err.String()
}

// called when formatting the err with fmt.Printf() like functions
func (err msgError) Format(f fmt.State, c rune) {
	var s string
	//s = fmt.Sprintf("msg(%p)", err)
	switch c {
	case 'v':
		s += fmt.Sprintf("%s:%s", err.source, err.String()) //source with "%s" -> basename
	case 'V':
		s += fmt.Sprintf("%v:%s", err.source, err.String()) //source with "%v" -> fullpath
	default:
		s += err.String() //no source
	}
	io.WriteString(f, s)

	if err.wrapped != nil {
		recurse := false
		if f.Flag('+') {
			recurse = true
			io.WriteString(f, " because ")
		}
		if f.Flag('-') {
			recurse = true
			io.WriteString(f, "\n")
		}
		if recurse {
			if formatter, ok := err.wrapped.(fmt.Formatter); ok {
				formatter.Format(f, c)
			} else {
				io.WriteString(f, err.wrapped.Error())
			}
		}
	}
} //msgError.Format()
