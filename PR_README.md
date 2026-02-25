# PR: Diagnostics/Logging Package, AutoDetect API Fix & CI Test Pipeline

**Branch:** `feature/diagnostics-logging`  
**Base:** `master`

---

## Summary

This PR delivers three areas of work:

1. **Diagnostics/Logging package** — A new `diagnostics/logging` package providing full-featured logging capabilities (file, memory, event, console) with the legacy `diagnostics` package preserved as a deprecated compatibility layer.
2. **AutoDetect API spelling fix** — Corrects the misspelled `Fom` → `From` in `NewSpeechRecognizerFromAutoDetectSourceLangConfig` and `NewSpeechSynthesizerFromAutoDetectSourceLangConfig`, with backward-compatible deprecated aliases.
3. **CI pipeline enhancement** — Adds automated test execution to the Azure DevOps pipeline using Key Vault secrets with log redaction.

---

## 1. Diagnostics/Logging Package

### New: `diagnostics/logging`

A full diagnostics surface providing feature parity with Java/Python/C# SDK bindings:

| Logger | Description |
|--------|-------------|
| `FileLogger` | Log to a file with append/overwrite modes |
| `MemoryLogger` | In-memory ring buffer with dump helpers (`DumpToSlice`, `DumpToFile`, `DumpToWriter`, `DumpToStderr`) |
| `EventLogger` | Callback-based log streaming for custom integrations |
| `ConsoleLogger` | Direct stdout/stderr logging |

Also includes:
- Trace helpers: `TraceError`, `TraceWarning`, `TraceInfo`, `TraceVerbose` (+ `WithCaller` variants)
- `Level` type with `Error`, `Warning`, `Info`, `Verbose` constants
- Per-logger `SetLevel()` and `SetFilters()` for fine-grained control

### Updated: `diagnostics` (deprecated compatibility layer)

The legacy `diagnostics` package is preserved with `Deprecated:` markers guiding users to `diagnostics/logging`. Existing callers continue to work without changes.

### Files Added/Changed

| File | Change |
|------|--------|
| `diagnostics/logging/doc.go` | Package documentation |
| `diagnostics/logging/level.go` | Level type and constants |
| `diagnostics/logging/error.go` | Error wrapper for native result codes |
| `diagnostics/logging/cfunctions.go` | CGo trampoline for event callbacks |
| `diagnostics/logging/file_logger.go` | FileLogger singleton |
| `diagnostics/logging/memory_logger.go` | MemoryLogger singleton |
| `diagnostics/logging/event_logger.go` | EventLogger singleton |
| `diagnostics/logging/console_logger.go` | ConsoleLogger singleton |
| `diagnostics/logging/spx_trace.go` | Trace helpers |
| `diagnostics/logging/logging_test.go` | 46 tests covering all loggers |
| `diagnostics/diagnostics.go` | Updated with deprecated markers |
| `diagnostics/diagnostics_test.go` | 3 legacy compatibility tests |
| `diagnostics/error.go` | Updated error handling |
| `diagnostics/README.md` | Architecture documentation |

### Test Coverage: 49 tests

| Area | Tests |
|------|-------|
| MemoryLogger | Start/Stop, DumpToFile, DumpToSlice, DumpToStderr, DumpToWriter, DumpOnExit, SetFilters, SetLevel, RoundTrip, error paths |
| EventLogger | Callback registration/replacement/unregister, SetFilters, SetLevel |
| FileLogger | Start/Stop, append/overwrite modes, double stop, SetFilters, SetLevel, path validation |
| ConsoleLogger | Start/Stop, stderr mode, SetFilters, SetLevel |
| Trace | All four levels, WithCaller variants, format args |
| Level | String() for all values |
| Error | Error message format |

---

## 2. AutoDetect API Spelling Fix

### Problem

The public API had a typo — `Fom` instead of `From` — in two constructor functions:
- `NewSpeechRecognizerFomAutoDetectSourceLangConfig`
- `NewSpeechSynthesizerFomAutoDetectSourceLangConfig`

### Solution

- **Corrected functions**: `NewSpeechRecognizerFromAutoDetectSourceLangConfig` and `NewSpeechSynthesizerFromAutoDetectSourceLangConfig` (new primary API)
- **Deprecated aliases**: The old misspelled names (`Fom`) delegate to the corrected functions, so existing callers are unaffected

### Files Changed

| File | Change |
|------|--------|
| `speech/speech_recognizer.go` | Renamed function + added deprecated alias |
| `speech/speech_synthesizer.go` | Renamed function + added deprecated alias |

### Tests

| Test | File | What it verifies |
|------|------|------------------|
| `TestRecognitionWithLanguageAutoDetection` | `speech/speech_recognizer_test.go` | Creates recognizer via `NewSpeechRecognizerFromAutoDetectSourceLangConfig` with `en-US`/`de-DE`, recognizes speech from `turn_on_the_lamp.wav` |
| `TestSynthesisWithLanguageAutoDetection` | `speech/speech_synthesizer_test.go` | Creates synthesizer via `NewSpeechSynthesizerFromAutoDetectSourceLangConfig` with open range, synthesizes Chinese text "你好，世界。", verifies audio > 1s |

Both tests pass against the live Speech service (westcentralus region).

### Samples

| Sample | File |
|--------|------|
| `RecognizeOnceFromAutoDetectSourceLangConfig` | `samples/recognizer/auto_detect.go` |
| `SynthesisFromAutoDetectSourceLangConfig` | `samples/synthesizer/auto_detect.go` |

Both samples are wired into `samples/main.go`.

---

## 3. CI Pipeline Enhancement

### Problem

The Go pipeline only built the code — it did not run any tests. Test credentials were not available in CI.

### Solution

Added three-file infrastructure following the JavaScript SDK's Key Vault pattern:

#### `ci/generate-subscription-file.yml` (new)
Reusable ADO template step that:
1. Downloads the `CarbonSubscriptionsJson` secret from the `CarbonSDK-CICD` Key Vault via `AzureKeyVault@2`
2. Writes it to `secrets/test.subscriptions.regions.json` via `file-creator@6`

#### `ci/load-build-secrets.sh` (new)
Bash script sourced at test time that:
1. Reads the subscriptions JSON using `jq`
2. Exports `SPEECH_SUBSCRIPTION_KEY` and `SPEECH_SUBSCRIPTION_REGION`
3. Provides a `global_redact` function (perl-based streaming filter) that replaces subscription keys with `***` in CI output
4. Validates that required fields are present and non-null

#### `ci/azure-pipelines.yml` (modified)
Updated test step that sources secrets, sets CGo flags, and pipes `go test` output through `global_redact`.

#### `.gitignore` (modified)
Added ignore rules to prevent accidental commit of secrets:
- `secrets/` — CI-generated directory containing subscription JSON
- `test.subscriptions.regions.json` — local subscription keys
- `test.certificates.json` — local certificates

### Security

Follows the same pattern as the [Carbon pipeline](https://dev.azure.com/speedme/_git/Carbon?path=/ci/pipeline/scripts/load-build-secrets.sh):

- `set +x` disables trace mode before any secret handling to prevent bash from echoing variable assignments
- `source ci/load-build-secrets.sh` loads secrets, defines `global_redact`, and redacts both raw and URL-encoded key variants
- `go test ... 2>&1 | global_redact` pipes all test output through a perl-based streaming filter that replaces subscription keys with `***`
- Test files redact secrets from memory log dumps at the source via `redactSecrets()` before `t.Log()`
- Secrets are fetched from Key Vault at runtime, not stored in the repo
- Service connection: `ADO -> Speech Services - DEV - SDK`

---

## Full Diff Summary

- 28 files changed, 2251 insertions(+), 176 deletions(-)

| Category | Files |
|----------|-------|
| Diagnostics/Logging | 14 files (new package + legacy updates + README) |
| AutoDetect API fix | 2 source files, 2 test files, 2 sample files, 1 main.go |
| CI Pipeline | 3 files (1 modified, 2 new) |
| Repo hygiene | 1 file (`.gitignore` updated) |

## Validation

| Check | Result |
|-------|--------|
| `go build ./audio ./common ./dialog ./speech` | PASS |
| `TestRecognitionWithLanguageAutoDetection` | PASS (live service) |
| `TestSynthesisWithLanguageAutoDetection` | PASS (live service, 1837ms audio) |
| YAML syntax validation | PASS (both pipeline files) |
| Handle leak check | ZERO leaked handles |
| Backward compatibility (deprecated aliases) | Preserved — old names delegate to new |
