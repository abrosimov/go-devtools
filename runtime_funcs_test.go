package devtools_test

import (
	"strings"
	"testing"

	"github.com/abrosimov/go-devtools"
)

func TestGetCalleeFunc(t *testing.T) {
	if devtools.GetCurrentFunc() != "github.com/abrosimov/go-devtools_test.TestGetCalleeFunc:11" {
		t.Errorf("current func is not TestGetCalleeFunc:11")
	}
	func() {
		if devtools.GetCurrentFunc() != "github.com/abrosimov/go-devtools_test.TestGetCalleeFunc.func1:15" {
			t.Errorf("current func is not TestGetCalleeFunc.func1:15")
		}
	}()
	func() {
		if devtools.GetCurrentFunc() != "github.com/abrosimov/go-devtools_test.TestGetCalleeFunc.func2:20" {
			t.Errorf("current func is not TestGetCalleeFunc.func2:20")
		}
	}()
}

func TestGetCallee(t *testing.T) {
	callee := devtools.GetCalleeFunc()
	// Because we want to get name of the function that called TestGetCallee,
	// and we don't want to stick to line numbers in "testing" framework
	// we're sticking only to the function name.
	parts := strings.Split(callee, ":")
	if parts[0] != "testing.tRunner" {
		t.Errorf("calle of the test function is not testing.tRunner")
	}

	func() {
		if callee = devtools.GetCalleeFunc(); callee != "github.com/abrosimov/go-devtools_test.TestGetCallee:40" {
			t.Errorf("caller of the test function is not TestGetCallee:40")
		}
	}()

	func() {
		if callee = devtools.GetCalleeFunc(); callee != "github.com/abrosimov/go-devtools_test.TestGetCallee:46" {
			t.Errorf("caller of the test function is not TestGetCallee:46")
		}
	}()
}
