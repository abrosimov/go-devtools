// Package debugx provides some useful functions for debugging.
package debugx

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// GetCurrentFunc returns the function name and the line number where GetCurrentFunc was called.
func GetCurrentFunc() string {
	const currentFuncPositionInStack = 3
	return getFuncNameInStack(currentFuncPositionInStack)
}

// GetCalleeFunc returns the function name and the line number of the function that called the function in which GetCalleeFunc was called.
func GetCalleeFunc() string {
	const calleeFuncPositionInStack = 4
	return getFuncNameInStack(calleeFuncPositionInStack)
}

func getFuncNameInStack(skip int) string {
	const stackSize = 15
	pc := make([]uintptr, stackSize)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d", frame.Function, frame.Line)
}

// GetStackTrace returns the stack trace of the current goroutine.
func GetStackTrace() string {
	const maxStackSize = 4096
	b := make([]byte, maxStackSize)
	n := runtime.Stack(b, false)
	s := string(b[:n])

	calls := strings.Split(strings.TrimSpace(s), "\n")
	captureFuncArgAddresses := regexp.MustCompile(`\(0x.*?\)`)

	/*
		Input format:
		goroutine 18 [running]:
		github.com/abrosimov/go-devtools/debugx.GetStackTrace()
			/Users/kirillabrosimov/Projects/go-devtools/debugx/runtime_funcs.go:34 +0x50
		github.com/abrosimov/go-devtools/debugx_test.TestGetStackTrace(0x140000b0680)
			/Users/kirillabrosimov/Projects/go-devtools/debugx/runtime_funcs_test.go:57 +0x24

		So we need to put tabbed lines to the previous line with removing unnecessary things.
	*/
	result := make([]string, 0, len(calls))
	lastWrittenIdx := 0
	skipNext := false
	for i := 0; i < len(calls); i++ {
		// skipNext and line comparison with "...debugx.GetStackTrace()" is needed to skip the line with GetStackTrace() call.
		if skipNext {
			skipNext = false
			continue
		}

		line := calls[i]
		if line == "github.com/abrosimov/go-devtools/debugx.GetStackTrace()" {
			skipNext = true
			continue
		}

		if line[0] == '\t' {
			split := strings.Split(strings.TrimSpace(line), " ")
			if len(split) == 0 {
				//nolint:godox // I really not sure that I need use log here, however, might be it'd be better to di it. But structured or not?
				// TODO: log error?
				continue
			}
			result[lastWrittenIdx] = fmt.Sprintf("%s at %s", result[lastWrittenIdx], split[0])
		} else {
			result = append(result, captureFuncArgAddresses.ReplaceAllString(line, "(...)"))
			lastWrittenIdx = len(result) - 1
		}
	}
	return strings.Join(result, "\n")
}
