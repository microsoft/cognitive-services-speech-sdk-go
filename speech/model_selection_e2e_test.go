// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

func TestMAITranscribe2StreamingModelSelectionE2E(t *testing.T) {
	if os.Getenv("ENABLE_MAI_TRANSCRIBE_2_STREAMING_TESTS") != "1" {
		t.Skip("ENABLE_MAI_TRANSCRIBE_2_STREAMING_TESTS is not set to 1; skipping preview-model test")
	}

	subscription := os.Getenv("SPEECH_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEECH_SUBSCRIPTION_REGION")
	if subscription == "" || region == "" {
		t.Fatal("SPEECH_SUBSCRIPTION_KEY and SPEECH_SUBSCRIPTION_REGION are required")
	}
	supportedRegions := map[string]bool{
		"centralus":     true,
		"eastus2":       true,
		"southeastasia": true,
		"swedencentral": true,
	}
	if !supportedRegions[strings.ToLower(region)] {
		t.Fatalf("MAI-Transcribe-2-Streaming is not available in region %q", region)
	}

	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Fatal("Unexpected error creating speech config: ", err)
	}
	defer config.Close()
	if err = config.SetModel("mai-transcribe-2-streaming"); err != nil {
		t.Fatal("Unexpected error setting model: ", err)
	}

	audioConfig, err := audio.NewAudioConfigFromWavFileInput("../test_files/turn_on_the_lamp.wav")
	if err != nil {
		t.Fatal("Unexpected error creating audio config: ", err)
	}
	defer audioConfig.Close()

	recognizer, err := NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Fatal("Unexpected error creating speech recognizer: ", err)
	}
	defer recognizer.Close()

	done := make(chan struct{}, 1)
	resultText := make(chan string, 1)
	canceled := make(chan string, 1)
	recognizer.Recognized(func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		if event.Result.Reason == common.RecognizedSpeech && event.Result.Text != "" {
			select {
			case resultText <- event.Result.Text:
			default:
			}
		}
	})
	recognizer.Canceled(func(event SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		select {
		case canceled <- fmt.Sprintf("%v: %s", event.ErrorCode, event.ErrorDetails):
		default:
		}
	})
	recognizer.SessionStopped(func(event SessionEventArgs) {
		defer event.Close()
		select {
		case done <- struct{}{}:
		default:
		}
	})

	if err = <-recognizer.StartContinuousRecognitionAsync(); err != nil {
		t.Fatal("Unexpected error starting recognition: ", err)
	}

	var text string
	select {
	case text = <-resultText:
	case details := <-canceled:
		t.Fatal("Recognition canceled: ", details)
	case <-done:
		t.Fatal("Recognition stopped without a final result")
	case <-time.After(60 * time.Second):
		t.Fatal("Timeout waiting for recognition result")
	}

	if err = <-recognizer.StopContinuousRecognitionAsync(); err != nil {
		t.Fatal("Unexpected error stopping recognition: ", err)
	}
	if !strings.Contains(strings.ToLower(text), "lamp") {
		t.Fatalf("Expected recognized text to contain %q, got %q", "lamp", text)
	}
}
