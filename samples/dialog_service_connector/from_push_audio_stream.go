package main

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/dialog"
	"fmt"
	"os"
	"bufio"
	"io"
	"time"
)

func pumpFileIntoStream(filename string, stream *audio.PushAudioInputStream) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err);
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1000)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			fmt.Println("Done reading file.")
			break
		}
		if err != nil {
			fmt.Println("Error reading file: ", err)
			break
		}
		err = stream.Write(buffer[0:n])
		if err != nil {
			fmt.Println("Error writing to the stream")
		}
	}
}


func main() {
	args := os.Args[1:]
	if (len(args) != 3) {
		fmt.Println("Input not valid")
		fmt.Println("Usage: ")
		fmt.Println(os.Args[0], " <subscription> <region> <file>")
		return
	}
	subscription := args[0]
	region := args[1]
	file := args[2]
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
		stream.CloseStream()
	}
	connector.Recognized(recognizedHandle)
	recognizingHandler := func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		fmt.Println("Recognizing ", event.Result.Text)
	}
	connector.Recognizing(recognizingHandler)
	pumpFileIntoStream(file, stream)
	connector.ListenOnceAsync()
	<- time.After(10 * time.Second)
}