// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

// #include <stdlib.h>
// #include <speechapi_c_diagnostics.h>
import "C"
import (
	"strings"
	"unsafe"
)

type consoleLogger struct{}

// ConsoleLogger is the process-wide console logging singleton.
var ConsoleLogger consoleLogger

// Start begins console logging. Pass true to log to stderr instead of stdout.
func (consoleLogger) Start(logToStderr ...bool) {
	toStderr := false
	if len(logToStderr) > 0 {
		toStderr = logToStderr[0]
	}
	C.diagnostics_log_console_start_logging(C.bool(toStderr))
}

// Stop ends console logging.
func (consoleLogger) Stop() {
	C.diagnostics_log_console_stop_logging()
}

// SetFilters sets case-sensitive filters; call with no args to clear.
func (consoleLogger) SetFilters(filters ...string) {
	joined := strings.Join(filters, ";")
	cFilters := C.CString(joined)
	defer C.free(unsafe.Pointer(cFilters))
	C.diagnostics_log_console_set_filters(cFilters)
}

// SetLevel sets the log level for console logging.
func (consoleLogger) SetLevel(level Level) {
	setLevel("console", level)
}
