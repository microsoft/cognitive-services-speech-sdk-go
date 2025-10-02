// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package dialog

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func createConnectorFromSubscriptionRegionAndAudioConfig(t *testing.T, subscription string, region string, audioConfig *audio.AudioConfig) *DialogServiceConnector {
	config, err := NewBotFrameworkConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	defer config.Close()
	connector, err := NewDialogServiceConnectorFromConfig(config, audioConfig)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	return connector
}

func createConnectorFromSubscriptionRegionAndFileInput(t *testing.T, subscription string, region string, file string) *DialogServiceConnector {
	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		t.Error("Got an error: ", err.Error())
		return nil
	}
	defer audioConfig.Close()
	return createConnectorFromSubscriptionRegionAndAudioConfig(t, subscription, region, audioConfig)
}

func createConnectorFromFileInput(t *testing.T, file string) *DialogServiceConnector {
	subscription := os.Getenv("DIALOG_SUBSCRIPTION_KEY")
	region := os.Getenv("DIALOG_SUBSCRIPTION_REGION")
	return createConnectorFromSubscriptionRegionAndFileInput(t, subscription, region, file)
}

func createConnectorFromAudioConfig(t *testing.T, audioConfig *audio.AudioConfig) *DialogServiceConnector {
	subscription := os.Getenv("DIALOG_SUBSCRIPTION_KEY")
	region := os.Getenv("DIALOG_SUBSCRIPTION_REGION")
	return createConnectorFromSubscriptionRegionAndAudioConfig(t, subscription, region, audioConfig)
}

func TestSessionEvents(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	sessionStartedFuture := make(chan bool)
	sessionStartedHandler := func(event speech.SessionEventArgs) {
		defer event.Close()
		sessionStartedFuture <- true
		id := event.SessionID
		t.Log("Started ", id)
	}
	sessionStoppedFuture := make(chan bool)
	sessionStoppedHandler := func(event speech.SessionEventArgs) {
		defer event.Close()
		sessionStoppedFuture <- true
		id := event.SessionID
		t.Log("Stopped ", id)
	}
	cancellationHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Log("Got a cancellation...")
		t.Log(event.ErrorDetails)
	}
	connector.SessionStarted(sessionStartedHandler)
	connector.SessionStopped(sessionStoppedHandler)
	connector.Canceled(cancellationHandler)
	future := connector.ListenOnceAsync()
	outcome := <-future
	defer outcome.Close()
	if outcome.Failed() {
		t.Error("Got an error: ", outcome.Error.Error())
		return
	}
	result := outcome.Result
	t.Log("Recognized: ", result.Text)
	select {
	case <-sessionStartedFuture:
		t.Log("Received a SessionStart event")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive SessionStart event.")
	}
	select {
	case <-sessionStoppedFuture:
		t.Log("Received a SessionStop event")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive SessionStop event.")
	}
}

func TestSpeechRecognitionEvents(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	recognizedFuture := make(chan string)
	recognizedHandle := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognized ", event.Result.Text)
		recognizedFuture <- "Recognized"
	}
	recognizingFuture := make(chan string)
	recognizingHandle := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing ", event.Result.Text)
		select {
		case recognizingFuture <- "Recognizing":
			t.Log("Notified listener.")
		default:
			t.Log("No one is listening, ignore...")
		}
	}
	cancellationHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Log("Got a cancellation...")
		t.Log(event.ErrorDetails)
	}
	connector.Recognized(recognizedHandle)
	connector.Recognizing(recognizingHandle)
	connector.Canceled(cancellationHandler)
	connector.ListenOnceAsync()
	select {
	case <-recognizingFuture:
		t.Log("Received at least one Recognizing event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't received Recognizing events.")
	}
	select {
	case <-recognizedFuture:
		t.Log("Received a Recognized event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive Recognizing event.")
	}
}

func TestCancellationEvent(t *testing.T) {
	region := os.Getenv("DIALOG_SUBSCRIPTION_REGION")
	connector := createConnectorFromSubscriptionRegionAndFileInput(t, "bad_suscription", region, "../test_files/turn_on_the_lamp.wav")
	defer connector.Close()
	future := make(chan string)
	cancelledHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		future <- "Received"
		t.Log("Received cancellation.")
	}
	connector.Canceled(cancelledHandler)
	connector.ListenOnceAsync()
	select {
	case <-future:
		t.Log("All good, received the event.")
	case <-time.After((5 * time.Second)):
		t.Error("Timeout, no event received")
	}
}

type testActivity struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func TestActivityReceivedEvent(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	future := make(chan string)
	activityReceivedHandler := func(event ActivityReceivedEventArgs) {
		defer event.Close()
		future <- "Received"
		t.Log("Received Activity")
	}
	cancellationHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Log("Got a cancellation...")
		t.Log(event.ErrorDetails)
	}
	connector.ActivityReceived(activityReceivedHandler)
	connector.Canceled(cancellationHandler)
	act := testActivity{Type: "message", Text: "Make this larger"}
	msg, _ := json.Marshal(act)
	connector.SendActivityAsync(string(msg))
	select {
	case <-future:
		t.Log("All good, received the event.")
	case <-time.After(5 * time.Second):
		t.Error("Timeout, no event received")
	}
}

func TestActivityWithAudio(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	future := make(chan bool)
	activityReceivedHandler := func(event ActivityReceivedEventArgs) {
		defer event.Close()
		var activity map[string]interface{}
		json.Unmarshal([]byte(event.Activity), &activity)
		messageType := activity["type"].(string)
		if messageType == "conversationUpdate" {
			t.Log("Got conversation update, ignoring")
			return
		}
		if event.HasAudio() {
			audio, err := event.GetAudio()
			if err != nil {
				t.Log("Got an error ", err.Error())
				future <- false
				return
			}
			i := 1
			for buffer, err := audio.Read(3200); (err == nil) && (len(buffer) > 0); buffer, err = audio.Read(3200) {
				t.Log("Got ", len(buffer), " bytes(", i, ")")
				i += 1
			}
			if err != nil {
				t.Log("Got an error ", err.Error())
				future <- false
				return
			}
			future <- true
		} else {
			future <- false
		}
	}
	cancellationHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Log("Got a cancellation...")
		t.Log(event.ErrorDetails)
	}
	connector.ActivityReceived(activityReceivedHandler)
	connector.Canceled(cancellationHandler)
	act := testActivity{Type: "message", Text: "what is the weather forecast in the mountain?"}
	msg, _ := json.Marshal(act)
	connector.SendActivityAsync(string(msg))
	select {
	case hasAudio := <-future:
		if !hasAudio {
			t.Error("No audio")
		} else {
			t.Log("Got audio")
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout, no event received")
	}
}

func TestConnectionFunctions(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	outcome := connector.ConnectAsync()
	err := <-outcome
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	t.Log("Connect Succeeded")
	outcome = connector.DisconnectAsync()
	err = <-outcome
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	t.Log("Disconnect Succeeded")
}

func TestSendActivity(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	act := testActivity{Type: "message", Text: "Make this larger"}
	msg, _ := json.Marshal(act)
	future := connector.SendActivityAsync(string(msg))
	outcome := <-future
	if outcome.Failed() {
		t.Error("Got an error ", outcome.Error.Error())
	} else {
		t.Log("Got interactionID ", outcome.InteractionID)
	}
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
	stream.Write(buffer[0:0]) // Force a final result at the end.
}

func TestFromPushInputStream(t *testing.T) {
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
	connector := createConnectorFromAudioConfig(t, audioConfig)
	if connector == nil {
		return
	}
	defer connector.Close()
	activityFuture := make(chan bool, 1)
	activityReceived := false
	activityReceivedHandler := func(event ActivityReceivedEventArgs) {
		defer event.Close()
		var activity map[string]interface{}
		json.Unmarshal([]byte(event.Activity), &activity)
		messageType := activity["type"].(string)
		if messageType == "conversationUpdate" {
			t.Log("Got conversation update, ignoring")
			return
		}
		t.Log(event.Activity)
		t.Log("Received Activity")
		if activityReceived {
			return
		}
		activityReceived = true
		activityFuture <- true
	}
	connector.ActivityReceived(activityReceivedHandler)
	recognizedFuture := make(chan bool, 1)
	recognizedHandle := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognized ", event.Result.Text)
		recognizedFuture <- true
	}
	connector.Recognized(recognizedHandle)
	recognizingHandler := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing ", event.Result.Text)
	}
	connector.Recognizing(recognizingHandler)
	cancellationHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Log("Got a cancellation...")
		t.Log(event.ErrorDetails)
	}
	connector.Canceled(cancellationHandler)
	pumpFileIntoStream(t, "../test_files/turn_on_the_lamp.wav", stream)
	connector.ListenOnceAsync()
	select {
	case correct := <-recognizedFuture:
		if correct {
			t.Log("All good, received expected recognition event.")
		} else {
			t.Error("Bad recognition")
		}
	case <-time.After(10 * time.Second):
		t.Error("Timeout, no recognition event received.")
	}
	select {
	case correct := <-activityFuture:
		if correct {
			t.Log("All good, received expected activity event.")
		} else {
			t.Error("Bad activity event")
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout, no activity event received.")
	}
}

func TestKeyword(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/whats_the_weather_like.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	model, err := speech.NewKeywordRecognitionModelFromFile("../test_files/kws.table")
	if err != nil {
		t.Error("Found an error: ", err)
	}
	defer model.Close()
	activityFuture := make(chan bool)
	activityReceivedHandler := func(event ActivityReceivedEventArgs) {
		defer event.Close()
		var activity map[string]interface{}
		json.Unmarshal([]byte(event.Activity), &activity)
		messageType := activity["type"].(string)
		if messageType == "conversationUpdate" {
			t.Log("Got conversation update, ignoring")
			return
		}
		t.Log(event.Activity)
		t.Log("Received Activity")
		select {
		case activityFuture <- true:
			t.Log("Notified listener.")
		default:
			t.Log("No one is listening, ignore...")
		}
	}
	connector.ActivityReceived(activityReceivedHandler)
	recognizedFuture := make(chan bool)
	keywordFuture := make(chan bool)
	recognizedHandle := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognized ", event.Result.Text)
		if event.Result.Reason == common.RecognizedKeyword {
			keywordFuture <- true
		} else {
			recognizedFuture <- true
		}
	}
	connector.Recognized(recognizedHandle)
	recognizingHandler := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing ", event.Result.Text)
	}
	connector.Recognizing(recognizingHandler)
	cancellationHandler := func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Log("Got a cancellation...")
		t.Log(event.ErrorDetails)
	}
	connector.Canceled(cancellationHandler)
	err = <-connector.StartKeywordRecognitionAsync(model)
	if err != nil {
		t.Error("Found an error: ", err)
	}
	select {
	case correct := <-keywordFuture:
		if correct {
			t.Log("All good, received expected keyword recognition event.")
		} else {
			t.Error("Bad keyword recognition")
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout, no keyword recognition event received.")
	}
	select {
	case correct := <-recognizedFuture:
		if correct {
			t.Log("All good, received expected recognition event.")
		} else {
			t.Error("Bad recognition")
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout, no recognition event received.")
	}
	select {
	case correct := <-activityFuture:
		if correct {
			t.Log("All good, received expected activity event.")
		} else {
			t.Error("Bad activity event")
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout, no activity event received.")
	}
}
