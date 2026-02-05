// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package conversation_transcriber

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// TranscribeFromFile performs conversation transcription from an audio file.
// This sample demonstrates how to transcribe a conversation from a WAV file with speaker identification.
func TranscribeFromFile(subscription string, region string, file string) {
	if file == "" {
		fmt.Println("Error: file path is required for this sample")
		return
	}

	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer audioConfig.Close()

	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer config.Close()

	transcriber, err := speech.NewConversationTranscriberFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer transcriber.Close()

	// Channel to signal when transcription is done
	done := make(chan bool)

	transcriber.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})

	transcriber.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
		done <- true
	})

	transcriber.Transcribing(func(event speech.ConversationTranscriptionEventArgs) {
		defer event.Close()
		fmt.Println("Transcribing - Speaker:", event.Result.SpeakerID, "Text:", event.Result.Text)
	})

	transcriber.Transcribed(func(event speech.ConversationTranscriptionEventArgs) {
		defer event.Close()
		if event.Result.Reason == common.RecognizedSpeech {
			fmt.Println("Transcribed - Speaker:", event.Result.SpeakerID, "Text:", event.Result.Text)
		} else if event.Result.Reason == common.NoMatch {
			fmt.Println("No speech could be recognized")
		}
	})

	transcriber.Canceled(func(event speech.ConversationTranscriptionCanceledEventArgs) {
		defer event.Close()
		fmt.Println("Canceled: ", event.ErrorDetails)
		fmt.Println("Cancellation Reason: ", event.Reason)
		if event.Reason == common.EndOfStream {
			fmt.Println("End of audio stream reached")
		}
		done <- true
	})

	// Start transcription
	transcriber.StartTranscribingAsync()

	// Wait for transcription to complete with timeout
	select {
	case <-done:
		fmt.Println("Transcription completed")
	case <-time.After(5 * time.Minute):
		fmt.Println("Transcription timed out")
	}

	transcriber.StopTranscribingAsync()
}
