package errors

import (
	"fmt"
	"io"

	"github.com/go-msvc/logger"
)

type IError interface {
	error
	Parent() error
	Caller() logger.Caller
	Message() string
}

//msError implements IError
type msError struct {
	parent error
	caller logger.Caller
	msg    string
}

func (e msError) Parent() error {
	return e.parent
}

func (e msError) Caller() logger.Caller {
	return e.caller
}

func (e msError) Message() string {
	return e.msg
}

func Is(err error, check error) bool {
	if err == check {
		return true
	}
	p := err
	for i := 0; i < 10; i++ {
		if stack, ok := p.(IError); ok {
			p = stack.Parent()
			if p == nil {
				break
			}
			if p == check {
				return true
			}
		}
	}
	return false
}

//implement error
func (e msError) Error() string {
	s := e.msg
	err := e.parent
	for err != nil {
		if e, ok := err.(IError); ok {
			s += " because " + e.Error()
			err = e.Parent()
		} else {
			s += " because " + err.Error()
			break
		}
	}
	return s
}

func (e msError) CallerError() string {
	s := fmt.Sprintf("%s", e.msg)
	err := e.parent
	for err != nil {
		if e, ok := err.(IError); ok {
			s += " because " + e.Error()
			err = e.Parent()
		} else {
			s += " because " + err.Error()
			break
		}
	}
	return s
}

//implement fmt.Formatter
func (e msError) Format(f fmt.State, c rune) {
	var s string
	switch c {
	case 's':
		s = e.msg

	case 'v':
		s = fmt.Sprintf("%v:%s", e.caller, e.msg)
	case 'V':
		s = fmt.Sprintf("%+V:%s", e.caller, e.msg)

	case 'f':
		s = fmt.Sprintf("%f:%s", e.caller, e.msg)
	case 'F':
		s = fmt.Sprintf("%+F:%s", e.caller, e.msg)

	default:
		s = e.Error()
	}
	io.WriteString(f, s)

	if e.parent != nil {
		stack := false
		if f.Flag('+') {
			stack = true
			io.WriteString(f, " because ")
		}
		if f.Flag('-') {
			stack = true
			io.WriteString(f, "\n")
		}
		if stack {
			if formatter, ok := e.parent.(fmt.Formatter); ok {
				formatter.Format(f, c)
			} else {
				io.WriteString(f, e.parent.Error())
			}
		}
	}
}

func Error(msg string) error {
	return msError{
		parent: nil,
		caller: logger.GetCaller(2),
		msg:    msg,
	}
}

func Errorf(format string, args ...interface{}) error {
	return msError{
		parent: nil,
		caller: logger.GetCaller(2),
		msg:    fmt.Sprintf(format, args...),
	}
}

func Wrap(err error, msg string) error {
	return wrap(err, 3, msg)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return wrap(err, 3, fmt.Sprintf(format, args...))
}

func wrap(err error, skip int, msg string) msError {
	if err == nil {
		Error(msg)
	}

	return msError{
		parent: err,
		msg:    msg,
		caller: logger.GetCaller(skip),
	}
}

// func Cause(err error) error {
// 	if err == nil {
// 		return err
// 	}

// 	if si, ok := err.(*stackItem); !ok {
// 		return err
// 	} else {
// 		if si.error == nil {
// 			return si
// 		} else {
// 			return Cause(si.error)
// 		}
// 	}

// } // Cause()

// func (w *stackItem) Unwrap() error {
// 	if w == nil {
// 		return nil
// 	}
// 	return w.error
// } // stackItem.Unwrap()
