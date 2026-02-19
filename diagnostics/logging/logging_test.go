// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package logging

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLevelString(t *testing.T) {
    tests := []struct {
        level Level
        want  string
    }{
        {Error, "error"},
        {Warning, "warning"},
        {Info, "info"},
        {Verbose, "verbose"},
        {Level(0xFF), "unknown"},
    }
    for _, tc := range tests {
        if got := tc.level.String(); got != tc.want {
            t.Errorf("Level(%d).String() = %q, want %q", int(tc.level), got, tc.want)
        }
    }
}

func TestLoggingError(t *testing.T) {
    err := newLoggingError("TestOp", 0x15)
    if err == nil {
        t.Fatal("expected non-nil error")
    }
    s := err.Error()
    if !strings.Contains(s, "TestOp") {
        t.Errorf("error string should contain operation name, got: %s", s)
    }
    if !strings.Contains(s, "0x15") {
        t.Errorf("error string should contain hex error code, got: %s", s)
    }
}

func TestFileLoggerStartEmptyPath(t *testing.T) {
    skipIfNoSDK(t)

    err := FileLogger.Start("")
    if err == nil {
        t.Fatal("expected error for empty file path")
    }
}

func TestFileLoggerStartWhitespacePath(t *testing.T) {
    skipIfNoSDK(t)

    err := FileLogger.Start("   ")
    if err == nil {
        t.Fatal("expected error for whitespace-only file path")
    }
}

func TestFileLoggerStartBadDirectory(t *testing.T) {
    skipIfNoSDK(t)

    err := FileLogger.Start(filepath.Join("Z:\\nonexistent_dir_12345", "test.log"))
    if err == nil {
        t.Fatal("expected error for non-existent directory")
    }
}

func TestFileLoggerSetLevel(t *testing.T) {
    skipIfNoSDK(t)

    tmpDir := t.TempDir()
    logPath := filepath.Join(tmpDir, "sdk_level_test.log")
    if err := FileLogger.Start(logPath); err != nil {
        t.Fatalf("Start: %v", err)
    }
    defer FileLogger.Stop()

    FileLogger.SetLevel(Error)
    FileLogger.SetLevel(Verbose)
}

func TestConsoleLoggerSetFilters(t *testing.T) {
    skipIfNoSDK(t)

    ConsoleLogger.Start()
    defer ConsoleLogger.Stop()

    ConsoleLogger.SetFilters("filterA", "filterB")
    ConsoleLogger.SetFilters()
}

func TestConsoleLoggerSetLevel(t *testing.T) {
    skipIfNoSDK(t)

    ConsoleLogger.Start()
    defer ConsoleLogger.Stop()

    ConsoleLogger.SetLevel(Warning)
    ConsoleLogger.SetLevel(Verbose)
}

func TestEventLoggerSetFilters(t *testing.T) {
    skipIfNoSDK(t)

    err := EventLogger.SetCallback(func(msg string) {})
    if err != nil {
        t.Fatalf("SetCallback: %v", err)
    }
    defer EventLogger.SetCallback(nil)

    if err := EventLogger.SetFilters("filterA"); err != nil {
        t.Fatalf("SetFilters: %v", err)
    }
    if err := EventLogger.SetFilters(); err != nil {
        t.Fatalf("SetFilters clear: %v", err)
    }
}

func TestEventLoggerSetLevel(t *testing.T) {
    skipIfNoSDK(t)

    err := EventLogger.SetCallback(func(msg string) {})
    if err != nil {
        t.Fatalf("SetCallback: %v", err)
    }
    defer EventLogger.SetCallback(nil)

    EventLogger.SetLevel(Error)
    EventLogger.SetLevel(Verbose)
}

func TestMemoryLoggerDumpToWriter(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceInfo("test dump to writer %d", 42)
    time.Sleep(50 * time.Millisecond)

    var buf strings.Builder
    err := MemoryLogger.DumpToWriter(&buf)
    if err != nil {
        t.Fatalf("DumpToWriter: %v", err)
    }
    _ = buf.String()
}

func TestMemoryLoggerDumpEmptyPath(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    err := MemoryLogger.Dump("")
    if err == nil {
        t.Fatal("expected error for empty dump path")
    }
}

func TestMemoryLoggerDumpWhitespacePath(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    err := MemoryLogger.Dump("   ")
    if err == nil {
        t.Fatal("expected error for whitespace-only dump path")
    }
}

func TestMemoryLoggerDumpBadDirectory(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    err := MemoryLogger.Dump(filepath.Join("Z:\\nonexistent_dir_99999", "dump.log"))
    if err == nil {
        t.Fatal("expected error for non-existent dump directory")
    }
}

func TestMemoryLoggerDumpToWriterNil(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    err := MemoryLogger.DumpToWriter(nil)
    if err == nil {
        t.Fatal("expected error for nil writer")
    }
}

func TestMemoryLoggerDumpOnExitBadDirectory(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    err := MemoryLogger.DumpOnExit(filepath.Join("Z:\\nonexistent_dir_99999", "dump.log"), "", false, false)
    if err == nil {
        t.Fatal("expected error for non-existent DumpOnExit directory")
    }
}

// skipIfNoSDK skips if SPEECH_SDK_AVAILABLE != "1".
func skipIfNoSDK(t *testing.T) {
    t.Helper()
    if os.Getenv("SPEECH_SDK_AVAILABLE") != "1" {
        t.Skip("Skipping integration test: SPEECH_SDK_AVAILABLE not set")
    }
}

func TestMemoryLoggerRoundTrip(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    time.Sleep(50 * time.Millisecond)

    lines := MemoryLogger.DumpToSlice()
    _ = lines
}

func TestMemoryLoggerSetFilters(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    MemoryLogger.SetFilters("filterA", "filterB")

    MemoryLogger.SetFilters()
}

func TestMemoryLoggerSetLevel(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    MemoryLogger.SetLevel(Warning)
    MemoryLogger.SetLevel(Verbose)
}

func TestFileLoggerStartStop(t *testing.T) {
    skipIfNoSDK(t)

    tmpDir := t.TempDir()
    logPath := filepath.Join(tmpDir, "sdk_test.log")

    if err := FileLogger.Start(logPath); err != nil {
        t.Fatalf("FileLogger.Start returned error: %v", err)
    }

    time.Sleep(50 * time.Millisecond)

    if err := FileLogger.Stop(); err != nil {
        t.Fatalf("FileLogger.Stop returned error: %v", err)
    }

    if _, err := os.Stat(logPath); os.IsNotExist(err) {
        t.Error("expected log file to be created")
    }
}

func TestFileLoggerStartAppendMode(t *testing.T) {
    skipIfNoSDK(t)

    tmpDir := t.TempDir()
    logPath := filepath.Join(tmpDir, "sdk_append_test.log")

    if err := os.WriteFile(logPath, []byte("MARKER\n"), 0644); err != nil {
        t.Fatalf("WriteFile: %v", err)
    }

    if err := FileLogger.Start(logPath, true); err != nil {
        t.Fatalf("Start(append) returned error: %v", err)
    }
    time.Sleep(50 * time.Millisecond)
    if err := FileLogger.Stop(); err != nil {
        t.Fatalf("Stop returned error: %v", err)
    }

    data, err := os.ReadFile(logPath)
    if err != nil {
        t.Fatalf("ReadFile: %v", err)
    }
    if !strings.HasPrefix(string(data), "MARKER") {
        t.Error("expected file to start with MARKER (append mode)")
    }
}

func TestFileLoggerSetFilters(t *testing.T) {
    skipIfNoSDK(t)

    tmpDir := t.TempDir()
    logPath := filepath.Join(tmpDir, "sdk_filter_test.log")

    if err := FileLogger.Start(logPath); err != nil {
        t.Fatalf("FileLogger.Start returned error: %v", err)
    }
    defer FileLogger.Stop()

    if err := FileLogger.SetFilters("filterA"); err != nil {
        t.Fatalf("SetFilters returned error: %v", err)
    }

    // Clear
    if err := FileLogger.SetFilters(); err != nil {
        t.Fatalf("SetFilters() clear returned error: %v", err)
    }
}

func TestConsoleLoggerStartStop(t *testing.T) {
    skipIfNoSDK(t)

    ConsoleLogger.Start()
    ConsoleLogger.Stop()
}

func TestConsoleLoggerStartToStderr(t *testing.T) {
    skipIfNoSDK(t)

    ConsoleLogger.Start(true)
    ConsoleLogger.Stop()
}

func TestEventLoggerCallback(t *testing.T) {
    skipIfNoSDK(t)

    var (
        mu       sync.Mutex
        received []string
    )

    err := EventLogger.SetCallback(func(msg string) {
        mu.Lock()
        defer mu.Unlock()
        received = append(received, msg)
    })
    if err != nil {
        t.Fatalf("SetCallback returned error: %v", err)
    }

    time.Sleep(100 * time.Millisecond)

    if err := EventLogger.SetCallback(nil); err != nil {
        t.Fatalf("SetCallback(nil) returned error: %v", err)
    }

    mu.Lock()
    defer mu.Unlock()
    t.Logf("received %d event log lines", len(received))
}

func TestSpxTraceNoSDKPanic(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceError("test error %d", 1)
    TraceWarning("test warning %d", 2)
    TraceInfo("test info %d", 3)
    TraceVerbose("test verbose %d", 4)

    TraceErrorWithCaller("fake.go", 42, "explicit error")
    TraceInfoWithCaller("fake.go", 100, "explicit info")
}

func TestMemoryLoggerDumpToFile(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceInfo("memory dump file test %d", 99)
    time.Sleep(50 * time.Millisecond)

    dumpPath := filepath.Join(t.TempDir(), "memory_dump.log")
    if err := MemoryLogger.Dump(dumpPath); err != nil {
        t.Fatalf("Dump: %v", err)
    }

    if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
        t.Error("expected dump file to be created")
    }
}

func TestMemoryLoggerDumpToSliceContent(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceInfo("slice_marker_%d", 7777)
    time.Sleep(50 * time.Millisecond)

    lines := MemoryLogger.DumpToSlice()
    found := false
    for _, line := range lines {
        if strings.Contains(line, "slice_marker_7777") {
            found = true
            break
        }
    }
    if !found {
        t.Errorf("expected DumpToSlice to contain trace marker, got %d lines", len(lines))
    }
}

func TestMemoryLoggerDumpToStderr(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceInfo("stderr dump test")
    time.Sleep(50 * time.Millisecond)

    if err := MemoryLogger.DumpToStderr(); err != nil {
        t.Fatalf("DumpToStderr: %v", err)
    }
}

func TestMemoryLoggerDumpOnExitValid(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    dumpPath := filepath.Join(t.TempDir(), "on_exit_dump.log")
    if err := MemoryLogger.DumpOnExit(dumpPath, "[PREFIX] ", false, false); err != nil {
        t.Fatalf("DumpOnExit: %v", err)
    }
}

func TestMemoryLoggerDumpOnExitStderrOnly(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    if err := MemoryLogger.DumpOnExit("", "", false, true); err != nil {
        t.Fatalf("DumpOnExit(stderr): %v", err)
    }
}

func TestMemoryLoggerStartStopIdempotent(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    MemoryLogger.Start()
    MemoryLogger.Stop()
    MemoryLogger.Stop()
}

func TestMemoryLoggerDumpToWriterContent(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceInfo("writer_marker_%d", 8888)
    time.Sleep(50 * time.Millisecond)

    var buf strings.Builder
    if err := MemoryLogger.DumpToWriter(&buf); err != nil {
        t.Fatalf("DumpToWriter: %v", err)
    }

    if !strings.Contains(buf.String(), "writer_marker_8888") {
        t.Errorf("expected DumpToWriter output to contain trace marker, got %d bytes", buf.Len())
    }
}

type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) {
    return 0, errors.New("intentional write failure")
}

func TestMemoryLoggerDumpToWriterError(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceInfo("trigger content for writer error test")
    time.Sleep(50 * time.Millisecond)

    err := MemoryLogger.DumpToWriter(failWriter{})
    if err == nil {
        lines := MemoryLogger.DumpToSlice()
        if len(lines) > 0 {
            t.Fatal("expected error from failing writer")
        }
    }
}

func TestMemoryLoggerDumpToWriterBytes(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceInfo("bytes_buffer_test_%d", 5555)
    time.Sleep(50 * time.Millisecond)

    var buf bytes.Buffer
    if err := MemoryLogger.DumpToWriter(&buf); err != nil {
        t.Fatalf("DumpToWriter: %v", err)
    }
    if buf.Len() == 0 {
        lines := MemoryLogger.DumpToSlice()
        if len(lines) > 0 {
            t.Error("expected non-empty buffer when memory log has content")
        }
    }
}

func TestMemoryLoggerSetFiltersMultiple(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    MemoryLogger.SetFilters("alpha", "beta", "gamma")
    MemoryLogger.SetFilters("single")
    MemoryLogger.SetFilters()
}

func TestMemoryLoggerSetLevelAll(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    MemoryLogger.SetLevel(Error)
    MemoryLogger.SetLevel(Warning)
    MemoryLogger.SetLevel(Info)
    MemoryLogger.SetLevel(Verbose)
}

func TestEventLoggerCallbackReceivesTrace(t *testing.T) {
    skipIfNoSDK(t)

    var (
        mu       sync.Mutex
        received []string
    )

    err := EventLogger.SetCallback(func(msg string) {
        mu.Lock()
        defer mu.Unlock()
        received = append(received, msg)
    })
    if err != nil {
        t.Fatalf("SetCallback: %v", err)
    }

    TraceInfo("event_logger_marker_%d", 3333)
    time.Sleep(150 * time.Millisecond)

    if err := EventLogger.SetCallback(nil); err != nil {
        t.Fatalf("SetCallback(nil): %v", err)
    }

    mu.Lock()
    defer mu.Unlock()
    t.Logf("received %d event log lines", len(received))
}

func TestEventLoggerCallbackReplacement(t *testing.T) {
    skipIfNoSDK(t)

    var count1, count2 int
    var mu sync.Mutex

    err := EventLogger.SetCallback(func(msg string) {
        mu.Lock()
        count1++
        mu.Unlock()
    })
    if err != nil {
        t.Fatalf("SetCallback #1: %v", err)
    }

    time.Sleep(50 * time.Millisecond)

    err = EventLogger.SetCallback(func(msg string) {
        mu.Lock()
        count2++
        mu.Unlock()
    })
    if err != nil {
        t.Fatalf("SetCallback #2: %v", err)
    }

    TraceInfo("after replacement %d", 1)
    time.Sleep(100 * time.Millisecond)

    if err := EventLogger.SetCallback(nil); err != nil {
        t.Fatalf("SetCallback(nil): %v", err)
    }

    mu.Lock()
    defer mu.Unlock()
    t.Logf("callback1=%d, callback2=%d", count1, count2)
}

func TestEventLoggerUnregisterWithoutRegister(t *testing.T) {
    skipIfNoSDK(t)

    if err := EventLogger.SetCallback(nil); err != nil {
        t.Fatalf("SetCallback(nil) without prior register: %v", err)
    }
}

func TestEventLoggerSetFiltersMultiple(t *testing.T) {
    skipIfNoSDK(t)

    err := EventLogger.SetCallback(func(msg string) {})
    if err != nil {
        t.Fatalf("SetCallback: %v", err)
    }
    defer EventLogger.SetCallback(nil)

    if err := EventLogger.SetFilters("alpha", "beta", "gamma"); err != nil {
        t.Fatalf("SetFilters multiple: %v", err)
    }
    if err := EventLogger.SetFilters(); err != nil {
        t.Fatalf("SetFilters clear: %v", err)
    }
}

func TestEventLoggerSetLevelAll(t *testing.T) {
    skipIfNoSDK(t)

    err := EventLogger.SetCallback(func(msg string) {})
    if err != nil {
        t.Fatalf("SetCallback: %v", err)
    }
    defer EventLogger.SetCallback(nil)

    EventLogger.SetLevel(Error)
    EventLogger.SetLevel(Warning)
    EventLogger.SetLevel(Info)
    EventLogger.SetLevel(Verbose)
}

// --- Additional File Logger tests ---

func TestFileLoggerStartOverwrite(t *testing.T) {
    skipIfNoSDK(t)

    tmpDir := t.TempDir()
    logPath := filepath.Join(tmpDir, "overwrite_test.log")

    if err := os.WriteFile(logPath, []byte("OLD_CONTENT\n"), 0644); err != nil {
        t.Fatalf("WriteFile: %v", err)
    }

    if err := FileLogger.Start(logPath); err != nil {
        t.Fatalf("Start: %v", err)
    }
    time.Sleep(50 * time.Millisecond)
    if err := FileLogger.Stop(); err != nil {
        t.Fatalf("Stop: %v", err)
    }

    data, err := os.ReadFile(logPath)
    if err != nil {
        t.Fatalf("ReadFile: %v", err)
    }
    if strings.HasPrefix(string(data), "OLD_CONTENT") {
        t.Error("expected file to be overwritten, but old content still present")
    }
}

func TestFileLoggerDoubleStop(t *testing.T) {
    skipIfNoSDK(t)

    logPath := filepath.Join(t.TempDir(), "double_stop.log")
    if err := FileLogger.Start(logPath); err != nil {
        t.Fatalf("Start: %v", err)
    }
    if err := FileLogger.Stop(); err != nil {
        t.Fatalf("Stop #1: %v", err)
    }
    FileLogger.Stop()
}

func TestFileLoggerSetLevelAll(t *testing.T) {
    skipIfNoSDK(t)

    logPath := filepath.Join(t.TempDir(), "level_all.log")
    if err := FileLogger.Start(logPath); err != nil {
        t.Fatalf("Start: %v", err)
    }
    defer FileLogger.Stop()

    FileLogger.SetLevel(Error)
    FileLogger.SetLevel(Warning)
    FileLogger.SetLevel(Info)
    FileLogger.SetLevel(Verbose)
}

func TestTraceWithCallerAllLevels(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceErrorWithCaller("test.go", 10, "error %s", "msg")
    TraceWarningWithCaller("test.go", 20, "warning %s", "msg")
    TraceInfoWithCaller("test.go", 30, "info %s", "msg")
    TraceVerboseWithCaller("test.go", 40, "verbose %s", "msg")
}

func TestTraceFormatArgs(t *testing.T) {
    skipIfNoSDK(t)

    MemoryLogger.Start()
    defer MemoryLogger.Stop()

    TraceInfo("int=%d float=%.2f str=%s", 42, 3.14, "hello")
    TraceError("no args")
    TraceWarning("single %v", true)
    time.Sleep(50 * time.Millisecond)

    lines := MemoryLogger.DumpToSlice()
    t.Logf("trace format test: %d lines captured", len(lines))
}
