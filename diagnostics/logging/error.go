// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

import "fmt"

type loggingError struct {
    operation string
    code      uintptr
}

func newLoggingError(operation string, code uintptr) error {
    return &loggingError{
        operation: operation,
        code:      code,
    }
}

func (e *loggingError) Error() string {
    return fmt.Sprintf("diagnostics logging operation '%s' failed with error code 0x%x", e.operation, e.code)
}
