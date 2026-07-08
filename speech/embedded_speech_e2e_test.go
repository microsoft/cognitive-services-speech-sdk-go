// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// These end-to-end tests exercise the embedded (offline) speech features against real licensed models.
// They only run when the EMBEDDED_MODELS_DIR environment variable points to a directory that contains the
// offline models and test audio using the following layout:
//
//	EMBEDDED_MODELS_DIR=/path/to/models
//
// The directory is expected to contain:
//   - rnnt/Models/SR                     (speech recognition model)
//   - rnnt/Models/ST                     (speech translation models)
//   - audio/whatstheweatherlike.wav      (English utterance)
//   - audio/de-de/CallTheFirstOne.wav    (German utterance)
//
// If the configured models use an empty license/key string, pass "" as the license. Running these tests
// also requires the embedded runtime extension shared libraries from the Speech SDK embedded package to be
// discoverable at load time (for example via LD_LIBRARY_PATH on Linux).
func embeddedModelsDir(t *testing.T) string {
	dir := os.Getenv("EMBEDDED_MODELS_DIR")
	if dir == "" {
		t.Skip("EMBEDDED_MODELS_DIR is not set; skipping embedded end-to-end test")
	}
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("EMBEDDED_MODELS_DIR %q is not accessible: %v", dir, err)
	}
	return dir
}

func TestEmbeddedSpeechRecognitionE2E(t *testing.T) {
	root := embeddedModelsDir(t)
	config, err := NewEmbeddedSpeechConfigFromPath(filepath.Join(root, "rnnt", "Models", "SR"))
	if err != nil {
		t.Fatal("Unexpected error creating embedded speech config: ", err)
	}
	defer config.Close()

	models, err := config.GetSpeechRecognitionModels()
	if err != nil {
		t.Fatal("Unexpected error listing recognition models: ", err)
	}
	if len(models) == 0 {
		t.Fatal("Expected at least one embedded speech recognition model")
	}
	model := models[0]
	t.Logf("Using recognition model: %s (locales: %v)", model.Name(), model.Locales())
	if err = config.SetSpeechRecognitionModel(model.Name(), ""); err != nil {
		t.Fatal("Unexpected error setting recognition model: ", err)
	}
	for _, m := range models {
		m.Close()
	}

	audioConfig, err := audio.NewAudioConfigFromWavFileInput(filepath.Join(root, "audio", "whatstheweatherlike.wav"))
	if err != nil {
		t.Fatal("Unexpected error creating audio config: ", err)
	}
	defer audioConfig.Close()

	recognizer, err := NewSpeechRecognizerFromConfig(config.GetSpeechConfig(), audioConfig)
	if err != nil {
		t.Fatal("Unexpected error creating speech recognizer: ", err)
	}
	defer recognizer.Close()

	select {
	case outcome := <-recognizer.RecognizeOnceAsync():
		if outcome.Error != nil {
			t.Fatal("Recognition failed: ", outcome.Error)
		}
		defer outcome.Result.Close()
		if outcome.Result.Reason != common.RecognizedSpeech {
			t.Fatalf("Unexpected reason: %v", outcome.Result.Reason)
		}
		t.Logf("Recognized: %s", outcome.Result.Text)
		if !strings.Contains(strings.ToLower(outcome.Result.Text), "weather") {
			t.Errorf("Expected recognized text to contain 'weather', got: %s", outcome.Result.Text)
		}
	case <-time.After(60 * time.Second):
		t.Fatal("Timeout waiting for recognition result")
	}
}

func TestEmbeddedSpeechTranslationE2E(t *testing.T) {
	root := embeddedModelsDir(t)
	config, err := NewEmbeddedSpeechConfigFromPath(filepath.Join(root, "rnnt", "Models", "ST"))
	if err != nil {
		t.Fatal("Unexpected error creating embedded speech config: ", err)
	}
	defer config.Close()

	models, err := config.GetSpeechTranslationModels()
	if err != nil {
		t.Fatal("Unexpected error listing translation models: ", err)
	}
	if len(models) == 0 {
		t.Fatal("Expected at least one embedded speech translation model")
	}
	// Pick a model that translates into English.
	var selected *SpeechTranslationModelInfo
	for _, m := range models {
		for _, target := range m.TargetLanguages() {
			if strings.HasPrefix(strings.ToLower(target), "en") {
				selected = m
				break
			}
		}
		if selected != nil {
			break
		}
	}
	if selected == nil {
		selected = models[0]
	}
	t.Logf("Using translation model: %s (%v -> %v)", selected.Name(), selected.SourceLanguages(), selected.TargetLanguages())
	if err = config.SetSpeechTranslationModel(selected.Name(), ""); err != nil {
		t.Fatal("Unexpected error setting translation model: ", err)
	}
	for _, m := range models {
		m.Close()
	}

	audioConfig, err := audio.NewAudioConfigFromWavFileInput(filepath.Join(root, "audio", "de-de", "CallTheFirstOne.wav"))
	if err != nil {
		t.Fatal("Unexpected error creating audio config: ", err)
	}
	defer audioConfig.Close()

	recognizer, err := NewTranslationRecognizerFromEmbeddedConfig(config, audioConfig)
	if err != nil {
		t.Fatal("Unexpected error creating translation recognizer: ", err)
	}
	defer recognizer.Close()

	select {
	case outcome := <-recognizer.RecognizeOnceAsync():
		if outcome.Error != nil {
			t.Fatal("Translation failed: ", outcome.Error)
		}
		defer outcome.Result.Close()
		if outcome.Result.Reason != common.TranslatedSpeech {
			t.Fatalf("Unexpected reason: %v", outcome.Result.Reason)
		}
		translations := outcome.Result.GetTranslations()
		t.Logf("Translations: %v", translations)
		if len(translations) == 0 {
			t.Error("Expected at least one translation")
		}
	case <-time.After(60 * time.Second):
		t.Fatal("Timeout waiting for translation result")
	}
}
