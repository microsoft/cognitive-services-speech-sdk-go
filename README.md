# Overview

This project contains Golang binding for the Microsoft Cognitive Service Speech SDK.

# Getting Started

Check the [Speech SDK Setup documentation for Go](https://docs.microsoft.com/azure/cognitive-services/speech-service/quickstarts/setup-platform?tabs=dotnet%2Cwindows%2Cjre%2Cbrowser&pivots=programming-language-go)

Get started with [speech-to-text sample for Go](https://docs.microsoft.com/azure/cognitive-services/speech-service/get-started-speech-to-text?tabs=windowsinstall&pivots=programming-language-go)

Get started with [text-to-speech sample for Go](https://docs.microsoft.com/azure/cognitive-services/speech-service/get-started-text-to-speech?tabs=script%2Cwindowsinstall&pivots=programming-language-go)

This project requires Go 1.13

# Features

## Language Detection

The Speech SDK supports automatic language detection for speech recognition. You can specify a list of candidate languages, and the SDK will detect which language is being spoken.

### Basic Usage

```go
import (
    "fmt"
    "github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
    "github.com/Microsoft/cognitive-services-speech-sdk-go/common"
    "github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// Create speech config
config, _ := speech.NewSpeechConfigFromSubscription("YourKey", "YourRegion")
defer config.Close()

// Create auto-detect config with candidate languages
autoDetectConfig, _ := speech.NewAutoDetectSourceLanguageConfigFromLanguages(
    []string{"en-US", "de-DE", "es-MX"})
defer autoDetectConfig.Close()

// Create recognizer with auto-detect
recognizer, _ := speech.NewSpeechRecognizerFomAutoDetectSourceLangConfig(
    config, autoDetectConfig, nil)
defer recognizer.Close()

// Recognize
result, _ := recognizer.RecognizeOnce()
defer result.Close()

// Get detected language using the helper class
autoDetectResult := speech.NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(result)
fmt.Printf("Detected: %s\n", autoDetectResult.Language)
fmt.Printf("Text: %s\n", result.Text)
```

### Language Detection Modes

- **At-Start** (default): Detects language at the beginning of audio (up to 10 candidate languages)
- **Continuous**: Detects language changes throughout audio (up to 4 candidate languages)

```go
// Enable Continuous mode for code-switching scenarios
config.SetProperty(common.SpeechServiceConnectionLanguageIDMode, "Continuous")

autoDetectConfig, _ := speech.NewAutoDetectSourceLanguageConfigFromLanguages(
    []string{"en-US", "de-DE", "es-MX", "ja-JP"})  // Up to 4 for Continuous mode

recognizer, _ := speech.NewSpeechRecognizerFomAutoDetectSourceLangConfig(
    config, autoDetectConfig, nil)

// In your event handler:
recognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
    autoDetectResult := speech.NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(&event.Result)
    fmt.Printf("Language: %s, Text: %s\n", autoDetectResult.Language, event.Result.Text)
})
```

### Examples

- [Recognize Once with Language Detection](samples/recognizer/language_detection_recognize_once.go)
- [Continuous Recognition with Continuous LID](samples/recognizer/language_detection_continuous_lid.go)

For more details, see [Language Detection Quick Reference](LANGUAGE_DETECTION_QUICK_REFERENCE.md).

# Reference

Reference documentation for these packages is available at http://aka.ms/csspeech/goref

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
