package debugx_test

import (
	"strings"
	"testing"

	"github.com/abrosimov/go-devtools/debugx"
)

func TestGetCalleeFunc(t *testing.T) {
	currentFunc := debugx.GetCurrentFunc()
	if currentFunc != "github.com/abrosimov/go-devtools/debugx_test.TestGetCalleeFunc:11" {
		t.Errorf("current func is not TestGetCalleeFunc:11, but %q given", currentFunc)
	}
	func() {
		currentFunc = debugx.GetCurrentFunc()
		if currentFunc != "github.com/abrosimov/go-devtools/debugx_test.TestGetCalleeFunc.func1:16" {
			t.Errorf("current func is not TestGetCalleeFunc.func1:15, but %q given", currentFunc)
		}
	}()
	func() {
		currentFunc = debugx.GetCurrentFunc()
		if currentFunc != "github.com/abrosimov/go-devtools/debugx_test.TestGetCalleeFunc.func2:22" {
			t.Errorf("current func is not TestGetCalleeFunc.func2:20, but %q given", currentFunc)
		}
	}()
}

func TestGetCallee(t *testing.T) {
	callee := debugx.GetCalleeFunc()
	// Because we want to get name of the function that called TestGetCallee,
	// and we don't want to stick to line numbers in "testing" framework
	// we're sticking only to the function name.
	parts := strings.Split(callee, ":")
	if parts[0] != "testing.tRunner" {
		t.Errorf("calle of the test function is not testing.tRunner")
	}

	func() {
		callee = debugx.GetCalleeFunc()
		if callee != "github.com/abrosimov/go-devtools/debugx_test.TestGetCallee:44" {
			t.Errorf("caller of the test function is not TestGetCallee:44, but %q given", callee)
		}
	}()

	func() {
		callee = debugx.GetCalleeFunc()
		if callee != "github.com/abrosimov/go-devtools/debugx_test.TestGetCallee:51" {
			t.Errorf("caller of the test function is not TestGetCallee:51, but %q given", callee)
		}
	}()
}

func TestGetStackTrace(t *testing.T) {
	stackTrace := debugx.GetStackTrace()
	if !strings.Contains(stackTrace, "github.com/abrosimov/go-devtools/debugx_test.TestGetStackTrace") {
		t.Errorf("Can't find TestGetStackTrace in the stack trace")
	}
}
