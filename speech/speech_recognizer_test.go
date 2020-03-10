//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package speech

import (
	"os"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
)

func createSpeechRecognizerFromSubscriptionRegionAndAudioConfig(t *testing.T, subscription string, region string, audioConfig *audio.AudioConfig) *SpeechRecognizer {
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	defer config.Close()
	recognizer, err := NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	return recognizer
}

func createSpeechRecognizerFromSubscriptionRegionAndFileInput(t *testing.T, subscription string, region string, file string) *SpeechRecognizer {
	audioConfig, err := audio.NewAudioConfigFromFileInput(file)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	defer audioConfig.Close()
	return createSpeechRecognizerFromSubscriptionRegionAndAudioConfig(t, subscription, region, audioConfig)
}

func createSpeechRecognizerFromFileInput(t *testing.T, file string) *SpeechRecognizer {
	subscription := os.Getenv("SR_SUBSCRIPTION_KEY")
	region := os.Getenv("SR_SUBSCRIPTION_REGION")
	return createSpeechRecognizerFromSubscriptionRegionAndFileInput(t, subscription, region, file)
}

func TestRecognizeOnce(t *testing.T) {
	recognizer := createSpeechRecognizerFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if recognizer == nil {
		return
	}
	defer recognizer.Close()
	recognizedFuture := make(chan string)
	recognizedHandler := func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognized: ", event.Result.Text)
		recognizedFuture <- "Recognized"
	}
	recognizingFuture := make(chan string)
	recognizingHandle := func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing: ", event.Result.Text)
		select {
		case recognizingFuture <- "Recognizing":
		default:
		}
	}
	recognizer.Recognized(recognizedHandler)
	recognizer.Recognizing(recognizingHandle)
	result := recognizer.RecognizeOnceAsync()
	select {
	case <-recognizingFuture:
		t.Log("Received at least one Recognizing event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive Recognizing event.")
	}
	select {
	case <-recognizedFuture:
		t.Log("Received a Recognized event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive Recognizing event.")
	}
	select {
	case <-result:
		t.Log("Result resolved.")
	case <-time.After(5 * time.Second):
		t.Error("Result didn't resolve.")
	}
}
