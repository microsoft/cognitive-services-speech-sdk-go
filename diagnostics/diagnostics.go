// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

// Deprecated: Use the diagnostics/logging sub-package instead.
package diagnostics

// #include <stdlib.h>
// #include <speechapi_c_diagnostics.h>
import "C"
import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/diagnostics/logging"
)

// Deprecated: Use logging.FileLogger.Start() instead.
func StartFileLogging(filename string, appendMode ...bool) error {
	return logging.FileLogger.Start(filename, appendMode...)
}

// Deprecated: Use logging.FileLogger.Stop() instead.
func StopFileLogging() error {
	return logging.FileLogger.Stop()
}

// Deprecated: Use logging.MemoryLogger.Start() instead.
func StartMemoryLogging() {
	logging.MemoryLogger.Start()
}

// Deprecated: Use logging.MemoryLogger.Stop() instead.
func StopMemoryLogging() {
	logging.MemoryLogger.Stop()
}

// Deprecated: Use logging.MemoryLogger.SetFilters() instead.
func SetMemoryLogFilters(filters string) {
	logging.MemoryLogger.SetFilters(filters)
}

// Deprecated: Use logging.MemoryLogger.DumpToSlice() instead.
func GetMemoryLogLineNumOldest() uint {
	return uint(C.diagnostics_log_memory_get_line_num_oldest())
}

// Deprecated: Use logging.MemoryLogger.DumpToSlice() instead.
func GetMemoryLogLineNumNewest() uint {
	return uint(C.diagnostics_log_memory_get_line_num_newest())
}

// Deprecated: Use logging.MemoryLogger.DumpToSlice() instead.
func GetMemoryLogLine(lineNum uint) string {
	cLine := C.diagnostics_log_memory_get_line(C.size_t(lineNum))
	if cLine == nil {
		return ""
	}
	return C.GoString(cLine)
}

// Deprecated: Use logging.MemoryLogger.DumpToStderr() instead.
func DumpMemoryLogToStderr() error {
	return logging.MemoryLogger.DumpToStderr()
}

// Deprecated: Use logging.MemoryLogger.Dump() instead.
func DumpMemoryLog(filename string, linePrefix string, emitToStdOut bool, emitToStdErr bool) error {
	var cFilename *C.char
	if filename != "" {
		cFilename = C.CString(filename)
		defer C.free(unsafe.Pointer(cFilename))
	}
	var cLinePrefix *C.char
	if linePrefix != "" {
		cLinePrefix = C.CString(linePrefix)
		defer C.free(unsafe.Pointer(cLinePrefix))
	}
	ret := uintptr(C.diagnostics_log_memory_dump(cFilename, cLinePrefix, C.bool(emitToStdOut), C.bool(emitToStdErr)))
	if ret != 0 {
		return newDiagnosticsError("dumpMemoryLog", ret)
	}
	return nil
}

// Deprecated: Use logging.MemoryLogger.DumpOnExit() instead.
func DumpMemoryLogOnExit(filename string, linePrefix string, emitToStdOut bool, emitToStdErr bool) error {
	return logging.MemoryLogger.DumpOnExit(filename, linePrefix, emitToStdOut, emitToStdErr)
}

// Deprecated: Use logging.ConsoleLogger.Start() instead.
func StartConsoleLogging(logToStderr ...bool) {
	logging.ConsoleLogger.Start(logToStderr...)
}

// Deprecated: Use logging.ConsoleLogger.Stop() instead.
func StopConsoleLogging() {
	logging.ConsoleLogger.Stop()
}

// Deprecated: Use logging.ConsoleLogger.SetFilters() instead.
func SetConsoleLogFilters(filters string) {
	logging.ConsoleLogger.SetFilters(filters)
}
