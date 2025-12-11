// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// AutoDetectSourceLanguageResult represents the result of automatic source language detection.
// It extracts the detected language from a recognition result.
//
// Example usage:
//
//	result, err := recognizer.RecognizeOnce()
//	if err != nil {
//	    return err
//	}
//	defer result.Close()
//
//	autoDetectResult := NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(result)
//	fmt.Printf("Detected language: %s\n", autoDetectResult.Language)
//
// Added in version 1.x.0
type AutoDetectSourceLanguageResult struct {
	// Language is the detected language in BCP-47 format (e.g., "en-US", "de-DE", "ja-JP").
	// If no language was detected, this will be an empty string.
	Language string
}

// NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult creates an instance of
// AutoDetectSourceLanguageResult from a speech recognition result.
//
// Parameters:
//   - result: The speech recognition result containing the detected language information.
//
// Returns:
//   - A new AutoDetectSourceLanguageResult instance with the detected language.
//     If result is nil or doesn't contain language information, Language will be empty.
//
// Example:
//
//	autoDetectConfig, _ := speech.NewAutoDetectSourceLanguageConfigFromLanguages(
//	    []string{"en-US", "de-DE"})
//	recognizer, _ := speech.NewSpeechRecognizerFomAutoDetectSourceLangConfig(
//	    config, autoDetectConfig, nil)
//	result, _ := recognizer.RecognizeOnce()
//	autoDetectResult := NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(result)
//	fmt.Println("Detected:", autoDetectResult.Language)
func NewAutoDetectSourceLanguageResultFromSpeechRecognitionResult(
	result *SpeechRecognitionResult) *AutoDetectSourceLanguageResult {

	if result == nil || result.Properties == nil {
		return &AutoDetectSourceLanguageResult{Language: ""}
	}

	return &AutoDetectSourceLanguageResult{
		Language: result.Properties.GetProperty(
			common.SpeechServiceConnectionAutoDetectSourceLanguageResult, ""),
	}
}

// NewAutoDetectSourceLanguageResultFromTranslationRecognitionResult creates an instance of
// AutoDetectSourceLanguageResult from a translation recognition result.
//
// Parameters:
//   - result: The translation recognition result containing the detected language information.
//
// Returns:
//   - A new AutoDetectSourceLanguageResult instance with the detected language.
//     If result is nil or doesn't contain language information, Language will be empty.
//
// Example:
//
//	autoDetectConfig, _ := speech.NewAutoDetectSourceLanguageConfigFromLanguages(
//	    []string{"en-US", "de-DE"})
//	recognizer, _ := speech.NewTranslationRecognizerFromAutoDetectSourceLangConfig(
//	    config, autoDetectConfig, nil)
//	result, _ := recognizer.RecognizeOnce()
//	autoDetectResult := NewAutoDetectSourceLanguageResultFromTranslationRecognitionResult(result)
//	fmt.Println("Source language:", autoDetectResult.Language)
func NewAutoDetectSourceLanguageResultFromTranslationRecognitionResult(
	result *TranslationRecognitionResult) *AutoDetectSourceLanguageResult {

	if result == nil || result.Properties == nil {
		return &AutoDetectSourceLanguageResult{Language: ""}
	}

	return &AutoDetectSourceLanguageResult{
		Language: result.Properties.GetProperty(
			common.SpeechServiceConnectionAutoDetectSourceLanguageResult, ""),
	}
}

// NewAutoDetectSourceLanguageResultFromConversationTranscriptionResult creates an instance of
// AutoDetectSourceLanguageResult from a conversation transcription result.
//
// Parameters:
//   - result: The conversation transcription result containing the detected language information.
//
// Returns:
//   - A new AutoDetectSourceLanguageResult instance with the detected language.
//     If result is nil or doesn't contain language information, Language will be empty.
//
// Example:
//
//	autoDetectConfig, _ := speech.NewAutoDetectSourceLanguageConfigFromLanguages(
//	    []string{"en-US", "de-DE"})
//	transcriber, _ := speech.NewConversationTranscriberFromAutoDetectSourceLangConfig(
//	    config, autoDetectConfig, nil)
//	result := // ... get result from event
//	autoDetectResult := NewAutoDetectSourceLanguageResultFromConversationTranscriptionResult(result)
//	fmt.Println("Detected:", autoDetectResult.Language)
func NewAutoDetectSourceLanguageResultFromConversationTranscriptionResult(
	result *ConversationTranscriptionResult) *AutoDetectSourceLanguageResult {

	if result == nil || result.Properties == nil {
		return &AutoDetectSourceLanguageResult{Language: ""}
	}

	return &AutoDetectSourceLanguageResult{
		Language: result.Properties.GetProperty(
			common.SpeechServiceConnectionAutoDetectSourceLanguageResult, ""),
	}
}
