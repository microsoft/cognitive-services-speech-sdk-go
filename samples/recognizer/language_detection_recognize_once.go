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

// LanguageDetectionRecognizeOnce demonstrates how to use the AutoDetectSourceLanguageResult
// helper class with single-shot recognition. The SDK will detect the language from a list
// of candidates and then perform recognition.
func LanguageDetectionRecognizeOnce(subscription string, region string, file string) {
	fmt.Println("Language Detection - Recognize Once")
	fmt.Println("====================================")
	fmt.Println("Recognizing from file:", file)
	fmt.Println()

	// Create speech config
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Error creating speech config:", err)
		return
	}
	defer config.Close()

	// Create auto-detect source language config with candidate languages
	// The SDK will detect which of these languages is spoken
	autoDetectConfig, err := speech.NewAutoDetectSourceLanguageConfigFromLanguages(
		[]string{"en-US", "de-DE", "es-MX", "fr-FR"})
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

	// Setup session events
	recognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})

	recognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
	})

	// Perform recognition
	fmt.Println("Candidate languages: en-US, de-DE, es-MX, fr-FR")
	fmt.Println("Recognizing...")
	fmt.Println()

	task := recognizer.RecognizeOnceAsync()
	var outcome speech.SpeechRecognitionOutcome

	select {
	case outcome = <-task:
	case <-time.After(10 * time.Second):
		fmt.Println("Timed out")
		return
	}

	if outcome.Error != nil {
		fmt.Println("Error during recognition:", outcome.Error)
		return
	}
	defer outcome.Close()

	result := outcome.Result

	// Check result
	fmt.Println()
	if result.Reason == common.RecognizedSpeech {
		// NEW: Use the AutoDetectSourceLanguageResult helper class
		// This is much easier than manually accessing properties!
		autoDetectResult := speech.NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(result)

		fmt.Println("=== Recognition Result ===")
		fmt.Printf("Detected Language: %s\n", autoDetectResult.Language)
		fmt.Printf("Recognized Text:   %s\n", result.Text)
		fmt.Printf("Duration:          %v\n", result.Duration)
		fmt.Printf("Offset:            %v\n", result.Offset)

		// NOTE: OLD WAY (no longer necessary):
		// language := result.Properties.GetProperty(
		//     common.SpeechServiceConnectionAutoDetectSourceLanguageResult, "")

	} else if result.Reason == common.NoMatch {
		fmt.Println("Speech could not be recognized.")
		fmt.Printf("Reason: %d\n", result.Reason)

	} else if result.Reason == common.Canceled {
		fmt.Println("Recognition was canceled")
		fmt.Printf("Reason: %d\n", result.Reason)
	} else {
		fmt.Printf("Unexpected result reason: %d\n", result.Reason)
	}
}
