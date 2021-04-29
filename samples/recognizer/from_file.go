// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package recognizer

import (
	"fmt"
	"time"
    "strings"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/samples/helpers"
)

func RecognizeOnceFromWavFile(subscription string, region string, file string) {
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
	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechRecognizer.Close()
	speechRecognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})
	speechRecognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
	})
	task := speechRecognizer.RecognizeOnceAsync()
	var outcome speech.SpeechRecognitionOutcome
	select {
	case outcome = <-task:
	case <-time.After(5 * time.Second):
		fmt.Println("Timed out")
		return
	}
	defer outcome.Close()
	if outcome.Error != nil {
		fmt.Println("Got an error: ", outcome.Error)
	}
	fmt.Println("Got a recognition!")
	fmt.Println(outcome.Result.Text)
}

func RecognizeOnceFromCompressedFile(subscription string, region string, file string) {
	var containerFormat audio.AudioStreamContainerFormat
	if strings.Contains(file, ".mulaw") {
		containerFormat = audio.MULAW
	} else if strings.Contains(file, ".alaw") {
		containerFormat = audio.ALAW
	} else if strings.Contains(file, ".mp3") {
		containerFormat = audio.MP3
	} else if strings.Contains(file, ".flac") {
		containerFormat = audio.FLAC
	} else if strings.Contains(file, ".opus") {
		containerFormat = audio.OGGOPUS
	} else {
		containerFormat = audio.ANY
	}
	format, err := audio.GetCompressedFormat(containerFormat)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer format.Close()
	stream, err := audio.CreatePushAudioInputStreamFromFormat(format)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer stream.Close()
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
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
	speechRecognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})
	speechRecognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
	})
	helpers.PumpFileIntoStream(file, stream)
	task := speechRecognizer.RecognizeOnceAsync()
	var outcome speech.SpeechRecognitionOutcome
	select {
	case outcome = <-task:
	case <-time.After(40 * time.Second):
		fmt.Println("Timed out")
		return
	}
	defer outcome.Close()
	if outcome.Error != nil {
		fmt.Println("Got an error: ", outcome.Error)
	}
	fmt.Println("Got a recognition!")
	fmt.Println(outcome.Result.Text)
}
