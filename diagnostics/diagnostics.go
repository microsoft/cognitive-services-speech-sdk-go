// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

// Deprecated: Use the diagnostics/logging sub-package instead.
package diagnostics

// #include <stdlib.h>
// #include <speechapi_c_property_bag.h>
// #include <speechapi_c_diagnostics.h>
import "C"
import "unsafe"

// Deprecated: Use logging.FileLogger.Start() instead.
func StartFileLogging(filename string, appendMode ...bool) error {
    appendToFile := false
    if len(appendMode) > 0 {
        appendToFile = appendMode[0]
    }

    var handle C.SPXHANDLE
    ret := uintptr(C.property_bag_create(&handle))
    if ret != 0 {
        return newDiagnosticsError("startFileLogging: property_bag_create", ret)
    }
    defer C.property_bag_release(handle)

    cFilename := C.CString(filename)
    defer C.free(unsafe.Pointer(cFilename))
    cFileKey := C.CString("SPEECH-LogFilename")
    defer C.free(unsafe.Pointer(cFileKey))
    ret = uintptr(C.property_bag_set_string(handle, -1, cFileKey, cFilename))
    if ret != 0 {
        return newDiagnosticsError("startFileLogging: set filename", ret)
    }

    appendStr := "0"
    if appendToFile {
        appendStr = "1"
    }
    cAppend := C.CString(appendStr)
    defer C.free(unsafe.Pointer(cAppend))
    cAppendKey := C.CString("SPEECH-AppendToLogFile")
    defer C.free(unsafe.Pointer(cAppendKey))
    ret = uintptr(C.property_bag_set_string(handle, -1, cAppendKey, cAppend))
    if ret != 0 {
        return newDiagnosticsError("startFileLogging: set append", ret)
    }

    ret = uintptr(C.diagnostics_log_start_logging(handle, nil))
    if ret != 0 {
        return newDiagnosticsError("startFileLogging", ret)
    }
    return nil
}

// Deprecated: Use logging.FileLogger.Stop() instead.
func StopFileLogging() error {
    ret := uintptr(C.diagnostics_log_stop_logging())
    if ret != 0 {
        return newDiagnosticsError("stopFileLogging", ret)
    }
    return nil
}

// Deprecated: Use logging.MemoryLogger.Start() instead.
func StartMemoryLogging() {
    C.diagnostics_log_memory_start_logging()
}

// Deprecated: Use logging.MemoryLogger.Stop() instead.
func StopMemoryLogging() {
    C.diagnostics_log_memory_stop_logging()
}

// Deprecated: Use logging.MemoryLogger.SetFilters() instead.
func SetMemoryLogFilters(filters string) {
    cFilters := C.CString(filters)
    defer C.free(unsafe.Pointer(cFilters))
    C.diagnostics_log_memory_set_filters(cFilters)
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
    ret := uintptr(C.diagnostics_log_memory_dump_to_stderr())
    if ret != 0 {
        return newDiagnosticsError("dumpMemoryLogToStderr", ret)
    }
    return nil
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
    ret := uintptr(C.diagnostics_log_memory_dump_on_exit(cFilename, cLinePrefix, C.bool(emitToStdOut), C.bool(emitToStdErr)))
    if ret != 0 {
        return newDiagnosticsError("dumpMemoryLogOnExit", ret)
    }
    return nil
}

// Deprecated: Use logging.ConsoleLogger.Start() instead.
func StartConsoleLogging(logToStderr ...bool) {
    toStderr := false
    if len(logToStderr) > 0 {
        toStderr = logToStderr[0]
    }
    C.diagnostics_log_console_start_logging(C.bool(toStderr))
}

// Deprecated: Use logging.ConsoleLogger.Stop() instead.
func StopConsoleLogging() {
    C.diagnostics_log_console_stop_logging()
}

// Deprecated: Use logging.ConsoleLogger.SetFilters() instead.
func SetConsoleLogFilters(filters string) {
    cFilters := C.CString(filters)
    defer C.free(unsafe.Pointer(cFilters))
    C.diagnostics_log_console_set_filters(cFilters)
}
