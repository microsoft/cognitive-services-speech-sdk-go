// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
)

// Helper function to get subscription key from environment
func getSubscription() string {
	return os.Getenv("SPEECH_SUBSCRIPTION_KEY")
}

// Helper function to get region from environment
func getRegion() string {
	return os.Getenv("SPEECH_SUBSCRIPTION_REGION")
}

// Helper function to create audio config from test file
func createAudioConfigFromTestFile(t *testing.T, filename string) (*audio.AudioConfig, error) {
	// Test files are in ../test_files/ relative to the speech package
	testFilePath := filepath.Join("..", "test_files", filename)
	return audio.NewAudioConfigFromWavFileInput(testFilePath)
}

func TestAutoDetectSourceLanguageResultFromSpeechRecognitionResult_NilInput(t *testing.T) {
	// Test with nil result - should not panic and return empty language
	autoDetectResult := NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(nil)

	if autoDetectResult == nil {
		t.Error("AutoDetectSourceLanguageResult should not be nil")
	}
	if autoDetectResult.Language != "" {
		t.Errorf("Language should be empty for nil result, got: %s", autoDetectResult.Language)
	}
}

func TestAutoDetectSourceLanguageResultFromTranslationRecognitionResult_NilInput(t *testing.T) {
	// Test with nil result - should not panic and return empty language
	autoDetectResult := NewAutoDetectSourceLanguageResultFromTranslationRecognitionResult(nil)

	if autoDetectResult == nil {
		t.Error("AutoDetectSourceLanguageResult should not be nil")
	}
	if autoDetectResult.Language != "" {
		t.Errorf("Language should be empty for nil result, got: %s", autoDetectResult.Language)
	}
}

func TestAutoDetectSourceLanguageResultFromConversationTranscriptionResult_NilInput(t *testing.T) {
	// Test with nil result - should not panic and return empty language
	autoDetectResult := NewAutoDetectSourceLanguageResultFromConversationTranscriptionResult(nil)

	if autoDetectResult == nil {
		t.Error("AutoDetectSourceLanguageResult should not be nil")
	}
	if autoDetectResult.Language != "" {
		t.Errorf("Language should be empty for nil result, got: %s", autoDetectResult.Language)
	}
}

// Integration test: Recognize once with language detection
func TestAutoDetectSourceLanguageResult_RecognizeOnce(t *testing.T) {
	// Skip if no subscription key
	subscription := getSubscription()
	if subscription == "" {
		t.Skip("COGNITIVE_SERVICE_SPEECH_SUBSCRIPTION_KEY not set")
	}

	region := getRegion()
	if region == "" {
		t.Skip("COGNITIVE_SERVICE_SPEECH_REGION not set")
	}

	// Create speech config
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Fatalf("Failed to create speech config: %v", err)
	}
	defer config.Close()

	// Create auto-detect config with candidate languages
	autoDetectConfig, err := NewAutoDetectSourceLanguageConfigFromLanguages(
		[]string{"en-US", "de-DE"})
	if err != nil {
		t.Fatalf("Failed to create auto-detect config: %v", err)
	}
	defer autoDetectConfig.Close()

	// Use test audio file (whats_the_weather_like.wav is English)
	audioConfig, err := createAudioConfigFromTestFile(t, "whats_the_weather_like.wav")
	if err != nil {
		t.Fatalf("Failed to create audio config: %v", err)
	}
	defer audioConfig.Close()

	// Create recognizer with auto-detect
	recognizer, err := NewSpeechRecognizerFomAutoDetectSourceLangConfig(
		config, autoDetectConfig, audioConfig)
	if err != nil {
		t.Fatalf("Failed to create recognizer: %v", err)
	}
	defer recognizer.Close()

	// Perform recognition
	outcome := <-recognizer.RecognizeOnceAsync()
	if outcome.Error != nil {
		t.Fatalf("Recognition failed: %v", outcome.Error)
	}
	defer outcome.Close()
	result := outcome.Result

	// Skip if no speech was recognized (e.g., service issues)
	if result.Reason != 3 { // RecognizedSpeech
		t.Skipf("Speech not recognized, reason: %d", result.Reason)
	}

	// Get detected language using the new helper class
	autoDetectResult := NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(result)

	// Verify we got a result
	if autoDetectResult == nil {
		t.Error("AutoDetectSourceLanguageResult should not be nil")
	}

	if autoDetectResult.Language == "" {
		t.Error("Language should not be empty")
	}

	// The test file is in English, so we expect en-US (case insensitive)
	detectedLower := strings.ToLower(autoDetectResult.Language)
	if detectedLower != "en-us" {
		t.Errorf("Expected en-US, got: %s", autoDetectResult.Language)
	}

	t.Logf("Successfully detected language: %s", autoDetectResult.Language)
	t.Logf("Recognized text: %s", result.Text)
}

// Integration test: Continuous recognition with language detection
func TestAutoDetectSourceLanguageResult_ContinuousRecognition(t *testing.T) {
	// Skip if no subscription key
	subscription := getSubscription()
	if subscription == "" {
		t.Skip("COGNITIVE_SERVICE_SPEECH_SUBSCRIPTION_KEY not set")
	}

	region := getRegion()
	if region == "" {
		t.Skip("COGNITIVE_SERVICE_SPEECH_REGION not set")
	}

	// Create speech config
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Fatalf("Failed to create speech config: %v", err)
	}
	defer config.Close()

	// Create auto-detect config
	autoDetectConfig, err := NewAutoDetectSourceLanguageConfigFromLanguages(
		[]string{"en-US", "de-DE"})
	if err != nil {
		t.Fatalf("Failed to create auto-detect config: %v", err)
	}
	defer autoDetectConfig.Close()

	// Use test audio file
	audioConfig, err := createAudioConfigFromTestFile(t, "whats_the_weather_like.wav")
	if err != nil {
		t.Fatalf("Failed to create audio config: %v", err)
	}
	defer audioConfig.Close()

	// Create recognizer
	recognizer, err := NewSpeechRecognizerFomAutoDetectSourceLangConfig(
		config, autoDetectConfig, audioConfig)
	if err != nil {
		t.Fatalf("Failed to create recognizer: %v", err)
	}
	defer recognizer.Close()

	// Track detected languages
	detectedLanguages := make([]string, 0)
	done := make(chan bool)

	// Setup recognized event handler
	recognizer.Recognized(func(event SpeechRecognitionEventArgs) {
		defer event.Close()

		if event.Result.Reason == 3 { // RecognizedSpeech
			// Use the helper class to get detected language
			autoDetectResult := NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(&event.Result)
			if autoDetectResult.Language != "" {
				detectedLanguages = append(detectedLanguages, autoDetectResult.Language)
				t.Logf("Detected language: %s, Text: %s", autoDetectResult.Language, event.Result.Text)
			}
		}
	})

	// Setup session stopped handler
	recognizer.SessionStopped(func(event SessionEventArgs) {
		defer event.Close()
		done <- true
	})

	// Start continuous recognition
	err = <-recognizer.StartContinuousRecognitionAsync()
	if err != nil {
		t.Fatalf("Failed to start continuous recognition: %v", err)
	}

	// Wait for completion
	<-done

	// Stop recognition
	err = <-recognizer.StopContinuousRecognitionAsync()
	if err != nil {
		t.Fatalf("Failed to stop continuous recognition: %v", err)
	}

	// Verify we detected at least one language
	if len(detectedLanguages) == 0 {
		t.Error("Should have detected at least one language")
	}

	if len(detectedLanguages) > 0 {
		// All detected languages should be en-US for English audio
		for _, lang := range detectedLanguages {
			langLower := strings.ToLower(lang)
			if langLower != "en-us" {
				t.Errorf("Expected en-US for English audio, got: %s", lang)
			}
		}
		t.Logf("Total utterances with language detection: %d", len(detectedLanguages))
	}
}

// Test that the helper class works with results that don't have language info
func TestAutoDetectSourceLanguageResult_NoLanguageInfo(t *testing.T) {
	// Skip if no subscription key
	subscription := getSubscription()
	if subscription == "" {
		t.Skip("COGNITIVE_SERVICE_SPEECH_SUBSCRIPTION_KEY not set")
	}

	region := getRegion()
	if region == "" {
		t.Skip("COGNITIVE_SERVICE_SPEECH_REGION not set")
	}

	// Create speech config
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Fatalf("Failed to create speech config: %v", err)
	}
	defer config.Close()

	// Use test audio file
	audioConfig, err := createAudioConfigFromTestFile(t, "whats_the_weather_like.wav")
	if err != nil {
		t.Fatalf("Failed to create audio config: %v", err)
	}
	defer audioConfig.Close()

	// Create recognizer WITHOUT auto-detect config
	// This should not include language detection info in the result
	recognizer, err := NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Fatalf("Failed to create recognizer: %v", err)
	}
	defer recognizer.Close()

	// Perform recognition
	outcome := <-recognizer.RecognizeOnceAsync()
	if outcome.Error != nil {
		t.Fatalf("Recognition failed: %v", outcome.Error)
	}
	defer outcome.Close()
	result := outcome.Result

	// Skip if no speech was recognized
	if result.Reason != 3 { // RecognizedSpeech
		t.Skipf("Speech not recognized, reason: %d", result.Reason)
	}

	// Try to get language - should be empty since we didn't use auto-detect
	autoDetectResult := NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(result)

	if autoDetectResult == nil {
		t.Error("AutoDetectSourceLanguageResult should not be nil")
	}

	// Language should be empty when auto-detect is not used
	if autoDetectResult.Language != "" {
		t.Errorf("Language should be empty when auto-detect is not used, got: %s", autoDetectResult.Language)
	}

	t.Logf("Recognized text without language detection: %s", result.Text)
}
