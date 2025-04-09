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

func setupConversation(t *testing.T) (teardown func()) {
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

func createConversationTranscriberFromSubscriptionRegionAndAudioConfig(t *testing.T, subscription string, region string, audioConfig *audio.AudioConfig) *ConversationTranscriber {
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	defer config.Close()

	// Create a conversation transcriber
	transcriber, err := NewConversationTranscriberFromConfig(config, audioConfig)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	return transcriber
}

func createConversationTranscriberFromAudioConfig(t *testing.T, audioConfig *audio.AudioConfig) *ConversationTranscriber {
	subscription := os.Getenv("SPEECH_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEECH_SUBSCRIPTION_REGION")
	return createConversationTranscriberFromSubscriptionRegionAndAudioConfig(t, subscription, region, audioConfig)
}

func createConversationTranscriberFromFileInput(t *testing.T, file string) *ConversationTranscriber {
	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	defer audioConfig.Close()
	return createConversationTranscriberFromAudioConfig(t, audioConfig)
}

func TestConversationTranscriberCreation(t *testing.T) {
	teardown := setupConversation(t)
	defer teardown()

	// Using environment variables for credentials
	subscription := os.Getenv("SPEECH_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEECH_SUBSCRIPTION_REGION")

	// Create a speech config
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Fatalf("Failed to create speech config: %v", err)
	}
	defer config.Close()

	// Create an audio config from single speaker file.
	audioConfig, err := audio.NewAudioConfigFromWavFileInput("../test_files/turn_on_the_lamp.wav")
	if err != nil {
		t.Skipf("Failed to create audio config from file: %v", err)
		return
	}
	defer audioConfig.Close()

	// Create a conversation transcriber
	transcriber, err := NewConversationTranscriberFromConfig(config, audioConfig)
	if err != nil {
		t.Fatalf("Failed to create conversation transcriber: %v", err)
	}
	defer transcriber.Close()

	// Verify transcriber was created successfully
	if transcriber.handle == nil {
		t.Fatalf("Transcriber handle is nil")
	}
}

func TestConversationTranscriberProperties(t *testing.T) {
	teardown := setupConversation(t)
	defer teardown()

	transcriber := createConversationTranscriberFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if transcriber == nil {
		t.Error("Transcriber creation failed")
		return
	}
	defer transcriber.Close()

	// Test property access
	// Set and get auth token
	testToken := "test_token_value"
	err := transcriber.SetAuthorizationToken(testToken)
	if err != nil {
		t.Fatalf("Failed to set auth token: %v", err)
	}

	token := transcriber.AuthorizationToken()
	if token != testToken {
		t.Fatalf("Auth token mismatch, expected %s, got %s", testToken, token)
	}
}

func TestConversationTranscriberEvents(t *testing.T) {
	teardown := setupConversation(t)
	defer teardown()

	transcriber := createConversationTranscriberFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if transcriber == nil {
		t.Error("Transcriber creation failed")
		return
	}
	defer transcriber.Close()

	sessionStartedFuture := make(chan bool)
	sessionStoppedFuture := make(chan bool)
	speechStartFuture := make(chan bool)
	speechEndFuture := make(chan bool)

	transcriber.SessionStarted(func(event SessionEventArgs) {
		defer event.Close()
		t.Log("SessionStarted")
		sessionStartedFuture <- true
	})
	transcriber.SessionStopped(func(event SessionEventArgs) {
		defer event.Close()
		t.Log("SessionStopped")
		sessionStoppedFuture <- true
	})
	transcriber.SpeechStartDetected(func(event RecognitionEventArgs) {
		defer event.Close()
		t.Log("SpeechStart")
		speechStartFuture <- true
	})
	transcriber.SpeechEndDetected(func(event RecognitionEventArgs) {
		defer event.Close()
		t.Log("SpeechEnd")
		speechEndFuture <- true
	})

	// Start transcribing
	err := <-transcriber.StartTranscribingAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}

	// Wait for session started event
	select {
	case <-sessionStartedFuture:
		t.Log("Received SessionStarted event.")
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SessionStarted event.")
	}

	// Wait for speech start event
	select {
	case <-speechStartFuture:
		t.Log("Received SpeechStart event.")
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SpeechStart event.")
	}

	// Wait for speech end event
	select {
	case <-speechEndFuture:
		t.Log("Received SpeechEnd event.")
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SpeechEnd event.")
	}

	// Stop transcribing
	err = <-transcriber.StopTranscribingAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}

	// Wait for session stopped event
	select {
	case <-sessionStoppedFuture:
		t.Log("Received SessionStopped event.")
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for SessionStopped event.")
	}
}

func TestConversationTranscriberSingleSpeaker(t *testing.T) {
	teardown := setupConversation(t)
	defer teardown()

	transcriber := createConversationTranscriberFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if transcriber == nil {
		t.Error("Transcriber creation failed")
		return
	}
	defer transcriber.Close()

	transcribingFuture := make(chan bool)
	transcribedFuture := make(chan bool)
	
	transcribedHandler := func(event ConversationTranscriptionEventArgs) {
		defer event.Close()
		t.Log("Transcribed text: ", event.Result.Text)
		t.Log("Speaker ID: ", event.Result.SpeakerId)
		transcribedFuture <- true
	}
	
	transcribingHandler := func(event ConversationTranscriptionEventArgs) {
		defer event.Close()
		t.Log("Transcribing text: ", event.Result.Text)
		t.Log("Speaker ID: ", event.Result.SpeakerId)
		select {
		case transcribingFuture <- true:
		default:
		}
	}
	
	transcriber.Transcribed(transcribedHandler)
	transcriber.Transcribing(transcribingHandler)
	
	// Start transcribing
	err := <-transcriber.StartTranscribingAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
	
	// Wait for transcribing event
	select {
	case <-transcribingFuture:
		t.Log("Received Transcribing event.")
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for Transcribing event.")
	}
	
	// Wait for transcribed event
	select {
	case <-transcribedFuture:
		t.Log("Received Transcribed event.")
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for Transcribed event.")
	}
	
	// Stop transcribing
	err = <-transcriber.StopTranscribingAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
}

// Helper function to get keys from a map for display purposes
func getKeysFromMap(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestConversationTranscriberContinuousRecognition(t *testing.T) {
	teardown := setupConversation(t)
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
	
	transcriber := createConversationTranscriberFromAudioConfig(t, audioConfig)
	if transcriber == nil {
		t.Error("Transcriber creation failed")
		return
	}
	defer transcriber.Close()
	
	firstResult := true
	transcribedFuture := make(chan string, 10)
	transcribingFuture := make(chan string, 10)
	sessionStoppedFuture := make(chan bool, 1)
	canceledFuture := make(chan bool, 1)
	
	// Channel to collect speaker IDs
	speakerIDsChan := make(chan string, 200)
	
	transcribedHandler := func(event ConversationTranscriptionEventArgs) {
		defer event.Close()
		firstResult = true
		t.Log("Transcribed: ", event.Result.Text)
		t.Log("Speaker ID: ", event.Result.SpeakerId)
		
		// Send speaker ID to the channel if it's not empty
		if event.Result.SpeakerId != "" && event.Result.SpeakerId != "Unknown" {
			speakerIDsChan <- event.Result.SpeakerId
		}
		
		transcribedFuture <- "Transcribed"
	}
	
	transcribingHandler := func(event ConversationTranscriptionEventArgs) {
		defer event.Close()
		t.Log("Transcribing: ", event.Result.Text)
		t.Log("Speaker ID: ", event.Result.SpeakerId)
		
		// Send speaker ID to the channel if it's not empty
		if event.Result.SpeakerId != "" && event.Result.SpeakerId != "Unknown" {
			speakerIDsChan <- event.Result.SpeakerId
		}
		
		if firstResult {
			firstResult = false
			transcribingFuture <- "Transcribing"
		}
	}
	
	transcriber.Transcribed(transcribedHandler)
	transcriber.Transcribing(transcribingHandler)
	transcriber.SessionStopped(func(event SessionEventArgs) {
		defer event.Close()
		t.Log("SessionStopped")
		sessionStoppedFuture <- true
	})
	transcriber.Canceled(func(event ConversationTranscriptionCanceledEventArgs) {
		t.Log("Canceled event fired")
		if event.Reason == common.EndOfStream {
			canceledFuture <- true
			return
		}
		t.Error("Canceled was not due to EOS " + event.ErrorDetails)
	})
	
	err = <-transcriber.StartTranscribingAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
	
	// Pump audio data into the stream
	pumpFileIntoStream(t, "../test_files/katiesteve_mono.wav", stream)
	pumpFileIntoStream(t, "../test_files/katiesteve_mono.wav", stream)
	pumpSilenceIntoStream(t, stream)
	stream.CloseStream()
	
	// Wait for first transcribing event
	select {
	case <-transcribingFuture:
		t.Log("Received first Transcribing event.")
	case <-time.After(30 * time.Second):
		t.Error("Didn't receive first Transcribing event.")
	}
	
	// Wait for first transcribed event
	select {
	case <-transcribedFuture:
		t.Log("Received first Transcribed event.")
	case <-time.After(30 * time.Second):
		t.Error("Didn't receive first Transcribed event.")
	}
	
	// Wait for second transcribing event
	select {
	case <-transcribingFuture:
		t.Log("Received second Transcribing event.")
	case <-time.After(30 * time.Second):
		t.Error("Didn't receive second Transcribing event.")
	}
	
	// Wait for second transcribed event
	select {
	case <-transcribedFuture:
		t.Log("Received second Transcribed event.")
	case <-time.After(30 * time.Second):
		t.Error("Didn't receive second Transcribed event.")
	}
	
	err = <-transcriber.StopTranscribingAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
	
	// Wait for session stopped event
	select {
	case <-sessionStoppedFuture:
		t.Log("Received SessionStopped event.")
	case <-time.After(30 * time.Second):
		t.Error("Timeout waiting for SessionStopped event.")
	}
	
	// Close the speaker ID channel to signal we're done collecting IDs
	close(speakerIDsChan)
	
	// Collect unique speaker IDs
	uniqueSpeakerIDs := make(map[string]bool)
	for speakerID := range speakerIDsChan {
		uniqueSpeakerIDs[speakerID] = true
	}
	
	// Verify that more than one speaker ID was detected
	if len(uniqueSpeakerIDs) <= 1 {
		t.Errorf("Expected more than 1 unique speaker ID, but got %d: %v", 
			len(uniqueSpeakerIDs), getKeysFromMap(uniqueSpeakerIDs))
	} else {
		t.Logf("Successfully detected %d unique speaker IDs: %v", 
			len(uniqueSpeakerIDs), getKeysFromMap(uniqueSpeakerIDs))
	}
}
