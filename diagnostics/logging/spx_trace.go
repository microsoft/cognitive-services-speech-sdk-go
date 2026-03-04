// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

// #include <stdlib.h>
// #include <speechapi_c_diagnostics.h>
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

const (
	traceTitleInfo    = "SPX_TRACE_INFO"
	traceTitleWarning = "SPX_TRACE_WARNING"
	traceTitleError   = "SPX_TRACE_ERROR"
	traceTitleVerbose = "SPX_TRACE_VERBOSE"
)

func traceLevel(l Level) C.int {
	if l != Error && l != Warning && l != Info && l != Verbose {
		return C.int(Verbose)
	}
	return C.int(l)
}

func traceMessage(level Level, title string, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	emitTrace(level, title, file, line, format, args...)
}

func emitTrace(level Level, title string, file string, line int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)

	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	cMsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cMsg))

	C.diagnostics_log_trace_string(traceLevel(level), cTitle, cFile, C.int(line), cMsg)
}

// setLevel sets the log level for the named logger.
func setLevel(loggerName string, level Level) {
	cLogger := C.CString(loggerName)
	defer C.free(unsafe.Pointer(cLogger))
	cLevel := C.CString(level.String())
	defer C.free(unsafe.Pointer(cLevel))
	C.diagnostics_set_log_level(cLogger, cLevel)
}

// TraceError emits an error-level trace message.
func TraceError(format string, args ...interface{}) {
	traceMessage(Error, traceTitleError, format, args...)
}

// TraceWarning emits a warning-level trace message.
func TraceWarning(format string, args ...interface{}) {
	traceMessage(Warning, traceTitleWarning, format, args...)
}

// TraceInfo emits an info-level trace message.
func TraceInfo(format string, args ...interface{}) {
	traceMessage(Info, traceTitleInfo, format, args...)
}

// TraceVerbose emits a verbose-level trace message.
func TraceVerbose(format string, args ...interface{}) {
	traceMessage(Verbose, traceTitleVerbose, format, args...)
}

// TraceErrorWithCaller emits an error-level trace with explicit caller.
func TraceErrorWithCaller(file string, line int, format string, args ...interface{}) {
	emitTrace(Error, traceTitleError, file, line, format, args...)
}

// TraceWarningWithCaller emits a warning-level trace with explicit caller.
func TraceWarningWithCaller(file string, line int, format string, args ...interface{}) {
	emitTrace(Warning, traceTitleWarning, file, line, format, args...)
}

// TraceInfoWithCaller emits an info-level trace with explicit caller.
func TraceInfoWithCaller(file string, line int, format string, args ...interface{}) {
	emitTrace(Info, traceTitleInfo, file, line, format, args...)
}

// TraceVerboseWithCaller emits a verbose-level trace with explicit caller.
func TraceVerboseWithCaller(file string, line int, format string, args ...interface{}) {
	emitTrace(Verbose, traceTitleVerbose, file, line, format, args...)
}
