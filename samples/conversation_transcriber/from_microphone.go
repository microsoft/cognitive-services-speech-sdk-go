// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package conversation_transcriber

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func sessionStartedHandler(event speech.SessionEventArgs) {
	defer event.Close()
	fmt.Println("Session Started (ID=", event.SessionID, ")")
}

func sessionStoppedHandler(event speech.SessionEventArgs) {
	defer event.Close()
	fmt.Println("Session Stopped (ID=", event.SessionID, ")")
}

func transcribingHandler(event speech.ConversationTranscriptionEventArgs) {
	defer event.Close()
	fmt.Println("Transcribing - Speaker:", event.Result.SpeakerID, "Text:", event.Result.Text)
}

func transcribedHandler(event speech.ConversationTranscriptionEventArgs) {
	defer event.Close()
	fmt.Println("Transcribed - Speaker:", event.Result.SpeakerID, "Text:", event.Result.Text)
}

func canceledHandler(event speech.ConversationTranscriptionCanceledEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation: ", event.ErrorDetails)
	fmt.Println("Cancellation Reason: ", event.Reason)
}

// ContinuousFromMicrophone performs continuous conversation transcription from the default microphone.
// This sample demonstrates how to transcribe a conversation with speaker identification.
func ContinuousFromMicrophone(subscription string, region string, _ string) {
	audioConfig, err := audio.NewAudioConfigFromDefaultMicrophoneInput()
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

	transcriber.SessionStarted(sessionStartedHandler)
	transcriber.SessionStopped(sessionStoppedHandler)
	transcriber.Transcribing(transcribingHandler)
	transcriber.Transcribed(transcribedHandler)
	transcriber.Canceled(canceledHandler)

	transcriber.StartTranscribingAsync()
	defer transcriber.StopTranscribingAsync()

	fmt.Println("Conversation transcription started. Press Enter to stop...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
