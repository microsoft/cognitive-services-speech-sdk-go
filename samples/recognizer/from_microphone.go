//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package recognizer

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

func recognizingHandler(event speech.SpeechRecognitionEventArgs) {
	defer event.Close()
	fmt.Println("Recognizing:", event.Result.Text)
}

func recognizedHandler(event speech.SpeechRecognitionEventArgs) {
	defer event.Close()
	fmt.Println("Recognized:", event.Result.Text)
}

func cancelledHandler(event speech.SpeechRecognitionCanceledEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation: ", event.ErrorDetails)
}

func ContinuousFromMicrophone(subscription string, region string, file string) {
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
	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechRecognizer.Close()
	speechRecognizer.SessionStarted(sessionStartedHandler)
	speechRecognizer.SessionStopped(sessionStoppedHandler)
	speechRecognizer.Recognizing(recognizingHandler)
	speechRecognizer.Recognized(recognizedHandler)
	speechRecognizer.Canceled(cancelledHandler)
	speechRecognizer.StartContinuousRecognitionAsync()
	defer speechRecognizer.StopContinuousRecognitionAsync()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
