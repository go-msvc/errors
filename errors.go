package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// The point at which an error was added to the error stack
type caller struct {
	pc       uintptr
	function string
	file     string
	line     int
} // caller

// Item on the error stack
type stackItem struct {
	error
	pos    int
	notTop bool
	msg    string
	caller *caller
} // stackItem

//New is the same as Errorf
func New(message string) error {
	return Errorf(message)
}

// Errorf creates a new error using printf formatting and captures the caller details
func Errorf(format string, args ...interface{}) error {
	return &stackItem{
		msg:    fmt.Sprintf(format, args...),
		caller: getCaller(2),
	}
}

//Wrapf is same as Errorf(), but wraps around an existing error
func Wrapf(err error, format string, args ...interface{}) error {
	return wrapf(err, format, 3, args...)
}

//Wrap is like Wrapf() but without formatting
func Wrap(err error, message string) error {
	return wrapf(err, message, 3)
}

func wrapf(err error, format string, skip int, args ...interface{}) error {
	if err == nil {
		return nil
	}
	pos := 0
	if si, ok := err.(*stackItem); ok {
		pos = si.pos + 1
	}
	return &stackItem{
		pos:    pos,
		error:  err,
		msg:    fmt.Sprintf(format, args...),
		caller: getCaller(skip),
	}
}

//Cause ...
func Cause(err error) error {
	if err == nil {
		return err
	}

	si, ok := err.(*stackItem)
	if !ok {
		return err
	}
	if si.error == nil {
		return si
	}
	return Cause(si.error)
}

func (w *stackItem) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if !w.notTop {
				fmt.Fprintf(s, "%s\n", w.msg)
			}

			io.WriteString(s, " ")
			w.caller.Format(s, verb)
			fmt.Fprintf(s, " [%d] %s\n", w.pos, w.msg)

			if w.error != nil {
				if si, ok := w.error.(*stackItem); ok {
					si.notTop = true
				}
				fmt.Fprintf(s, "%+v", w.error)
			}
			return
		}

		fallthrough
	case 's':
		io.WriteString(s, w.msg)
		if w.error != nil {
			fmt.Fprintf(s, ": %s", w.error)
		}
	case 'q':
		fmt.Fprintf(s, "%q", w.msg)
	}
}

func (w *stackItem) Error() string {
	if w != nil {
		return fmt.Sprintf("%s", w)
	}
	return ""
}

//Format formats the caller according to the fmt.Formatter interface.
func (caller caller) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			fn := caller.function
			const fnw = 16
			const fw = 22
			if len(fn) > fnw {
				fn = fn[len(fn)-fnw:]
			}
			fileName := func() string {
				idx := strings.LastIndexByte(caller.file, '/')
				if idx == -1 {
					return caller.file
				}
				idx = strings.LastIndexByte(caller.file[:idx], '/')
				if idx == -1 {
					return caller.file
				}
				return caller.file[idx+1:]
			}() + ":" + strconv.Itoa(caller.line)
			if len(fileName) > fw {
				fileName = fileName[len(fileName)-fw:]
			}
			fmt.Fprintf(s, "%*s %*s", fw, fileName, fnw, fn)
		default:
			io.WriteString(s, path.Base(caller.file))
		}
	case 'd':
		fmt.Fprintf(s, "%d", caller.line)
	case 'n':
		io.WriteString(s, funcname(caller.function))
	case 'v':
		caller.Format(s, 's')
	}
}

func getCaller(skip int) *caller {
	if pc, file, line, ok := runtime.Caller(skip); ok {
		fnName := ""
		if fn := runtime.FuncForPC(pc); fn != nil {
			fnName = fn.Name()
		} else {
			fnName = "unknown"
		}
		return &caller{
			pc:       pc,
			function: fnName,
			file:     file,
			line:     line,
		}
	}
	return &caller{}
}

//funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
