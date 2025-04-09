// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func main() {
	// Replace with your own subscription key and region
	subscription := os.Getenv("SPEECH_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEECH_REGION")

	if subscription == "" || region == "" {
		fmt.Println("Please set SPEECH_SUBSCRIPTION_KEY and SPEECH_REGION environment variables")
		return
	}

	// Create a speech config
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Error creating speech config:", err)
		return
	}
	defer config.Close()

	// Create audio config from microphone
	audioConfig, err := audio.NewAudioConfigFromDefaultMicrophoneInput()
	if err != nil {
		fmt.Println("Error creating audio config:", err)
		return
	}
	defer audioConfig.Close()

	// Create a conversation transcriber
	transcriber, err := speech.NewConversationTranscriberFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("Error creating conversation transcriber:", err)
		return
	}
	defer transcriber.Close()

	// Set up event handlers
	transcriber.SessionStarted(func(event speech.SessionEventArgs) {
		fmt.Println("Session started")
	})

	transcriber.SessionStopped(func(event speech.SessionEventArgs) {
		fmt.Println("Session stopped")
	})

	transcriber.SpeechStartDetected(func(event speech.RecognitionEventArgs) {
		fmt.Println("Speech start detected")
	})

	transcriber.SpeechEndDetected(func(event speech.RecognitionEventArgs) {
		fmt.Println("Speech end detected")
	})

	transcriber.Transcribing(func(event speech.ConversationTranscriptionEventArgs) {
		fmt.Printf("Transcribing: [%s] %s\n", event.Result().SpeakerId, event.Result().Text)
	})

	transcriber.Transcribed(func(event speech.ConversationTranscriptionEventArgs) {
		fmt.Printf("Transcribed: [%s] %s\n", event.Result().SpeakerId, event.Result().Text)
	})

	transcriber.Canceled(func(event speech.ConversationTranscriptionCanceledEventArgs) {
		fmt.Printf("Canceled: Reason=%v, ErrorCode=%v, ErrorDetails=%s\n", 
			event.Reason(), event.ErrorCode(), event.ErrorDetails())
	})

	// Start continuous transcription
	errorChan := transcriber.StartTranscribingAsync()
	err = <-errorChan
	if err != nil {
		fmt.Println("Error starting transcription:", err)
		return
	}
	fmt.Println("Transcription started. Speak now...")

	// Set up a signal channel to handle interrupts (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Wait for signal
	<-sigChan
	fmt.Println("Stopping transcription...")

	// Stop transcription
	errorChan = transcriber.StopTranscribingAsync()
	err = <-errorChan
	if err != nil {
		fmt.Println("Error stopping transcription:", err)
		return
	}

	// Wait for final events to process
	time.Sleep(2 * time.Second)
	fmt.Println("Transcription stopped")
}