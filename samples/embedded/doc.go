// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

// Package embedded demonstrates embedded (offline) speech recognition, synthesis and translation.
//
// It provides the following samples:
//   - RecognizeOnceFromWavFile:       single-shot speech recognition from a WAV file.
//   - RecognizeContinuousFromWavFile: continuous speech recognition from a WAV file.
//   - SynthesisToWavFile:             speech synthesis to a WAV file.
//   - TranslateOnceFromWavFile:       single-shot speech translation from a WAV file.
//
// Embedded speech runs fully on-device and does not use the cloud Speech service. It requires:
//   - The embedded native runtime extensions to be present next to the core Speech SDK library.
//   - Licensed speech recognition models, synthesis voices and/or translation models installed on the device.
//
// The three positional command line arguments are repurposed for the embedded samples as
// (modelPath, modelOrVoiceName, file). The model/voice license text is read from an environment
// variable so it does not have to be passed on the command line:
//   - EMBEDDED_SPEECH_MODEL_LICENSE for speech recognition models, synthesis voices and translation models.
//     The same license applies to embedded speech recognition, synthesis and translation.
//
// If you do not know the exact model or voice name, call EmbeddedSpeechConfig.GetSpeechRecognitionModels
// or GetSpeechTranslationModels to enumerate what is installed under modelPath, then pass one of the
// returned Name() values.
package embedded
