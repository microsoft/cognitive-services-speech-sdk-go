// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package recognizer

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// LanguageDetectionContinuousLID demonstrates how to use continuous language detection with
// Continuous Language ID mode, which can detect language changes mid-stream.
// Perfect for code-switching scenarios or multilingual audio.
func LanguageDetectionContinuousLID(subscription string, region string, file string) {
	fmt.Println("Language Detection - Continuous with Continuous LID Mode")
	fmt.Println("=========================================================")
	fmt.Println("Audio file:", file)
	fmt.Println()

	// Create speech config
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Error creating speech config:", err)
		return
	}
	defer config.Close()

	// Enable Continuous Language ID mode
	// This allows detection of language changes throughout the audio
	fmt.Println("Enabling Continuous Language ID mode...")
	config.SetProperty(common.SpeechServiceConnectionLanguageIDMode, "Continuous")

	// Create auto-detect config with up to 4 candidate languages (limit for Continuous mode)
	// For At-Start mode (default), you can specify up to 10 languages
	autoDetectConfig, err := speech.NewAutoDetectSourceLanguageConfigFromLanguages(
		[]string{"en-US", "de-DE", "es-MX", "ja-JP"})
	if err != nil {
		fmt.Println("Error creating auto-detect config:", err)
		return
	}
	defer autoDetectConfig.Close()

	// Create audio config from file
	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		fmt.Println("Error creating audio config:", err)
		return
	}
	defer audioConfig.Close()

	// Create speech recognizer with auto-detect config
	recognizer, err := speech.NewSpeechRecognizerFomAutoDetectSourceLangConfig(
		config, autoDetectConfig, audioConfig)
	if err != nil {
		fmt.Println("Error creating recognizer:", err)
		return
	}
	defer recognizer.Close()

	// Track detected languages and language switches
	detectedLanguages := make([]string, 0)
	languageSwitches := 0
	lastLanguage := ""

	// Setup session started event
	recognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("\n=== Session Started ===")
		fmt.Println("Candidate languages: en-US, de-DE, es-MX, ja-JP")
		fmt.Println("Listening for speech...\n")
	})

	// Setup recognizing event (intermediate results)
	recognizer.Recognizing(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()

		if event.Result.Reason == common.RecognizingSpeech {
			// Use the new helper class to get detected language
			autoDetectResult := speech.NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(&event.Result)

			if autoDetectResult.Language != "" {
				fmt.Printf("RECOGNIZING [%s]: %s\n", autoDetectResult.Language, event.Result.Text)
			}
		}
	})

	// Setup recognized event (final results)
	recognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()

		if event.Result.Reason == common.RecognizedSpeech {
			// Use the new helper class to get detected language
			autoDetectResult := speech.NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(&event.Result)

			if autoDetectResult.Language != "" {
				// Check if language changed
				if lastLanguage != "" && lastLanguage != autoDetectResult.Language {
					languageSwitches++
					fmt.Printf("\n*** Language switched: %s → %s ***\n\n",
						lastLanguage, autoDetectResult.Language)
				}

				detectedLanguages = append(detectedLanguages, autoDetectResult.Language)
				lastLanguage = autoDetectResult.Language

				fmt.Printf("RECOGNIZED [%s]: %s\n", autoDetectResult.Language, event.Result.Text)
			} else {
				fmt.Printf("RECOGNIZED: %s (no language detected)\n", event.Result.Text)
			}
		} else if event.Result.Reason == common.NoMatch {
			fmt.Println("NOMATCH: Speech could not be recognized")
		}
	})

	// Setup canceled event
	recognizer.Canceled(func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		fmt.Printf("\nCANCELED: Reason=%d\n", event.Reason)
		if event.Reason == common.Error {
			fmt.Printf("CANCELED: ErrorCode=%d\n", event.ErrorCode)
			fmt.Printf("CANCELED: ErrorDetails=%s\n", event.ErrorDetails)
		}
	})

	// Setup session stopped event
	done := make(chan bool)
	recognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("\n=== Session Stopped ===")
		done <- true
	})

	// Start continuous recognition
	fmt.Println("Starting continuous recognition with language detection...")
	fmt.Println("Candidate languages: en-US, de-DE, es-MX, ja-JP")
	fmt.Println("Listening for speech with language detection...\n")

	err = <-recognizer.StartContinuousRecognitionAsync()
	if err != nil {
		fmt.Println("Error starting continuous recognition:", err)
		return
	}

	// Wait for session to stop (or timeout after 30 seconds)
	select {
	case <-done:
		// Session stopped naturally
	case <-time.After(30 * time.Second):
		fmt.Println("\nTimeout reached, stopping recognition...")
	}

	// Stop continuous recognition
	err = <-recognizer.StopContinuousRecognitionAsync()
	if err != nil {
		fmt.Println("Error stopping continuous recognition:", err)
		return
	}

	// Print summary
	fmt.Println("\n=== Summary ===")
	fmt.Printf("Total utterances recognized: %d\n", len(detectedLanguages))
	fmt.Printf("Language switches detected: %d\n", languageSwitches)

	if len(detectedLanguages) > 0 {
		// Count each language
		languageCounts := make(map[string]int)
		for _, lang := range detectedLanguages {
			languageCounts[lang]++
		}

		fmt.Println("\nLanguages detected:")
		for lang, count := range languageCounts {
			fmt.Printf("  %s: %d utterances (%.1f%%)\n",
				lang, count, float64(count)/float64(len(detectedLanguages))*100)
		}
	}
}
