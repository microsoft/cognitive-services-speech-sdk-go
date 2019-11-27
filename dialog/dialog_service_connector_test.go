package dialog

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"testing"
	"os"
)

func createConnectorFromFileInput(t *testing.T, file string) *DialogServiceConnector {
	subscription := os.Getenv("TEST_SUBSCRIPTION_KEY")
	region := os.Getenv("TEST_SUBSCRIPTION_REGION")
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
	if outcome.Error != nil {
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