package devtools

import (
	"fmt"
	"runtime"
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
	return s
}
