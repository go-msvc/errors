package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
)

type Caller interface {
	fmt.Stringer
	Package() string
	PackageFile() string
	Function() string
	File() string
	Line() int
}

type caller struct {
	file       string
	line       int
	pkgDotFunc string
}

func GetCaller(skip int) Caller {
	c := caller{
		file:       "",
		line:       -1,
		pkgDotFunc: "",
	}

	var pc uintptr
	var ok bool
	if pc, c.file, c.line, ok = runtime.Caller(skip); !ok {
		return c
	}

	if fn := runtime.FuncForPC(pc); fn != nil {
		c.pkgDotFunc = fn.Name()
	}
	return c
} //GetCaller()

func (c caller) String() string {
	return fmt.Sprintf("%s(%d)", path.Base(c.file), c.line)
}

// with Function: "github.com/go-msvc/ms_test.TestCaller"
// return "github.com/go-msvc/ms_test"
func (c caller) Package() string {
	if i := strings.LastIndex(c.pkgDotFunc, "."); i >= 0 {
		return c.pkgDotFunc[:i]
	}
	return ""
}

// return "github.com/go-msvc/ms_test/my_test.go"
func (c caller) PackageFile() string {
	if i := strings.LastIndex(c.pkgDotFunc, "."); i >= 0 {
		return c.pkgDotFunc[:i] + "/" + path.Base(c.file)
	}
	return ""
}

// with Function: "github.com/go-msvc/ms_test.TestCaller"
// return "github.com/go-msvc/ms_test"
func (c caller) Function() string {
	if i := strings.LastIndex(c.pkgDotFunc, "."); i >= 0 {
		return c.pkgDotFunc[i+1:]
	}
	return ""
}

// return full file name on system where code is built...
func (c caller) File() string {
	return c.file
}

func (c caller) Line() int {
	return c.line
}

// %s -> basefile(line)
// %v -> fullpath(line)
// %#:#s -> min and max len, align right
// %-#:#s -> min and max len, align left
func (caller caller) Format(f fmt.State, c rune) {
	var s string
	switch c {
	case 'v': //full name
		s = fmt.Sprintf("%s(%d)", caller.PackageFile(), caller.line)
	default: //base name
		s = fmt.Sprintf("%s(%d)", path.Base(caller.file), caller.line)
	}

	l := len(s)
	if maxLen, ok := f.Precision(); ok && l > maxLen {
		s = s[:maxLen]
		l = maxLen
	}
	if f.Flag('-') {
		//left align
		if minLen, ok := f.Width(); ok && len(s) < minLen {
			s += spaces(minLen - l)
		}
	} else {
		//right align
		l := len(s)
		if minLen, ok := f.Width(); ok && len(s) < minLen {
			s = spaces(minLen-l) + s
		}
	}
	io.WriteString(f, s)
} // caller.Format()

func spaces(n int) string {
	s := "                                                    "
	for len(s) < n {
		s += s
	}
	return s[:n]
}
