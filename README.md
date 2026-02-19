# Overview

This project contains Go bindings for the [Microsoft Cognitive Services Speech SDK](https://docs.microsoft.com/azure/cognitive-services/speech-service/).

# Getting Started

Check the [Speech SDK Setup documentation for Go](https://learn.microsoft.com/azure/ai-services/speech-service/quickstarts/setup-platform?pivots=programming-language-go)

Get started with [speech-to-text sample for Go](https://learn.microsoft.com/azure/ai-services/speech-service/get-started-speech-to-text?pivots=programming-language-go)

Get started with [text-to-speech sample for Go](https://learn.microsoft.com/azure/ai-services/speech-service/get-started-text-to-speech?pivots=programming-language-go)

This project requires Go 1.18 or later.

# Packages

| Package | Description |
|---------|-------------|
| `speech` | Speech recognition, synthesis, translation |
| `audio` | Audio configuration and streams |
| `dialog` | Dialog service connector |
| `common` | Shared types and properties |
| `diagnostics` | Legacy diagnostics (deprecated) |
| `diagnostics/logging` | **Diagnostics logging** — file, memory, event, console loggers and trace helpers |

# Diagnostics Logging

The `diagnostics/logging` package provides process-wide logging for debugging and diagnostics:

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

See [diagnostics/README.md](diagnostics/README.md) for architecture details.

# Reference

Reference documentation for these packages is available at http://aka.ms/csspeech/goref

# Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
