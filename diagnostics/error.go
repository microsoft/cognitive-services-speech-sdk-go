package diagnostics

// #include <azac_error.h>
// #include <azac_api_c_common.h>
import "C"

import "fmt"

type diagnosticsError struct {
	operation string
	code      uintptr
}

func newDiagnosticsError(operation string, code uintptr) error {
	return &diagnosticsError{
		operation: operation,
		code:      code,
	}
}

func (e *diagnosticsError) Error() string {
	return fmt.Sprintf("diagnostics operation '%s' failed with error code %d", e.operation, e.code)
}
