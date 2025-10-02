// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"bufio"
	"io"
	"os"
	"strings"
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

func createSpeechRecognizerFromAudioConfig(t *testing.T, audioConfig *audio.AudioConfig) *SpeechRecognizer {
	subscription := os.Getenv("SPEECH_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEECH_SUBSCRIPTION_REGION")
	return createSpeechRecognizerFromSubscriptionRegionAndAudioConfig(t, subscription, region, audioConfig)
}

func createSpeechRecognizerFromFileInput(t *testing.T, file string) *SpeechRecognizer {
	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	defer audioConfig.Close()
	return createSpeechRecognizerFromAudioConfig(t, audioConfig)
}

func pumpFileIntoStream(t *testing.T, filename string, stream *audio.PushAudioInputStream) {
	file, err := os.Open(filename)
	if err != nil {
		t.Error("Error opening file: ", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1000)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			t.Log("Done reading file.")
			break
		}
		if err != nil {
			t.Error("Error reading file: ", err)
			break
		}
		err = stream.Write(buffer[0:n])
		if err != nil {
			t.Error("Error writing to the stream")
		}
	}
}

func pumpSilenceIntoStream(t *testing.T, stream *audio.PushAudioInputStream) {
	buffer := make([]byte, 1000)
	for i := range buffer {
		buffer[i] = 0
	}
	for i := 0; i < 16; i++ {
		err := stream.Write(buffer)
		if err != nil {
			t.Error("Error writing to the stream")
		}
	}
}

func TestSessionEvents(t *testing.T) {
	recognizer := createSpeechRecognizerFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if recognizer == nil {
		t.Error("Recognizer creation failed")
		return
	}
	defer recognizer.Close()
	sessionStartedFuture := make(chan bool)
	sessionStoppedFuture := make(chan bool)
	recognizer.SessionStarted(func(event SessionEventArgs) {
		defer event.Close()
		t.Log("SessionStarted")
		sessionStartedFuture <- true
	})
	recognizer.SessionStopped(func(event SessionEventArgs) {
		defer event.Close()
		t.Log("SessionStarted")
		sessionStoppedFuture <- true
	})
	recognizer.RecognizeOnceAsync()
	select {
	case <-sessionStartedFuture:
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SessionStarted event.")
	}
	select {
	case <-sessionStoppedFuture:
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SessionStopped event.")
	}
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

func TestContinuousRecognition(t *testing.T) {
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
	recognizer := createSpeechRecognizerFromAudioConfig(t, audioConfig)
	if recognizer == nil {
		return
	}
	defer recognizer.Close()
	firstResult := true
	recognizedFuture := make(chan string)
	recognizingFuture := make(chan string)
	recognizedHandler := func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		firstResult = true
		t.Log("Recognized: ", event.Result.Text)
		recognizedFuture <- "Recognized"
	}
	recognizingHandle := func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing: ", event.Result.Text)
		if firstResult {
			firstResult = false
			recognizingFuture <- "Recognizing"
		}
	}
	recognizer.Recognized(recognizedHandler)
	recognizer.Recognizing(recognizingHandle)
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
	err = <-recognizer.StopContinuousRecognitionAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
}

func testPhraseList(t *testing.T, with_grammar bool) {
	recognizer := createSpeechRecognizerFromFileInput(t, "../test_files/peloozoid.wav")
	if recognizer == nil {
		return
	}
	defer recognizer.Close()
	phraseListGrammar, err := NewPhraseListGrammarFromRecognizer(recognizer)
	if err != nil {
		t.Error("Grammar creation failed")
	}
	defer phraseListGrammar.Close()
	if with_grammar {
		phraseListGrammar.AddPhrase("peloozoid")
	}
	var result *SpeechRecognitionResult
	select {
	case outcome := <-recognizer.RecognizeOnceAsync():
		if outcome.Error != nil {
			t.Error("Received an error")
		}
		result = outcome.Result
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for result.")
	}
	defer result.Close()
	if strings.Contains(strings.ToLower(result.Text), "peloozoid") != with_grammar {
		t.Log(result.Text)
		t.Errorf("Result doesn't match expectation (expected 'peloozoid', got '%s')", result.Text)
	}
}

func TestPhraseListGrammarWithoutGrammar(t *testing.T) {
	testPhraseList(t, false)
}

func TestPhraseListGrammarWithGrammar(t *testing.T) {
	testPhraseList(t, true)
}
