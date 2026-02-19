// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

// #include <stdlib.h>
// #include <speechapi_c_diagnostics.h>
import "C"
import (
    "errors"
    "io"
    "os"
    "path/filepath"
    "strings"
    "unsafe"
)

type memoryLogger struct{}

// MemoryLogger is the process-wide memory logging singleton.
var MemoryLogger memoryLogger

// Start begins memory logging.
func (memoryLogger) Start() {
    C.diagnostics_log_memory_start_logging()
}

// Stop ends memory logging.
func (memoryLogger) Stop() {
    C.diagnostics_log_memory_stop_logging()
}

// SetFilters sets case-sensitive filters; call with no args to clear.
func (memoryLogger) SetFilters(filters ...string) {
    joined := strings.Join(filters, ";")
    cFilters := C.CString(joined)
    defer C.free(unsafe.Pointer(cFilters))
    C.diagnostics_log_memory_set_filters(cFilters)
}

// SetLevel sets the log level for memory logging.
func (memoryLogger) SetLevel(level Level) {
    cLogger := C.CString("memory")
    defer C.free(unsafe.Pointer(cLogger))
    cLevel := C.CString(level.String())
    defer C.free(unsafe.Pointer(cLevel))
    C.diagnostics_set_log_level(cLogger, cLevel)
}

// Dump writes the buffer contents to filePath.
func (memoryLogger) Dump(filePath string) error {
    if strings.TrimSpace(filePath) == "" {
        return newLoggingError("MemoryLogger.Dump", 0x005) // SPXERR_INVALID_ARG
    }

    dir := filepath.Dir(filePath)
    if dir != "" && dir != "." {
        if info, err := os.Stat(dir); err != nil || !info.IsDir() {
            return newLoggingError("MemoryLogger.Dump: invalid directory", 0x005)
        }
    }

    cFilePath := C.CString(filePath)
    defer C.free(unsafe.Pointer(cFilePath))
    cPrefix := C.CString("")
    defer C.free(unsafe.Pointer(cPrefix))
    ret := uintptr(C.diagnostics_log_memory_dump(cFilePath, cPrefix, C.bool(false), C.bool(false)))
    if ret != 0 {
        return newLoggingError("MemoryLogger.Dump", ret)
    }
    return nil
}

// DumpToWriter writes the buffer contents to w.
func (memoryLogger) DumpToWriter(w io.Writer) error {
    if w == nil {
        return errors.New("MemoryLogger.DumpToWriter: writer must not be nil")
    }
    start := uint(C.diagnostics_log_memory_get_line_num_oldest())
    stop := uint(C.diagnostics_log_memory_get_line_num_newest())
    for i := start; i < stop; i++ {
        cLine := C.diagnostics_log_memory_get_line(C.size_t(i))
        if cLine == nil {
            continue
        }
        line := C.GoString(cLine)
        _, err := io.WriteString(w, line)
        if err != nil {
            return err
        }
    }
    return nil
}

// DumpToSlice returns the buffer contents as a string slice.
func (memoryLogger) DumpToSlice() []string {
    start := uint(C.diagnostics_log_memory_get_line_num_oldest())
    stop := uint(C.diagnostics_log_memory_get_line_num_newest())
    output := make([]string, 0, stop-start)
    for i := start; i < stop; i++ {
        cLine := C.diagnostics_log_memory_get_line(C.size_t(i))
        if cLine == nil {
            continue
        }
        output = append(output, C.GoString(cLine))
    }
    return output
}

// DumpToStderr writes the buffer to stderr.
func (memoryLogger) DumpToStderr() error {
    ret := uintptr(C.diagnostics_log_memory_dump_to_stderr())
    if ret != 0 {
        return newLoggingError("MemoryLogger.DumpToStderr", ret)
    }
    return nil
}

// DumpOnExit schedules a buffer dump on process exit.
func (memoryLogger) DumpOnExit(filePath string, linePrefix string, emitToStdOut bool, emitToStdErr bool) error {
    if filePath != "" {
        dir := filepath.Dir(filePath)
        if dir != "" && dir != "." {
            if info, err := os.Stat(dir); err != nil || !info.IsDir() {
                return newLoggingError("MemoryLogger.DumpOnExit: invalid directory", 0x005)
            }
        }
    }

    var cFilePath *C.char
    if filePath != "" {
        cFilePath = C.CString(filePath)
        defer C.free(unsafe.Pointer(cFilePath))
    }
    var cLinePrefix *C.char
    if linePrefix != "" {
        cLinePrefix = C.CString(linePrefix)
        defer C.free(unsafe.Pointer(cLinePrefix))
    }
    ret := uintptr(C.diagnostics_log_memory_dump_on_exit(cFilePath, cLinePrefix, C.bool(emitToStdOut), C.bool(emitToStdErr)))
    if ret != 0 {
        return newLoggingError("MemoryLogger.DumpOnExit", ret)
    }
    return nil
}
