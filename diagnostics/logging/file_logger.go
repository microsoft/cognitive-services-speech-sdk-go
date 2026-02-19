// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_property_bag.h>
// #include <speechapi_c_diagnostics.h>
import "C"
import (
    "os"
    "path/filepath"
    "strings"
    "unsafe"
)

type fileLogger struct{}

// FileLogger is the process-wide file logging singleton.
var FileLogger fileLogger

// Start begins logging to filePath. Pass true for appendMode to append.
func (fileLogger) Start(filePath string, appendMode ...bool) error {
    if strings.TrimSpace(filePath) == "" {
        return newLoggingError("FileLogger.Start", 0x005) // SPXERR_INVALID_ARG
    }

    directory := filepath.Dir(filePath)
    if directory != "" && directory != "." {
        if info, err := os.Stat(directory); err != nil || !info.IsDir() {
            return newLoggingError("FileLogger.Start: invalid directory", 0x005)
        }
    }

    appendToFile := false
    if len(appendMode) > 0 {
        appendToFile = appendMode[0]
    }

    var handle C.SPXHANDLE
    ret := uintptr(C.property_bag_create(&handle))
    if ret != 0 {
        return newLoggingError("FileLogger.Start: property_bag_create", ret)
    }
    defer C.property_bag_release(handle)

    cFilePath := C.CString(filePath)
    defer C.free(unsafe.Pointer(cFilePath))
    cKey := C.CString("SPEECH-LogFilename")
    defer C.free(unsafe.Pointer(cKey))
    ret = uintptr(C.property_bag_set_string(handle, -1, cKey, cFilePath))
    if ret != 0 {
        return newLoggingError("FileLogger.Start: set filename", ret)
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
        return newLoggingError("FileLogger.Start: set append", ret)
    }

    ret = uintptr(C.diagnostics_log_start_logging(handle, nil))
    if ret != 0 {
        return newLoggingError("FileLogger.Start", ret)
    }
    return nil
}

// Stop ends file logging.
func (fileLogger) Stop() error {
    ret := uintptr(C.diagnostics_log_stop_logging())
    if ret != 0 {
        return newLoggingError("FileLogger.Stop", ret)
    }
    return nil
}

// SetFilters sets case-sensitive filters; call with no args to clear.
func (fileLogger) SetFilters(filters ...string) error {
    var handle C.SPXHANDLE
    ret := uintptr(C.property_bag_create(&handle))
    if ret != 0 {
        return newLoggingError("FileLogger.SetFilters: property_bag_create", ret)
    }
    defer C.property_bag_release(handle)

    joined := strings.Join(filters, ";")
    cFilters := C.CString(joined)
    defer C.free(unsafe.Pointer(cFilters))
    cKey := C.CString("SPEECH-LogFileFilters")
    defer C.free(unsafe.Pointer(cKey))
    ret = uintptr(C.property_bag_set_string(handle, -1, cKey, cFilters))
    if ret != 0 {
        return newLoggingError("FileLogger.SetFilters: set filters", ret)
    }

    ret = uintptr(C.diagnostics_log_apply_properties(handle, nil))
    if ret != 0 {
        return newLoggingError("FileLogger.SetFilters", ret)
    }
    return nil
}

// SetLevel sets the log level for file logging.
func (fileLogger) SetLevel(level Level) {
    cLogger := C.CString("file")
    defer C.free(unsafe.Pointer(cLogger))
    cLevel := C.CString(level.String())
    defer C.free(unsafe.Pointer(cLevel))
    C.diagnostics_set_log_level(cLogger, cLevel)
}
