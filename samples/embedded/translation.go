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

// TranslateOnceFromWavFile performs a single embedded (offline) speech translation from a WAV file.
//
// Arguments:
//   - modelPath: path to a directory that contains offline speech translation models.
//   - modelName: name of the translation model to use (see EmbeddedSpeechConfig.GetSpeechTranslationModels).
//   - file:      path to the input WAV file.
//
// The model license text is read from the EMBEDDED_SPEECH_MODEL_LICENSE environment variable.
func TranslateOnceFromWavFile(modelPath string, modelName string, file string) {
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

	// List the translation models available in the configured path.
	models, err := config.GetSpeechTranslationModels()
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	for _, model := range models {
		fmt.Printf("Found model: %s (%v -> %v)\n", model.Name(), model.SourceLanguages(), model.TargetLanguages())
		model.Close()
	}

	license := os.Getenv("EMBEDDED_SPEECH_MODEL_LICENSE")
	if err = config.SetSpeechTranslationModel(modelName, license); err != nil {
		fmt.Println("Got an error: ", err)
		return
	}

	translationRecognizer, err := speech.NewTranslationRecognizerFromEmbeddedConfig(config, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer translationRecognizer.Close()

	translationRecognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})
	translationRecognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
	})

	task := translationRecognizer.RecognizeOnceAsync()
	var outcome speech.TranslationRecognitionOutcome
	select {
	case outcome = <-task:
	case <-time.After(15 * time.Second):
		fmt.Println("Timed out")
		return
	}
	if outcome.Error != nil {
		fmt.Println("Got an error: ", outcome.Error)
		return
	}
	defer outcome.Result.Close()
	fmt.Println("Got a translation!")
	for language, translation := range outcome.Result.GetTranslations() {
		fmt.Printf("  %s: %s\n", language, translation)
	}
}
