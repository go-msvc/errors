package errors_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/go-msvc/errors/v2"
	"github.com/stretchr/testify/assert"
)

func TestCaller(t *testing.T) {
	lineNr := 14
	c := errors.GetCaller(1)
	t.Logf("Pkg=%s, File=%s, Line=%d, Func=%s", c.Package(), c.File(), c.Line(), c.Function())
	if c.Package() != "github.com/go-msvc/errors/v2_test" {
		t.Fatalf("Package=%s != github.com/go-msvc/errors/v2_test", c.Package())
	}
	if path.Base(c.File()) != "caller_test.go" {
		t.Fatalf("Package=%s != caller_test.go", c.File())
	}
	if c.Line() != lineNr {
		t.Fatalf("Line=%d != %d", c.Line(), lineNr)
	}
	if c.Function() != "TestCaller" {
		t.Fatalf("Function=%s != TestCaller", c.Function())
	}

	//test formatting:
	base := fmt.Sprintf("caller_test.go(%d)", lineNr)
	tests := []struct {
		format        string
		expectedValue string
	}{
		{"%s", base},
		{"%v", "github.com/go-msvc/errors/v2_test/" + base},
		{"%30s", fmt.Sprintf("%30s", base)},   //shorter and padded
		{"%-30s", fmt.Sprintf("%-30s", base)}, //shorter and padded
		{"%10s", fmt.Sprintf("%10s", base)},   //longer than 10
		{"%-10s", fmt.Sprintf("%-10s", base)}, //longer than 10
		{"%10.10s", fmt.Sprintf("%10.10s", base)},
		{"%-10.10s", fmt.Sprintf("%-10.10s", base)},
		{"%30.10s", fmt.Sprintf("%30.10s", base)},
		{"%-30.10s", fmt.Sprintf("%-30.10s", base)},
		{"%10.30s", fmt.Sprintf("%10.30s", base)},
		{"%-10.30s", fmt.Sprintf("%-10.30s", base)},
	}
	for index, test := range tests {
		t.Run(fmt.Sprintf("test[%d]:%+v", index, test), func(t *testing.T) {
			s := fmt.Sprintf(test.format, c)
			assert.Equal(t, test.expectedValue, s)
			//t.Logf("OK: fmt.Sprintf(\"%s\", caller) -> \"%s\"", test.format, s)
		})
	}
}
