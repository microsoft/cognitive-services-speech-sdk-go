// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

// #include <stdlib.h>
// #include <speechapi_c_diagnostics.h>
//
// uintptr_t cgo_register_event_callback(void);
// uintptr_t cgo_unregister_event_callback(void);
import "C"
import (
    "strings"
    "sync"
    "unsafe"
)

// EventCallback is the function signature for event logger callbacks.
type EventCallback func(message string)

var (
    callbackMu     sync.Mutex
    storedCallback EventCallback
)

type eventLogger struct{}

// EventLogger is the process-wide event logging singleton.
var EventLogger eventLogger

//export goEventLoggerCallback
func goEventLoggerCallback(logLine *C.char) {
    callbackMu.Lock()
    cb := storedCallback
    callbackMu.Unlock()
    if cb != nil {
        cb(C.GoString(logLine))
    }
}

// SetCallback registers or removes the event logging handler (nil to unregister).
func (eventLogger) SetCallback(callback EventCallback) error {
    callbackMu.Lock()
    defer callbackMu.Unlock()

    if callback != nil {
        storedCallback = callback
        ret := uintptr(C.cgo_register_event_callback())
        if ret != 0 {
            storedCallback = nil
            return newLoggingError("EventLogger.SetCallback", ret)
        }
    } else {
        ret := uintptr(C.cgo_unregister_event_callback())
        storedCallback = nil
        if ret != 0 {
            return newLoggingError("EventLogger.SetCallback(nil)", ret)
        }
    }
    return nil
}

// SetFilters sets case-sensitive filters; call with no args to clear.
func (eventLogger) SetFilters(filters ...string) error {
    joined := strings.Join(filters, ";")
    cFilters := C.CString(joined)
    defer C.free(unsafe.Pointer(cFilters))
    ret := uintptr(C.diagnostics_logmessage_set_filters(cFilters))
    if ret != 0 {
        return newLoggingError("EventLogger.SetFilters", ret)
    }
    return nil
}

// SetLevel sets the log level for event logging.
func (eventLogger) SetLevel(level Level) {
    cLogger := C.CString("event")
    defer C.free(unsafe.Pointer(cLogger))
    cLevel := C.CString(level.String())
    defer C.free(unsafe.Pointer(cLevel))
    C.diagnostics_set_log_level(cLogger, cLevel)
}
