# Go Diagnostics Binding Architecture

This document describes the architecture used for the Go Diagnostics API in Speech SDK bindings.

## Goals

- Provide feature parity with Java/Python/C# diagnostics surfaces.
- Keep Go API idiomatic while preserving native SDK behavior.
- Isolate CGo/native interop details behind small, focused abstractions.
- Maintain backward compatibility for existing users of the legacy `diagnostics` package.

## Package Layout

### `diagnostics/logging` (new primary API)

This package contains the full diagnostics surface:

- `FileLogger` for file logging
- `MemoryLogger` for memory-ring logging and dump helpers
- `EventLogger` for callback-based log streaming
- `ConsoleLogger` for stdout/stderr logging
- `TraceError`, `TraceWarning`, `TraceInfo`, `TraceVerbose` (+ `WithCaller` variants)
- `Level` enum-like type (`Error`, `Warning`, `Info`, `Verbose`)

Files are intentionally split by logger/type to keep ownership and review scope small.

### `diagnostics` (deprecated compatibility layer)

This package keeps old entry points available and forwards behavior to equivalent native diagnostics APIs.

- Marked with `Deprecated:` comments.
- Intended only for compatibility during migration.
- New work should target `diagnostics/logging`.

## Architectural Decisions

### 1) Singleton logger model

Each logger is a package-level singleton (`FileLogger`, `MemoryLogger`, etc.) matching the native SDK's process-wide semantics and the shape used in other language bindings.

### 2) Thin wrappers around native C APIs

Go methods map ~1:1 to native diagnostics APIs from `speechapi_c_diagnostics.h`, keeping behavior changes centralized in the native layer.

### 3) Explicit CGo boundary ownership

All callback bridge code lives in `logging/cfunctions.go`: C trampoline -> exported Go function -> mutex-guarded dispatch. This avoids glue duplication and keeps thread-safety constraints visible.

### 4) Property bag for file logger configuration

`FileLogger.Start` passes `SPEECH-LogFilename` and `SPEECH-AppendToLogFile` via a temporary property bag, mirroring the native/C++/Java pattern.

### 5) Backward compatibility strategy

Legacy `diagnostics` package stays available with `Deprecated:` markers guiding users to `diagnostics/logging`.

## Error Handling

- Native result codes (`SPXHR`/`AZACHR`) surface as Go `error` values via lightweight wrappers.
- Native `void` APIs remain fire-and-forget. Invalid caller input is validated early where useful.

## Concurrency

- Native event callbacks arrive on SDK worker threads.
- `EventLogger` protects state with a mutex; callback handlers should be fast and non-blocking.

## Quick Start

```go
import "github.com/Microsoft/cognitive-services-speech-sdk-go/diagnostics/logging"

// File logging
logging.FileLogger.Start("/tmp/speech.log")
defer logging.FileLogger.Stop()

// Memory logging with dump
logging.MemoryLogger.Start()
defer logging.MemoryLogger.Stop()
logging.TraceInfo("recognized: %s", result.Text)
lines := logging.MemoryLogger.DumpToSlice()

// Event-based logging
logging.EventLogger.SetCallback(func(msg string) {
    fmt.Println(msg)
})
defer logging.EventLogger.SetCallback(nil)

// Console logging
logging.ConsoleLogger.Start()
defer logging.ConsoleLogger.Stop()

// Set log level
logging.FileLogger.SetLevel(logging.Error)
```

## Testing

Integration tests are gated by `SPEECH_SDK_AVAILABLE=1`; pure unit tests (e.g. `TestLevelString`, `TestLoggingError`) run unconditionally.

## Local Validation

Requirements: Go toolchain, CGo-compatible C compiler, Speech SDK headers and native library.

Key environment variables: `CGO_ENABLED=1`, `CGO_CFLAGS` (header path), `CGO_LDFLAGS` (lib path), `SPEECH_SDK_AVAILABLE=1`.

Run tests via CMake:

```bash
cmake --build build --target go-tests --config Release
cmake --build build --target go-tests-race --config Release
```

## Extension Guidelines

1. Add new capabilities in `diagnostics/logging` first.
2. Keep Go naming idiomatic but parity-aligned with Java/Python/C#.
3. Mirror native signatures closely; add tests for success and error paths.
4. Only add compatibility shims in `diagnostics` when needed for non-breaking upgrades.

## Known Constraints

- Logging is process-wide by native SDK design.
- Event callback is a single registration point per process.
- Integration tests are opt-in via `SPEECH_SDK_AVAILABLE=1`.
