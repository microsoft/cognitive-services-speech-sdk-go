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

// RecognizeContinuousFromWavFile performs continuous embedded (offline) recognition from a WAV file.
//
// Arguments:
//   - modelPath: path to a directory that contains offline speech recognition models.
//   - modelName: name of the recognition model to use (see EmbeddedSpeechConfig.GetSpeechRecognitionModels).
//   - file:      path to the input WAV file.
//
// The model license text is read from the EMBEDDED_SPEECH_MODEL_LICENSE environment variable.
func RecognizeContinuousFromWavFile(modelPath string, modelName string, file string) {
	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
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
	if err = config.SetSpeechRecognitionModel(modelName, license); err != nil {
		fmt.Println("Got an error: ", err)
		return
	}

	// The embedded config wraps a regular SpeechConfig, so reuse the standard recognizer factory.
	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(config.GetSpeechConfig(), audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechRecognizer.Close()

	done := make(chan struct{})
	speechRecognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})
	speechRecognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
		close(done)
	})
	speechRecognizer.Recognizing(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		fmt.Println("Recognizing:", event.Result.Text)
	})
	speechRecognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		fmt.Println("Recognized:", event.Result.Text)
	})
	speechRecognizer.Canceled(func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		if event.Reason == common.EndOfStream {
			fmt.Println("Reached the end of the audio stream.")
			return
		}
		fmt.Printf("Canceled: reason=%v, error=%v, details=%s\n", event.Reason, event.ErrorCode, event.ErrorDetails)
	})

	if err = <-speechRecognizer.StartContinuousRecognitionAsync(); err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechRecognizer.StopContinuousRecognitionAsync()

	// Wait for the session to stop (end of the WAV file) or time out.
	select {
	case <-done:
	case <-time.After(60 * time.Second):
		fmt.Println("Timed out")
	}
}
