// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package diagnostics

import (
    "os"
    "path/filepath"
    "testing"
)

func skipIfNoSDK(t *testing.T) {
    t.Helper()
    if os.Getenv("SPEECH_SDK_AVAILABLE") != "1" {
        t.Skip("Skipping integration test: SPEECH_SDK_AVAILABLE not set")
    }
}

func TestStartStopFileLogging(t *testing.T) {
    skipIfNoSDK(t)

    logPath := filepath.Join(t.TempDir(), "legacy_diagnostics.log")
    if err := StartFileLogging(logPath); err != nil {
        t.Fatalf("StartFileLogging: %v", err)
    }
    if err := StopFileLogging(); err != nil {
        t.Fatalf("StopFileLogging: %v", err)
    }
}

func TestStartStopFileLoggingAppendMode(t *testing.T) {
    skipIfNoSDK(t)

    logPath := filepath.Join(t.TempDir(), "legacy_diagnostics_append.log")
    if err := StartFileLogging(logPath, true); err != nil {
        t.Fatalf("StartFileLogging(append): %v", err)
    }
    if err := StopFileLogging(); err != nil {
        t.Fatalf("StopFileLogging: %v", err)
    }
}

func TestStartConsoleLoggingDefault(t *testing.T) {
    skipIfNoSDK(t)

    StartConsoleLogging()
    StopConsoleLogging()
}
