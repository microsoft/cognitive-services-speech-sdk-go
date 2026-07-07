// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package embedded

import (
	"fmt"
	"os"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// RecognizeOnceFromWavFile performs a single embedded (offline) recognition from a WAV file.
//
// Arguments:
//   - modelPath: path to a directory that contains offline speech recognition models.
//   - modelName: name of the recognition model to use (see EmbeddedSpeechConfig.GetSpeechRecognitionModels).
//   - file:      path to the input WAV file.
//
// The model license text is read from the EMBEDDED_SPEECH_MODEL_LICENSE environment variable.
func RecognizeOnceFromWavFile(modelPath string, modelName string, file string) {
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

	// List the recognition models available in the configured path.
	models, err := config.GetSpeechRecognitionModels()
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	for _, model := range models {
		fmt.Printf("Found model: %s (locales: %v)\n", model.Name(), model.Locales())
		model.Close()
	}

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

	speechRecognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})
	speechRecognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
	})

	task := speechRecognizer.RecognizeOnceAsync()
	var outcome speech.SpeechRecognitionOutcome
	select {
	case outcome = <-task:
	case <-time.After(15 * time.Second):
		fmt.Println("Timed out")
		return
	}
	defer outcome.Close()
	if outcome.Error != nil {
		fmt.Println("Got an error: ", outcome.Error)
		return
	}
	fmt.Println("Got a recognition!")
	fmt.Println(outcome.Result.Text)
}
