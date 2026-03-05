// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

// Level defines the severity level for logging.
type Level int

const (
    // Error severity level.
    Error Level = 0x02

    // Warning severity level.
    Warning Level = 0x04

    // Info severity level.
    Info Level = 0x08

    // Verbose severity level.
    Verbose Level = 0x10
)

// String returns the level name.
func (l Level) String() string {
    switch l {
    case Error:
        return "error"
    case Warning:
        return "warning"
    case Info:
        return "info"
    case Verbose:
        return "verbose"
    default:
        return "unknown"
    }
}
