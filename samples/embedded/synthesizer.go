// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package embedded

import (
	"fmt"
	"os"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// SynthesisToWavFile performs embedded (offline) speech synthesis to a WAV file.
//
// Arguments:
//   - modelPath: path to a directory that contains offline synthesis voices.
//   - voiceName: name of the synthesis voice to use.
//   - file:      path to the output WAV file.
//
// The voice license text is read from the EMBEDDED_SPEECH_MODEL_LICENSE environment variable.
func SynthesisToWavFile(modelPath string, voiceName string, file string) {
	audioConfig, err := audio.NewAudioConfigFromWavFileOutput(file)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer audioConfig.Close()

	config, err := speech.NewEmbeddedSpeechConfigFromPath(modelPath)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer config.Close()

	license := os.Getenv("EMBEDDED_SPEECH_MODEL_LICENSE")
	if err = config.SetSpeechSynthesisVoice(voiceName, license); err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	// Offline synthesis produces raw PCM audio.
	if err = config.SetSpeechSynthesisOutputFormat(common.Riff24Khz16BitMonoPcm); err != nil {
		fmt.Println("Got an error: ", err)
		return
	}

	// The embedded config wraps a regular SpeechConfig, so reuse the standard synthesizer factory.
	speechSynthesizer, err := speech.NewSpeechSynthesizerFromConfig(config.GetSpeechConfig(), audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechSynthesizer.Close()

	speechSynthesizer.SynthesisStarted(func(event speech.SpeechSynthesisEventArgs) {
		defer event.Close()
		fmt.Println("Synthesis started.")
	})
	speechSynthesizer.SynthesisCompleted(func(event speech.SpeechSynthesisEventArgs) {
		defer event.Close()
		fmt.Printf("Synthesized, audio length %d.\n", len(event.Result.AudioData))
	})
	speechSynthesizer.SynthesisCanceled(func(event speech.SpeechSynthesisEventArgs) {
		defer event.Close()
		fmt.Println("Received a cancellation.")
	})

	task := speechSynthesizer.SpeakTextAsync("Hello from embedded speech synthesis.")
	var outcome speech.SpeechSynthesisOutcome
	select {
	case outcome = <-task:
	case <-time.After(60 * time.Second):
		fmt.Println("Timed out")
		return
	}
	defer outcome.Close()
	if outcome.Error != nil {
		fmt.Println("Got an error: ", outcome.Error)
		return
	}
	if outcome.Result.Reason != common.SynthesizingAudioCompleted {
		cancellation, _ := speech.NewCancellationDetailsFromSpeechSynthesisResult(outcome.Result)
		fmt.Printf("Synthesis canceled: reason=%v.\n", cancellation.Reason)
		if cancellation.Reason == common.Error {
			fmt.Printf("Synthesis canceled: error=%v, details=%s\n", cancellation.ErrorCode, cancellation.ErrorDetails)
		}
		return
	}
	fmt.Printf("Synthesis finished, audio written to %s.\n", file)
}
