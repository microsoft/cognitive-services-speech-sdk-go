// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package diagnostics

// #include <azac_error.h>
// #include <azac_api_c_common.h>
//
// static const char* diagnostics_get_error_message(uintptr_t code) {
//     return error_get_message((AZAC_HANDLE)code);
// }
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
    msg := C.GoString(C.diagnostics_get_error_message(C.uintptr_t(e.code)))
    return fmt.Sprintf("diagnostics operation '%s' failed with error code 0x%x (%s)", e.operation, e.code, msg)
}
