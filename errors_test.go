package errors_test

import (
	"fmt"
	"testing"

	"github.com/go-msvc/errors"
)

func TestError(t *testing.T) {
	lineNr := 12
	err := errors.Errorf("cannot open")

	//test formatting:
	tests := []struct {
		format        string
		expectedValue string
	}{
		{"%s", fmt.Sprintf("cannot open")},

		//default expected behaviour: %v -> show file name reference
		{"%v", fmt.Sprintf("errors_test.go(%d):cannot open", lineNr)},
		{"%V", fmt.Sprintf("github.com/go-msvc/errors_test/errors_test.go(%d):cannot open", lineNr)},

		//function name reference
		{"%f", fmt.Sprintf("TestError(%d):cannot open", lineNr)},
		{"%F", fmt.Sprintf("github.com/go-msvc/errors_test.TestError(%d):cannot open", lineNr)},
	}
	for index, test := range tests {
		var s string
		s = fmt.Sprintf(test.format, err)
		if s != test.expectedValue {
			t.Fatalf("test[%d] fmt.Sprintf(\"%s\", caller) -> \"%s\" != \"%s\"", index, test.format, s, test.expectedValue)
		}
		t.Logf("test[%d] OK: fmt.Sprintf(\"%s\", caller) -> \"%s\"", index, test.format, s)
	}
}

func get() error {
	return errors.Errorf("cannot open")
}

func TestStack(t *testing.T) {
	getLineNr := 40 //line when get() fails
	lineNr := 47    //line where we call Wrapf()
	if err := get(); err != nil {
		err = errors.Wrapf(err, "cannot get")
		//test formatting:
		tests := []struct {
			format        string
			expectedValue string
		}{
			{"%s", fmt.Sprintf("cannot get")},
			{"%+s", fmt.Sprintf("cannot get because cannot open")}, //+ for full stack

			//default expected behaviour: %v -> show file name reference
			{"%v", fmt.Sprintf("errors_test.go(%d):cannot get", lineNr)},
			{"%V", fmt.Sprintf("github.com/go-msvc/errors_test/errors_test.go(%d):cannot get", lineNr)},
			//+ to get full stack in one line
			{"%+v", fmt.Sprintf("errors_test.go(%d):cannot get because errors_test.go(%d):cannot open", lineNr, getLineNr)},
			{"%+V", fmt.Sprintf("github.com/go-msvc/errors_test/errors_test.go(%d):cannot get because github.com/go-msvc/errors_test/errors_test.go(%d):cannot open", lineNr, getLineNr)},
			//- to get full stack in multiple lines
			{"%-v", fmt.Sprintf("errors_test.go(%d):cannot get\nerrors_test.go(%d):cannot open", lineNr, getLineNr)},
			{"%-V", fmt.Sprintf("github.com/go-msvc/errors_test/errors_test.go(%d):cannot get\ngithub.com/go-msvc/errors_test/errors_test.go(%d):cannot open", lineNr, getLineNr)},

			//function name reference
			{"%f", fmt.Sprintf("TestStack(%d):cannot get", lineNr)},
			{"%F", fmt.Sprintf("github.com/go-msvc/errors_test.TestStack(%d):cannot get", lineNr)},
			//+ to get full stack in one line
			{"%+f", fmt.Sprintf("TestStack(%d):cannot get because get(%d):cannot open", lineNr, getLineNr)},
			{"%+F", fmt.Sprintf("github.com/go-msvc/errors_test.TestStack(%d):cannot get because github.com/go-msvc/errors_test.get(%d):cannot open", lineNr, getLineNr)},
			//- to get full stack in multiple lines
			{"%-f", fmt.Sprintf("TestStack(%d):cannot get\nget(%d):cannot open", lineNr, getLineNr)},
			{"%-F", fmt.Sprintf("github.com/go-msvc/errors_test.TestStack(%d):cannot get\ngithub.com/go-msvc/errors_test.get(%d):cannot open", lineNr, getLineNr)},
		}
		for index, test := range tests {
			var s string
			s = fmt.Sprintf(test.format, err)
			if s != test.expectedValue {
				t.Fatalf("test[%d] fmt.Sprintf(\"%s\", caller) -> \"%s\" != \"%s\"", index, test.format, s, test.expectedValue)
			}
			t.Logf("test[%d] OK: fmt.Sprintf(\"%s\", caller) -> \"%s\"", index, test.format, s)
		}
	}
}
