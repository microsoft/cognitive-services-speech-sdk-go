package diagnostics

// #include <stdlib.h>
// #include <speechapi_c_diagnostics.h>
import "C"
import "unsafe"

// StartMemoryLogging starts logging to memory
func StartMemoryLogging() {
	C.diagnostics_log_memory_start_logging()
}

// StopMemoryLogging stops logging to memory
func StopMemoryLogging() {
	C.diagnostics_log_memory_stop_logging()
}

// SetMemoryLogFilters sets filters for memory logging
func SetMemoryLogFilters(filters string) {
	cFilters := C.CString(filters)
	defer C.free(unsafe.Pointer(cFilters))
	C.diagnostics_log_memory_set_filters(cFilters)
}

// GetMemoryLogLineNumOldest gets the line number of the oldest memory log entry
func GetMemoryLogLineNumOldest() uint {
	return uint(C.diagnostics_log_memory_get_line_num_oldest())
}

// GetMemoryLogLineNumNewest gets the line number of the newest memory log entry
func GetMemoryLogLineNumNewest() uint {
	return uint(C.diagnostics_log_memory_get_line_num_newest())
}

// GetMemoryLogLine gets a specific line from the memory log
func GetMemoryLogLine(lineNum uint) string {
	cLine := C.diagnostics_log_memory_get_line(C.size_t(lineNum))
	if cLine == nil {
		return ""
	}
	return C.GoString(cLine)

}

// DumpMemoryLogToStderr dumps the memory log to stderr
func DumpMemoryLogToStderr() error {
	ret := uintptr(C.diagnostics_log_memory_dump_to_stderr())
	if ret != 0 {
		return newDiagnosticsError("dumpMemoryLogToStderr", ret)
	}
	return nil
}

// DumpMemoryLog dumps the memory log to a file and/or standard output
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

// DumpMemoryLogOnExit dumps the memory log when the program exits
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

// StartConsoleLogging starts logging to the console
func StartConsoleLogging(logToStderr bool) {
	C.diagnostics_log_console_start_logging(C.bool(logToStderr))
}

// StopConsoleLogging stops logging to the console
func StopConsoleLogging() {
	C.diagnostics_log_console_stop_logging()
}

// SetConsoleLogFilters sets filters for console logging
func SetConsoleLogFilters(filters string) {
	cFilters := C.CString(filters)
	defer C.free(unsafe.Pointer(cFilters))
	C.diagnostics_log_console_set_filters(cFilters)
}
