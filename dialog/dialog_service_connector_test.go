package dialog

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"testing"
	"fmt"
	"os"
)

func TestSessionEvents(t *testing.T) {
	subscription := os.Getenv("TEST_SUBSCRIPTION_KEY")
	region := os.Getenv("TEST_SUBSCRIPTION_REGION")
	var err error
	var audioConfig *audio.AudioConfig
	audioConfig, err = audio.NewAudioConfigFromFileInput("../test_files/turn_on_the_lamp.wav")
	if err != nil {
		msg := fmt.Sprint("Got an error: ", err.Error())
		t.Error(msg)
		return
	}
	var config *BotFrameworkConfig
	config, err = NewBotFrameworkConfigFromSubscription(subscription, region)
	if err != nil {
		msg := fmt.Sprint("Got an error: ", err.Error())
		t.Error(msg)
		return
	}
	var connector *DialogServiceConnector
	connector, err = NewDialogServiceConnectorFromConfig(config, audioConfig)
	if err != nil {
		msg := fmt.Sprint("Got an error: ", err.Error())
		t.Error(msg)
		return
	}
	sessionStartedHandler := func(event speech.SessionEventArgs) {
		id := event.SessionID()
		fmt.Println("Started ", id)
	}
	sessionStoppedHandler := func(event speech.SessionEventArgs) {
		id := event.SessionID()
		fmt.Println("Stopped ", id)
	}
	connector.SessionStarted(sessionStartedHandler)
	connector.SessionStopped(sessionStoppedHandler)
	future := connector.ListenOnceAsync()
	outcome := <- future
	if outcome.Error != nil {
		msg := fmt.Sprint("Got an error: ", err.Error())
		t.Error(msg)
		return
	}
	result := outcome.Result
	fmt.Println("Recognized: ", result.Text)
}