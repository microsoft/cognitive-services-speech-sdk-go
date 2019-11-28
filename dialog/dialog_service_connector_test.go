package dialog

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"testing"
	"time"
	"os"
	"encoding/json"
)

func createConnectorFromSubscriptionRegionAndFileInput(t *testing.T, subscription string, region string, file string) *DialogServiceConnector {
	var err error
	var audioConfig *audio.AudioConfig
	audioConfig, err = audio.NewAudioConfigFromFileInput(file)
	if err != nil {
		t.Error("Got an error: ", err.Error())
		return nil
	}
	defer audioConfig.Close()
	var config *BotFrameworkConfig
	config, err = NewBotFrameworkConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err.Error())
		return nil
	}
	defer config.Close()
	var connector *DialogServiceConnector
	connector, err = NewDialogServiceConnectorFromConfig(config, audioConfig)
	if err != nil {
		t.Error("Got an error: ", err.Error())
		return nil
	}
	return connector
}

func createConnectorFromFileInput(t *testing.T, file string) *DialogServiceConnector {
	subscription := os.Getenv("TEST_SUBSCRIPTION_KEY")
	region := os.Getenv("TEST_SUBSCRIPTION_REGION")
	return createConnectorFromSubscriptionRegionAndFileInput(t, subscription, region, file)
}

func TestSessionEvents(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	receivedSessionStarted := false
	sessionStartedHandler := func(event speech.SessionEventArgs) {
		receivedSessionStarted = true
		id := event.SessionID()
		t.Log("Started ", id)
	}
	receivedSessionStopped := false
	sessionStoppedHandler := func(event speech.SessionEventArgs) {
		receivedSessionStopped = true
		id := event.SessionID()
		t.Log("Stopped ", id)
	}
	connector.SessionStarted(sessionStartedHandler)
	connector.SessionStopped(sessionStoppedHandler)
	future := connector.ListenOnceAsync()
	outcome := <- future
	defer outcome.Close()
	if outcome.Failed() {
		t.Error("Got an error: ", outcome.Error.Error())
		return
	}
	result := outcome.Result
	t.Log("Recognized: ", result.Text)
	if !receivedSessionStarted {
		t.Error("Didn't receive SessionStart event.")
	}
	if !receivedSessionStopped {
		t.Error("Didn't receive SessionStopped event.")
	}
}

func TestSpeechRecognitionEvents(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	receivedRecognized := false
	recognizedHandle := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		receivedRecognized = true
		t.Log("Recognized ", event.Result.Text)
	}
	receivedRecognizing := false
	recognizingHandle := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		receivedRecognizing = true
		t.Log("Recognizing ", event.Result.Text)
	}
	connector.Recognized(recognizedHandle)
	connector.Recognizing(recognizingHandle)
	future := connector.ListenOnceAsync()
	outcome := <- future
	defer outcome.Close()
	if outcome.Failed() {
		t.Error("Got an error: ", outcome.Error.Error())
		return
	}
	if !receivedRecognized {
		t.Error("Didn't receive Recognized event.")
	}
	if !receivedRecognizing {
		t.Error("Didn't receive Recognizing event.")
	}
}

func TestCancellationEvent(t *testing.T) {
	region := os.Getenv("TEST_SUBSCRIPTION_REGION")
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
	case <- future:
		t.Log("All good, received the event.")
	case <-time.After((5 * time.Second)):
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
	err := <- outcome
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	t.Log("Connect Succeeded")
	outcome = connector.DisconnectAsync()
	err = <- outcome
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	t.Log("Disconnect Succeeded")
}

type testActivity struct {
	Type string    `json:"type"`
	Text string `json:"text"`
}

func TestSendActivity(t *testing.T) {
	connector := createConnectorFromFileInput(t, "../test_files/turn_on_the_lamp.wav")
	if connector == nil {
		return
	}
	defer connector.Close()
	act := testActivity{ Type: "message", Text: "Make this larger" }
	msg, _ := json.Marshal(act)
	future := connector.SendActivityAsync(string(msg))
	outcome := <- future
	if outcome.Failed() {
		t.Error("Got an error ", outcome.Error.Error())
	} else {
		t.Log("Got interactionID ", outcome.InteractionID)
	}
}