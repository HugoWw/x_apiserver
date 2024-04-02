package util

import (
	"runtime"
	"strings"
)

// GetCallerName Gets the name of the method specified in the stack to call
// The argument skip is the number of stack frames to skip before recording,
// with 0 identifying the frame for Callers itself and 1 identifying the caller of Callers.
// The argument excludes is to filter out unwanted stack calls
func GetCallerName(skip int, excludes []string) string {
	var fn string

	pc := make([]uintptr, 20)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc)
OUTER:
	for i := 0; i < n; i++ {
		frame, more := frames.Next()
		// fmt.Printf("********  %s\n", frame.Function)
		if !more {
			break
		}
		fpath := frame.Function
		for _, exclude := range excludes {
			if strings.Contains(fpath, exclude) {
				continue OUTER
			}
		}
		slash := strings.LastIndex(fpath, "/")
		if slash == -1 {
			fn = fpath
		} else {
			fn = fpath[slash+1:]
		}
		return fn
	}

	return fn
}
