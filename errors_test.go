package errors

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const linkBc = " because "
const linkNl = "\n"

func TestNew(t *testing.T) {
	const msg = "some-error"
	expErrLine := GetCaller(1).Line() + 1     //+1 because err is defined in the next line of this test
	errs := []BaseError{New(msg), Error(msg)} //both functions does exactly the same
	for idx, err := range errs {
		t.Run(fmt.Sprintf("test[%d]", idx), func(t *testing.T) {
			t.Logf("err=(%T)%+v", err, err)

			assert.Equal(t, err.Source().Line(), expErrLine)
			assert.Equal(t, msg, err.String())
			assert.Equal(t, msg, err.Error()) //recursive err.String() == err.Error() because this err does not wrap another

			//error without source
			assert.Equal(t, "error is "+err.String(), fmt.Sprintf("error is %s", err)) //add "error is " to avoid compiler warning
			assert.Equal(t, err.String(), fmt.Sprintf("%+s", err))
			assert.Equal(t, err.String(), fmt.Sprintf("%-s", err))

			//error with source as basename
			assert.Equal(t, fmt.Sprintf("errors_test.go(%d):%s", expErrLine, err.String()), fmt.Sprintf("%v", err))
			assert.Equal(t, fmt.Sprintf("errors_test.go(%d):%s", expErrLine, err.String()), fmt.Sprintf("%+v", err))
			assert.Equal(t, fmt.Sprintf("errors_test.go(%d):%s", expErrLine, err.String()), fmt.Sprintf("%-v", err))

			//error with source as fullpath
			assert.Equal(t, fmt.Sprintf("github.com/go-msvc/errors/v2/errors_test.go(%d):%s", expErrLine, err.String()), fmt.Sprintf("%V", err))
			assert.Equal(t, fmt.Sprintf("github.com/go-msvc/errors/v2/errors_test.go(%d):%s", expErrLine, err.String()), fmt.Sprintf("%+V", err))
			assert.Equal(t, fmt.Sprintf("github.com/go-msvc/errors/v2/errors_test.go(%d):%s", expErrLine, err.String()), fmt.Sprintf("%-V", err))
		})
	}
}

func TestWrapExternalError(t *testing.T) {
	//an external error
	_, origErr := os.Open("/some/file/that/does/not/exist")

	//wrap it here
	const msg = "cannot open file"
	expErrLine := GetCaller(1).Line() + 1 //+1 because err is defined in the next line of this test
	err := Wrap(origErr, msg)

	assert.Equal(t, err.Source().Line(), expErrLine)
	assert.Equal(t, msg, err.String())
	assert.Equal(t, msg+linkBc+origErr.Error(), err.Error())

	//'s' -> error without source
	assert.Equal(t, "xx"+err.String(), fmt.Sprintf("xx%s", err))                  //no recursion
	assert.Equal(t, err.String()+linkBc+origErr.Error(), fmt.Sprintf("%+s", err)) //+ recurse on wrapped errors with because
	assert.Equal(t, err.String()+linkNl+origErr.Error(), fmt.Sprintf("%-s", err)) //+ recurse on wrapped errors with newline

	//'v' -> error with source as basename
	errWithSource := fmt.Sprintf("errors_test.go(%d):%s", expErrLine, err.String())
	assert.Equal(t, errWithSource, fmt.Sprintf("%v", err))                         //no recursion
	assert.Equal(t, errWithSource+linkBc+origErr.Error(), fmt.Sprintf("%+v", err)) //+ recurse on wrapped errors with because
	assert.Equal(t, errWithSource+linkNl+origErr.Error(), fmt.Sprintf("%-v", err)) //+ recurse on wrapped errors with newline

	//'V' -> error with source as fullpath
	errWithSource = fmt.Sprintf("github.com/go-msvc/errors/v2/errors_test.go(%d):%s", expErrLine, err.String())
	assert.Equal(t, errWithSource, fmt.Sprintf("%V", err))                         //no recursion
	assert.Equal(t, errWithSource+linkBc+origErr.Error(), fmt.Sprintf("%+V", err)) //+ recurse on wrapped errors with because
	assert.Equal(t, errWithSource+linkNl+origErr.Error(), fmt.Sprintf("%-V", err)) //+ recurse on wrapped errors with newline
}

func TestWrappedAndFormatted(t *testing.T) {
	const id = "444332"
	expMsg1 := fmt.Sprintf("not found(%s)", id)
	expErrLine1 := GetCaller(1).Line() + 1 //+1 because err is defined in the next line of this test
	origErr := Errorf("not found(%s)", id)

	assert.Equal(t, origErr.Source().Line(), expErrLine1)
	assert.Equal(t, expMsg1, origErr.String())
	assert.Equal(t, expMsg1, origErr.Error())

	const value2 = 50
	expMsg2 := fmt.Sprintf("invalid size(%d)", value2)
	expErrLine2 := GetCaller(1).Line() + 1 //+1 because err is defined in the next line of this test
	err := Wrapf(origErr, "invalid size(%d)", value2)

	assert.Equal(t, err.Source().Line(), expErrLine2)
	assert.Equal(t, expMsg2, err.String())
	assert.Equal(t, expMsg2+linkBc+expMsg1, err.Error())

	//'s' -> error without source
	assert.Equal(t, "yy"+err.String(), fmt.Sprintf("yy%s", err))                  //no recursion
	assert.Equal(t, err.String()+linkBc+origErr.Error(), fmt.Sprintf("%+s", err)) //+ recurse on wrapped errors with because
	assert.Equal(t, err.String()+linkNl+origErr.Error(), fmt.Sprintf("%-s", err)) //+ recurse on wrapped errors with newline

	//'v' -> error with source as basename
	errOrgWithSource := fmt.Sprintf("errors_test.go(%d):%s", expErrLine1, origErr.String())
	errWithSource := fmt.Sprintf("errors_test.go(%d):%s", expErrLine2, err.String())
	assert.Equal(t, errWithSource, fmt.Sprintf("%v", err))                          //no recursion
	assert.Equal(t, errWithSource+linkBc+errOrgWithSource, fmt.Sprintf("%+v", err)) //+ recurse on wrapped errors with because
	assert.Equal(t, errWithSource+linkNl+errOrgWithSource, fmt.Sprintf("%-v", err)) //+ recurse on wrapped errors with newline

	//'V' -> error with source as fullpath
	errOrgWithSource = fmt.Sprintf("github.com/go-msvc/errors/v2/errors_test.go(%d):%s", expErrLine1, origErr.String())
	errWithSource = fmt.Sprintf("github.com/go-msvc/errors/v2/errors_test.go(%d):%s", expErrLine2, err.String())
	assert.Equal(t, errWithSource, fmt.Sprintf("%V", err))                          //no recursion
	assert.Equal(t, errWithSource+linkBc+errOrgWithSource, fmt.Sprintf("%+V", err)) //+ recurse on wrapped errors with because
	assert.Equal(t, errWithSource+linkNl+errOrgWithSource, fmt.Sprintf("%-V", err)) //+ recurse on wrapped errors with newline
}

// func TestErrorIs(t *testing.T) {
// 	//define two type of errors
// 	e1 := Error("error1")
// 	e2 := Error("error2")

// 	//something fail using a type of error, and then get passed up to wrap again
// 	err := Wrapf(e1, "failed some")
// 	err = Wrapf(err, "failed more")

// 	//confirm it is e1 and not e2
// 	if errors.Is(err, e1) {
// 		t.Logf("Is e1, good")
// 	} else {
// 		t.Fatalf("Is not e1, bad")
// 	}

// 	if !errors.Is(err, e2) {
// 		t.Logf("Is not e2, good")
// 	} else {
// 		t.Fatalf("Is e1, bad")
// 	}
// }
