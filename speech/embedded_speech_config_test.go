// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"testing"
)

func TestEmbeddedConfigFromPath(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPath("models")
	if err != nil {
		t.Error("Unexpected error creating embedded speech config: ", err)
		return
	}
	defer config.Close()
	if config.GetSpeechConfig() == nil {
		t.Error("Underlying speech config should not be nil")
	}
}

func TestEmbeddedConfigFromPaths(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPaths([]string{"models1", "models2"})
	if err != nil {
		t.Error("Unexpected error creating embedded speech config: ", err)
		return
	}
	defer config.Close()
	if config.GetSpeechConfig() == nil {
		t.Error("Underlying speech config should not be nil")
	}
}

func TestEmbeddedConfigFromPathsEmpty(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPaths([]string{})
	if err == nil {
		t.Error("Expected an error when no paths are provided")
	}
	if config != nil {
		t.Error("Expected a nil config when no paths are provided")
	}
}

func TestEmbeddedConfigRecognitionModel(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPath("models")
	if err != nil {
		t.Error("Unexpected error creating embedded speech config: ", err)
		return
	}
	defer config.Close()
	name := "en-US model"
	if err = config.SetSpeechRecognitionModel(name, "license text"); err != nil {
		t.Error("Unexpected error setting recognition model: ", err)
	}
	if config.GetSpeechRecognitionModelName() != name {
		t.Error("Recognition model name not properly set")
	}
}

func TestEmbeddedConfigSynthesisVoice(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPath("models")
	if err != nil {
		t.Error("Unexpected error creating embedded speech config: ", err)
		return
	}
	defer config.Close()
	name := "en-US voice name"
	if err = config.SetSpeechSynthesisVoice(name, "license text"); err != nil {
		t.Error("Unexpected error setting synthesis voice: ", err)
	}
	if config.GetSpeechSynthesisVoiceName() != name {
		t.Error("Synthesis voice name not properly set")
	}
}

func TestEmbeddedConfigKeywordModel(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPath("models")
	if err != nil {
		t.Error("Unexpected error creating embedded speech config: ", err)
		return
	}
	defer config.Close()
	name := "computer"
	if err = config.SetKeywordRecognitionModel(name, "license text"); err != nil {
		t.Error("Unexpected error setting keyword model: ", err)
	}
	if config.GetKeywordRecognitionModelName() != name {
		t.Error("Keyword model name not properly set")
	}
}

func TestEmbeddedConfigGetSpeechRecognitionModels(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPath("models")
	if err != nil {
		t.Error("Unexpected error creating embedded speech config: ", err)
		return
	}
	defer config.Close()
	// With no models present in the path this returns an empty list rather than an error.
	models, err := config.GetSpeechRecognitionModels()
	if err != nil {
		t.Error("Unexpected error listing recognition models: ", err)
		return
	}
	for _, model := range models {
		model.Close()
	}
}

func TestEmbeddedConfigTranslationModel(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPath("models")
	if err != nil {
		t.Error("Unexpected error creating embedded speech config: ", err)
		return
	}
	defer config.Close()
	name := "en-US to many"
	if err = config.SetSpeechTranslationModel(name, "license text"); err != nil {
		t.Error("Unexpected error setting translation model: ", err)
	}
	if config.GetSpeechTranslationModelName() != name {
		t.Error("Translation model name not properly set")
	}
}

func TestEmbeddedConfigGetSpeechTranslationModels(t *testing.T) {
	config, err := NewEmbeddedSpeechConfigFromPath("models")
	if err != nil {
		t.Error("Unexpected error creating embedded speech config: ", err)
		return
	}
	defer config.Close()
	// With no models present in the path this returns an empty list rather than an error.
	models, err := config.GetSpeechTranslationModels()
	if err != nil {
		t.Error("Unexpected error listing translation models: ", err)
		return
	}
	for _, model := range models {
		model.Close()
	}
}
