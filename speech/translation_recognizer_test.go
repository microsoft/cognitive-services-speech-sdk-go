// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/diagnostics"
)

func setup(t *testing.T) (teardown func()) {
	logLineAtStart := diagnostics.GetMemoryLogLineNumNewest()
	diagnostics.StartMemoryLogging()

	return func() {
		diagnostics.StopMemoryLogging()

		if t.Failed() {
			logLineAtEnd := diagnostics.GetMemoryLogLineNumNewest()

			var logLines strings.Builder

			for i := logLineAtStart; i < logLineAtEnd; i++ {
				logLines.WriteString(diagnostics.GetMemoryLogLine(i))
			}

			t.Log(logLines.String())
		}
	}
}

func createTranslationRecognizerFromSubscriptionRegionAndAudioConfig(t *testing.T, subscription string, region string, audioConfig *audio.AudioConfig) *TranslationRecognizer {
	config, err := NewSpeechTranslationConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	defer config.Close()

	err = config.SetSpeechRecognitionLanguage("en-us")
	if err != nil {
		t.Error("Got an error setting source language: ", err)
		return nil
	}

	// Add target languages for translation
	err = config.AddTargetLanguage("es") // Spanish
	if err != nil {
		t.Error("Got an error adding target language: ", err)
		return nil
	}
	err = config.AddTargetLanguage("fr") // French
	if err != nil {
		t.Error("Got an error adding target language: ", err)
		return nil
	}

	recognizer, err := NewTranslationRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	return recognizer
}

func createTranslationRecognizerFromAudioConfig(t *testing.T, audioConfig *audio.AudioConfig) *TranslationRecognizer {
	subscription := os.Getenv("SPEECH_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEECH_SUBSCRIPTION_REGION")
	return createTranslationRecognizerFromSubscriptionRegionAndAudioConfig(t, subscription, region, audioConfig)
}

func createTranslationRecognizerFromFileInput(t *testing.T, file string) *TranslationRecognizer {
	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	defer audioConfig.Close()
	return createTranslationRecognizerFromAudioConfig(t, audioConfig)
}

func TestTranslationSessionEvents(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	recognizer := createTranslationRecognizerFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if recognizer == nil {
		t.Error("Recognizer creation failed")
		return
	}
	defer recognizer.Close()
	sessionStartedFuture := make(chan bool, 1)
	sessionStoppedFuture := make(chan bool, 1)
	speechStartFuture := make(chan bool, 1)
	speechEndFuture := make(chan bool, 1)

	recognizer.SessionStarted(func(event SessionEventArgs) {
		defer event.Close()
		t.Log("SessionStarted")
		sessionStartedFuture <- true
	})
	recognizer.SessionStopped(func(event SessionEventArgs) {
		defer event.Close()
		t.Log("SessionStopped")
		sessionStoppedFuture <- true
	})
	recognizer.SpeechStartDetected(func(event RecognitionEventArgs) {
		defer event.Close()
		t.Log("SpeechStart")
		speechStartFuture <- true
	})
	recognizer.SpeechEndDetected(func(event RecognitionEventArgs) {
		defer event.Close()
		t.Log("SpeechEnd")
		speechEndFuture <- true
	})

	recognizer.RecognizeOnceAsync()
	select {
	case <-sessionStartedFuture:
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SessionStarted event.")
	}
	select {
	case <-speechStartFuture:
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SpeechStart event.")
	}
	select {
	case <-speechEndFuture:
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SpeechEnd event.")
	}
	select {
	case <-sessionStoppedFuture:
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SessionStopped event.")
	}
}

func TestTranslationRecognizeOnce(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	recognizer := createTranslationRecognizerFromFileInput(t, "../test_files/myVoiceIsMyPassportVerifyMe01.wav")
	if recognizer == nil {
		return
	}
	defer recognizer.Close()
	recognizedFuture := make(chan string)
	recognizedHandler := func(event TranslationRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognized text: ", event.Result.Text)
		translations := event.Result.GetTranslations()
		t.Log("Spanish translation: ", translations["es"])
		t.Log("French translation: ", translations["fr"])
		recognizedFuture <- "Recognized"
	}
	recognizingFuture := make(chan string)
	recognizingHandle := func(event TranslationRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing text: ", event.Result.Text)
		translations := event.Result.GetTranslations()
		t.Log("Spanish translation: ", translations["es"])
		t.Log("French translation: ", translations["fr"])
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
		t.Error("Didn't receive Recognized event.")
	}
	select {
	case outcome := <-result:
		if outcome.Error != nil {
			t.Error("Got an error: ", outcome.Error)
		} else {
			translations := outcome.Result.GetTranslations()
			if len(translations) == 0 {
				t.Error("No translations received")
			}
			t.Log("Result translations: ", translations)
		}
	case <-time.After(5 * time.Second):
		t.Error("Result didn't resolve.")
	}
}

func TestTranslationContinuousRecognition(t *testing.T) {
	ImplTranslationContinuousRecognition(t, false)
}

func TestTranslationContinuousRecognitionEos(t *testing.T) {
	ImplTranslationContinuousRecognition(t, true)
}

func ImplTranslationContinuousRecognition(t *testing.T, runToEnd bool) {
	teardown := setup(t)
	defer teardown()

	format, err := audio.GetDefaultInputFormat()
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer format.Close()
	stream, err := audio.CreatePushAudioInputStreamFromFormat(format)
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer stream.Close()
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer audioConfig.Close()
	recognizer := createTranslationRecognizerFromAudioConfig(t, audioConfig)
	if recognizer == nil {
		return
	}
	defer recognizer.Close()
	firstResult := true
	recognizedFuture := make(chan string, 10)
	recognizingFuture := make(chan string, 10)
	canceledFuture := make(chan bool)

	recognizedHandler := func(event TranslationRecognitionEventArgs) {
		defer event.Close()
		firstResult = true
		t.Log("Recognized text: ", event.Result.Text)
		translations := event.Result.GetTranslations()
		t.Log("Spanish translation: ", translations["es"])
		t.Log("French translation: ", translations["fr"])
		recognizedFuture <- "Recognized"
	}
	recognizingHandle := func(event TranslationRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing text: ", event.Result.Text)
		translations := event.Result.GetTranslations()
		t.Log("Spanish translation: ", translations["es"])
		t.Log("French translation: ", translations["fr"])
		if firstResult {
			firstResult = false
			recognizingFuture <- "Recognizing"
		}
	}
	recognizer.Recognized(recognizedHandler)
	recognizer.Recognizing(recognizingHandle)
	recognizer.Canceled(func(event TranslationRecognitionCanceledEventArgs) {
		t.Log("Canceled event fired")
		if event.Reason == common.EndOfStream {
			canceledFuture <- true
			return
		}

		t.Error("Canceled was not due to EOS " + event.ErrorDetails)
	})

	err = <-recognizer.StartContinuousRecognitionAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
	pumpFileIntoStream(t, "../test_files/turn_on_the_lamp.wav", stream)
	pumpFileIntoStream(t, "../test_files/turn_on_the_lamp.wav", stream)
	pumpSilenceIntoStream(t, stream)
	stream.CloseStream()
	select {
	case <-recognizingFuture:
		t.Log("Received first Recognizing event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive first Recognizing event.")
	}
	select {
	case <-recognizedFuture:
		t.Log("Received first Recognized event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive first Recognized event.")
	}
	select {
	case <-recognizingFuture:
		t.Log("Received second Recognizing event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive second Recognizing event.")
	}
	select {
	case <-recognizedFuture:
		t.Log("Received second Recognized event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive second Recognized event.")
	}
	if !runToEnd {
		err = <-recognizer.StopContinuousRecognitionAsync()
		if err != nil {
			t.Error("Got error: ", err)
		}
	} else {
		select {
		case <-canceledFuture:
			t.Log("Cancled EOS")
		case <-time.After(5 * time.Second):
			t.Error("Didn't receive Canceled event.")
		}
	}
}

func TestTranslationSynthesis(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	recognizer := createTranslationRecognizerFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if recognizer == nil {
		return
	}
	defer recognizer.Close()

	// Set voice name for synthesis
	config := recognizer.Properties
	config.SetProperty(common.SpeechServiceConnectionTranslationVoice, "es-ES-ElviraNeural")

	synthesisFuture := make(chan bool, 2)
	synthesisCompleteFuture := make(chan bool, 2)

	synthesisHandler := func(event TranslationSynthesisEventArgs) {
		defer event.Close()
		audioData := event.Result.GetAudioData()
		if len(audioData) > 0 && event.Result.Reason == common.SynthesizingAudio {
			t.Log(time.Now().String()+"Received synthesized audio data of length: ", len(audioData))
			synthesisFuture <- true
		} else if event.Result.Reason == common.SynthesizingAudioCompleted {
			t.Log("Synthesis is complete")
			synthesisCompleteFuture <- true
		}
	}

	recognizer.Synthesizing(synthesisHandler)

	select {
	case outcome := <-recognizer.RecognizeOnceAsync():
		if outcome.Error != nil {
			t.Error("Got an error: ", outcome.Error)
		}
		t.Log(time.Now().String() + "Got result")
	case <-time.After(9 * time.Second):
		t.Error("Recognition result didn't resolve.")
	}

	select {
	case <-synthesisFuture:
		t.Log("Received synthesis event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive synthesis event.")
	}

	select {
	case <-synthesisCompleteFuture:
		t.Log("Received synthesis complete event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive synthesis event.")
	}
}
