// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package dialog_service_connector

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/dialog"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/samples/helpers"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func ListenOnceFromStream(subscription string, region string, file string) {
	stream, err := audio.CreatePushAudioInputStream()
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
	config, err := dialog.NewBotFrameworkConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer config.Close()
	connector, err := dialog.NewDialogServiceConnectorFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer connector.Close()
	sessionStartedHandler := func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started")
	}
	sessionStoppedHandler := func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Stopped")
	}
	connector.SessionStarted(sessionStartedHandler)
	connector.SessionStopped(sessionStoppedHandler)
	activityReceivedHandler := func(event dialog.ActivityReceivedEventArgs) {
		defer event.Close()
		fmt.Println("Received an activity.")
	}
	connector.ActivityReceived(activityReceivedHandler)
	recognizedHandle := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		fmt.Println("Recognized ", event.Result.Text)
	}
	connector.Recognized(recognizedHandle)
	recognizingHandler := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		fmt.Println("Recognizing ", event.Result.Text)
	}
	connector.Recognizing(recognizingHandler)
	helpers.PumpFileIntoStream(file, stream)
	connector.ListenOnceAsync()
	<-time.After(10 * time.Second)
}
