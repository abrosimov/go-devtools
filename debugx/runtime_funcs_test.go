package debugx_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/abrosimov/go-devtools/debugx"
)

func TestGetCurrentFunc(t *testing.T) {
	tests := []struct {
		name         string
		getFunc      func() string
		wantContains string
	}{
		{
			name: "from test function",
			//nolint:gocritic // lambda needed to test stack depth
			getFunc: func() string {
				return debugx.GetCurrentFunc()
			},
			wantContains: "TestGetCurrentFunc",
		},
		{
			name: "from nested function",
			//nolint:gocritic // lambda needed to test stack depth
			getFunc: func() string {
				//nolint:gocritic // nested lambda tests different stack depth
				nested := func() string {
					return debugx.GetCurrentFunc()
				}
				return nested()
			},
			wantContains: "TestGetCurrentFunc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.getFunc()
			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("GetCurrentFunc() = %q, want to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestGetCalleeFunc(t *testing.T) {
	t.Run("called from test context", func(t *testing.T) {
		// When called directly from test, caller should be from testing package
		got := debugx.GetCalleeFunc()
		if !strings.Contains(got, "testing.") {
			t.Errorf("GetCalleeFunc() = %q, want to contain %q", got, "testing.")
		}
	})

	t.Run("called from nested function", func(t *testing.T) {
		//nolint:gocritic // lambda needed to test caller detection
		helper := func() string {
			return debugx.GetCalleeFunc()
		}
		got := helper()
		if !strings.Contains(got, "TestGetCalleeFunc") {
			t.Errorf("GetCalleeFunc() = %q, want to contain %q", got, "TestGetCalleeFunc")
		}
	})
}

func TestGetStackTrace(t *testing.T) {
	tests := []struct {
		name         string
		wantContains []string
	}{
		{
			name: "stack trace contains goroutine and test function",
			wantContains: []string{
				"goroutine",
				"TestGetStackTrace",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := debugx.GetStackTrace()
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("GetStackTrace() missing %q in stack trace", want)
				}
			}
		})
	}
}

// ExampleGetCurrentFunc demonstrates how to get the current function name and line number.
func ExampleGetCurrentFunc() {
	funcInfo := debugx.GetCurrentFunc()
	// The output will contain the function name and line number
	// Format: package.function:line
	fmt.Println(strings.Contains(funcInfo, "ExampleGetCurrentFunc"))
	// Output: true
}

// ExampleGetCalleeFunc demonstrates how to get the caller function information.
func ExampleGetCalleeFunc() {
	//nolint:gocritic // lambda needed to demonstrate caller detection
	helperFunction := func() string {
		return debugx.GetCalleeFunc()
	}

	caller := helperFunction()
	// The caller is ExampleGetCalleeFunc
	fmt.Println(strings.Contains(caller, "ExampleGetCalleeFunc"))
	// Output: true
}

// ExampleGetStackTrace demonstrates how to get a formatted stack trace.
func ExampleGetStackTrace() {
	trace := debugx.GetStackTrace()
	// Stack trace contains goroutine information and function calls
	fmt.Println(strings.Contains(trace, "goroutine"))
	fmt.Println(strings.Contains(trace, "ExampleGetStackTrace"))
	// Output:
	// true
	// true
}
