# Embedded (offline) speech samples for Go

These samples show how to use the Microsoft Cognitive Services Speech SDK for Go to run speech
**recognition**, **synthesis** and **translation** fully on-device, without connecting to the cloud
Speech service.

| Sample | Function | Scenario |
| ------ | -------- | -------- |
| [recognizer.go](recognizer.go) | `RecognizeOnceFromWavFile` | Single-shot speech-to-text from a WAV file. |
| [continuous.go](continuous.go) | `RecognizeContinuousFromWavFile` | Continuous speech-to-text from a WAV file. |
| [synthesizer.go](synthesizer.go) | `SynthesisToWavFile` | Text-to-speech to a WAV file. |
| [translation.go](translation.go) | `TranslateOnceFromWavFile` | Single-shot speech translation from a WAV file. |

## Prerequisites

Embedded speech is a **Limited Access** feature. Before you can build and run these samples you need:

1. **Approval for embedded speech.** Access to the embedded runtime and to models/voices requires
   registration and approval from Microsoft. See
   [Embedded Speech](https://learn.microsoft.com/azure/ai-services/speech-service/embedded-speech).
2. **The embedded native Speech SDK package** for your platform (contains the core library *and* the
   embedded runtime extensions). See [Native libraries](#native-libraries).
3. **Licensed models and/or voices** installed on the device. See [Models and voices](#models-and-voices).
4. A working **cgo** toolchain (a C compiler such as `gcc`/`clang`), because the SDK is a cgo binding.
5. **Go 1.13** or later.

## Native libraries

Embedded speech needs the on-device runtime extensions that are **not** part of the standard
(cloud-only) Speech SDK package. Download the **embedded** Speech SDK package. It is a superset of the
standard package: it contains everything the standard package has, plus the offline runtimes:

- `libMicrosoft.CognitiveServices.Speech.core` — the core library.
- `...extension.embedded.sr` and `...extension.embedded.sr.runtime` — offline recognition/translation.
- `...extension.embedded.tts` and `...extension.embedded.tts.runtime` — offline synthesis.
- `...extension.onnxruntime` — neural inference used by the offline runtimes.

Because the embedded package is a superset, an embedded user only needs to download the **embedded**
package — not both the standard and embedded packages.

> The native package ships the **engine only**. The actual models and voices are separate, licensed
> assets (see below).

## Models and voices

The speech recognition models, translation models and synthesis voices are **not** included in the
native SDK package. They are separately licensed assets that you obtain through the Limited Access
program and install into a directory on the device. Each sample takes the path to that directory as
its first argument.

If you do not know the exact model or voice name, the samples that recognize/translate call
`EmbeddedSpeechConfig.GetSpeechRecognitionModels` / `GetSpeechTranslationModels` and print the
installed names and locales at startup. You can also enumerate them yourself:

```go
config, _ := speech.NewEmbeddedSpeechConfigFromPath("/path/to/models")
defer config.Close()

models, _ := config.GetSpeechRecognitionModels() // or GetSpeechTranslationModels()
for _, model := range models {
    fmt.Printf("%s (locales: %v)\n", model.Name(), model.Locales())
    model.Close()
}
```

## Environment variables

The model / voice **license text** is read from an environment variable so it does not have to be
passed on the command line. The same license applies to recognition, synthesis and translation:

| Variable | Description |
| -------- | ----------- |
| `EMBEDDED_SPEECH_MODEL_LICENSE` | License text (or the path to it) for the models and voices you use. |

The model / voice **path** and **name** are passed as command-line arguments to the sample runner
(see [Running the samples](#running-the-samples)).

## Building

The SDK is a cgo binding, so the Go toolchain must know where the native headers and libraries are at
build time, and the shared libraries must be discoverable at run time. On Linux, point
`CGO_CFLAGS` / `CGO_LDFLAGS` at the embedded Speech SDK package and add its `lib` folder for your
architecture to the loader path. For x64:

```bash
export SPEECHSDK_ROOT=/path/to/SpeechSDK-Embedded-Linux
export CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
export CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$SPEECHSDK_ROOT/lib/x64"

go build ./...
```

For 64-bit Arm use `lib/arm64`; for 32-bit Arm use `lib/arm32`.

## Running the samples

The samples are invoked through the shared sample runner in [../main.go](../main.go). The runner takes
four positional arguments; for the embedded samples the first three are repurposed as
`(modelPath, modelOrVoiceName, wavFile)` and the **sample name comes last**:

```bash
# Provide the license once for the whole session.
export EMBEDDED_SPEECH_MODEL_LICENSE="<your model/voice license text>"

# Speech-to-text (single shot)
go run . /path/to/models "en-US recognition model name" input.wav embedded:RecognizeOnceFromWavFile

# Speech-to-text (continuous)
go run . /path/to/models "en-US recognition model name" input.wav embedded:RecognizeContinuousFromWavFile

# Speech translation (single shot)
go run . /path/to/translation/models "translation model name" input.wav embedded:TranslateOnceFromWavFile

# Text-to-speech (writes a WAV file)
go run . /path/to/voices "en-US voice name" output.wav embedded:SynthesisToWavFile
```

The input WAV files for recognition and translation should be 16 kHz (or 8 kHz) mono PCM. The
synthesis sample writes a 24 kHz 16-bit mono PCM WAV file.

## Troubleshooting

- **`SPXERR_...` / "model load error" at recognition or synthesis start.** The model or voice was not
  found at the configured path, or the name does not match an installed model/voice. Let the sample
  print the installed models, or enumerate them as shown in [Models and voices](#models-and-voices),
  and pass an exact `Name()` value.
- **Link errors for `Microsoft.CognitiveServices.Speech.core` at build time.** `CGO_LDFLAGS` is not
  pointing at the `lib/<arch>` folder of the **embedded** package, or you are using the standard
  (cloud-only) package.
- **`error while loading shared libraries` at run time.** The `lib/<arch>` folder is not on
  `LD_LIBRARY_PATH`, so the loader cannot find the core library or the embedded runtime extensions.
- **Recognition/synthesis works online but not offline.** You are linking against the standard package,
  which does not contain the embedded runtime extensions. Use the embedded package (see
  [Native libraries](#native-libraries)).

## See also

- Package documentation: [doc.go](doc.go)
- Embedded speech section in the repository [README](../../README.md#embedded-offline-speech)
- [Embedded Speech overview](https://learn.microsoft.com/azure/ai-services/speech-service/embedded-speech)
