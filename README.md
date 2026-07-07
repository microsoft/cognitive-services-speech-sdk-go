# Overview

This project contains Golang binding for the Microsoft Cognitive Service Speech SDK.

# Getting Started

Check the [Speech SDK Setup documentation for Go](https://docs.microsoft.com/azure/cognitive-services/speech-service/quickstarts/setup-platform?tabs=dotnet%2Cwindows%2Cjre%2Cbrowser&pivots=programming-language-go)

Get started with [speech-to-text sample for Go](https://docs.microsoft.com/azure/cognitive-services/speech-service/get-started-speech-to-text?tabs=windowsinstall&pivots=programming-language-go)

Get started with [text-to-speech sample for Go](https://docs.microsoft.com/azure/cognitive-services/speech-service/get-started-text-to-speech?tabs=script%2Cwindowsinstall&pivots=programming-language-go)

This project requires Go 1.13

# Embedded (offline) speech

The SDK also supports embedded (offline) speech recognition, synthesis and translation, which run fully
on-device without using the cloud Speech service. Use `speech.NewEmbeddedSpeechConfigFromPath` (or
`NewEmbeddedSpeechConfigFromPaths`) to point at the directory that contains your offline models and
voices, then reuse the standard recognizer, synthesizer and translation recognizer factories:

```go
config, err := speech.NewEmbeddedSpeechConfigFromPath("/path/to/models")
if err != nil {
    // handle error
}
defer config.Close()

// Select a recognition model and its license text.
config.SetSpeechRecognitionModel("en-US model name", os.Getenv("EMBEDDED_SPEECH_MODEL_LICENSE"))

// The embedded config wraps a regular SpeechConfig, so pass GetSpeechConfig() to the existing factory.
recognizer, err := speech.NewSpeechRecognizerFromConfig(config.GetSpeechConfig(), audioConfig)

// For embedded translation, select a translation model and use the dedicated factory.
config.SetSpeechTranslationModel("translation model name", os.Getenv("EMBEDDED_SPEECH_MODEL_LICENSE"))
translationRecognizer, err := speech.NewTranslationRecognizerFromEmbeddedConfig(config, audioConfig)
```

Embedded speech has additional runtime requirements:

- The embedded native runtime extensions (for example
  `libMicrosoft.CognitiveServices.Speech.extension.embedded.sr`/`.tts` and their `onnxruntime`
  dependency) must be present next to the core Speech SDK library on the load path.
- Licensed speech recognition models, synthesis voices and/or translation models must be installed on
  the device and passed to the config via the model/voice path.
- Embedded speech is a Limited Access feature that requires approval from Microsoft.

## Discovering installed models

If you don't know the exact model or voice name, enumerate what is installed under the configured
path and use one of the returned names:

```go
models, err := config.GetSpeechRecognitionModels() // or GetSpeechTranslationModels()
if err != nil {
    // handle error
}
for _, model := range models {
    fmt.Printf("%s (locales: %v)\n", model.Name(), model.Locales())
    model.Close()
}
config.SetSpeechRecognitionModel(models[0].Name(), os.Getenv("EMBEDDED_SPEECH_MODEL_LICENSE"))
```

## Building and running

Embedded speech uses cgo, so the Go toolchain must be told where the native Speech SDK headers and
libraries live at build time, and the shared libraries must be discoverable at run time. On Linux/x64,
point `CGO_CFLAGS`/`CGO_LDFLAGS` at the embedded Speech SDK package and add its `lib` folder to the
loader path:

```bash
export SPEECHSDK_ROOT=/path/to/SpeechSDK-Embedded-Linux
export CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
export CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$SPEECHSDK_ROOT/lib/x64"

go build ./...
```

The [samples/embedded](samples/embedded) package provides runnable functions. Using the sample runner
in [samples](samples), the arguments are `(modelPath, modelOrVoiceName, wavFile)` followed by the
sample name:

```bash
# Speech-to-text
go run . /path/to/models "en-US model name" input.wav embedded:RecognizeOnceFromWavFile

# Speech translation
go run . /path/to/translation/models "translation model name" input.wav embedded:TranslateOnceFromWavFile

# Text-to-speech
go run . /path/to/voices "en-US voice name" output.wav embedded:SynthesisToWavFile
```

The model/voice license text is read from the `EMBEDDED_SPEECH_MODEL_LICENSE` environment variable
(the same license applies to recognition, synthesis and translation) so it is not passed on the command line.

See the [embedded samples guide](samples/embedded/README.md) for prerequisites, native library setup,
model/voice installation, build flags and a walkthrough of each sample.

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
