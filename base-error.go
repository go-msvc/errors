package errors

import (
	"fmt"
)

// all errors types defined in this package embeds this to offer BaseError interface
type BaseError interface {
	error
	fmt.Stringer
	// fmt.Formatter
	Source() Caller
}

type baseError struct {
	source  Caller
	wrapped error //nil when not wrapping
}

// override this for each error type in this package
// this is used recursively in the Error() string method
// which is only defined in baseError
// func (err baseError) String() string {
// 	return "error at " + err.source.String()
// }

// Error() returns String() of this and all wrapped errors
func (err baseError) Error() string {
	if err.wrapped != nil {
		return err.wrapped.Error()
	}
	return "base-error" //err.String()
}

func (err baseError) Unwrap() error {
	return err.wrapped
}

func (err baseError) Source() Caller {
	return err.source
}

// implement fmt.Formatter
// func (err baseError) Format(f fmt.State, c rune) {
// var s string
// s = fmt.Sprintf("base(%p)", err)
// switch c {
// case 'v':
// 	s += fmt.Sprintf("%v:%s", err.source, err.Error())
// default: //also for 's'
// 	s += "...baseNoString..." //err.String()
// }
// io.WriteString(f, s)

// if err.wrapped != nil {
// 	stack := false
// 	if f.Flag('+') {
// 		stack = true
// 		io.WriteString(f, " because ")
// 	}
// 	if f.Flag('-') {
// 		stack = true
// 		io.WriteString(f, "\n")
// 	}
// 	if stack {
// 		if formatter, ok := err.wrapped.(fmt.Formatter); ok {
// 			formatter.Format(f, c)
// 		} else {
// 			io.WriteString(f, err.wrapped.Error())
// 		}
// 	}
// }
// }
