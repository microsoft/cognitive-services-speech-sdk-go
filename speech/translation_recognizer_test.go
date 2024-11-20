package speech

import (
	"os"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
)

func TestStartContinuousRecognitionAsync(t *testing.T) {
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
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer audioConfig.Close()
	config, err := NewSpeechTranslationConfigFromSubscription(os.Getenv("SPEECH_SUBSCRIPTION_KEY"), os.Getenv("SPEECH_SUBSCRIPTION_REGION"))
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	defer config.Close()
	err = config.SetSpeechRecognitionLanguage("en-US")
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	err = config.AddTargetLanguage("en")
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	err = config.SetTranslationVoiceName("en-US-AndrewMultilingualNeural")
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	recognizer, err := NewTranslationRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Error("Got an error ", err.Error())
	}
	if recognizer == nil {
		return
	}
	defer recognizer.Close()
	firstResult := true
	recognizedFuture := make(chan string)
	recognizingFuture := make(chan string)
	synthesizingFuture := make(chan string)
	recognizedHandler := func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		firstResult = true
		t.Log("Recognized: ", event.Result.Text)
		recognizedFuture <- "Recognized"
	}
	recognizingHandle := func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		t.Log("Recognizing: ", event.Result.Text)
		if firstResult {
			firstResult = false
			recognizingFuture <- "Recognizing"
		}
	}
	synthesizingHandle := func(event TranslationSynthesisEventArgs) {
		defer event.Close()
		t.Log("Synthesizing: ", len(event.Result.AudioData))
		synthesizingFuture <- "Synthesizing"
	}
	recognizer.Recognized(recognizedHandler)
	recognizer.Recognizing(recognizingHandle)
	recognizer.Synthesizing(synthesizingHandle)
	err = <-recognizer.StartContinuousRecognitionAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
	pumpFileIntoStream(t, "../test_files/turn_on_the_lamp.wav", stream)
	pumpFileIntoStream(t, "../test_files/turn_on_the_lamp.wav", stream)
	pumpSilenceIntoStream(t, stream)
	stream.CloseStream()
	select {
	case <-recognizingFuture:
		t.Log("Received first Recognizing event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive first Recognizing event.")
	}
	select {
	case <-recognizedFuture:
		t.Log("Received first Recognized event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive first Recognized event.")
	}
	select {
	case <-synthesizingFuture:
		t.Log("Received first Synthesizing event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive first Synthesizing event.")
	}
	select {
	case <-recognizingFuture:
		t.Log("Received second Recognizing event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive second Recognizing event.")
	}
	select {
	case <-recognizedFuture:
		t.Log("Received second Recognized event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive second Recognized event.")
	}
	select {
	case <-synthesizingFuture:
		t.Log("Received second Synthesizing event.")
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive second Synthesizing event.")
	}
	err = <-recognizer.StopContinuousRecognitionAsync()
	if err != nil {
		t.Error("Got error: ", err)
	}
}
