// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package recognizer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

type SDKWrapperEventType int

const (
	Cancellation SDKWrapperEventType = iota
	Recognizing
	Recognized
)

type SDKWrapperEvent struct {
	EventType    SDKWrapperEventType
	Cancellation *speech.SpeechRecognitionCanceledEventArgs
	Recognized   *speech.SpeechRecognitionEventArgs
	Recognizing  *speech.SpeechRecognitionEventArgs
}

func (event *SDKWrapperEvent) Close() {
	if event.Cancellation != nil {
		event.Cancellation.Close()
	}
	if event.Recognizing != nil {
		event.Recognizing.Close()
	}
	if event.Recognized != nil {
		event.Recognized.Close()
	}
}

type SDKWrapper struct {
	stream     *audio.PushAudioInputStream
	recognizer *speech.SpeechRecognizer
	started    int32
}

func NewWrapper(subscription string, region string) (*SDKWrapper, error) {
	format, err := audio.GetDefaultInputFormat()
	if err != nil {
		return nil, err
	}
	defer format.Close()
	stream, err := audio.CreatePushAudioInputStreamFromFormat(format)
	if err != nil {
		return nil, err
	}
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		stream.Close()
		return nil, err
	}
	defer audioConfig.Close()
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		stream.Close()
		return nil, err
	}
	defer config.Close()
	recognizer, err := speech.NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		stream.Close()
		return nil, err
	}
	wrapper := new(SDKWrapper)
	wrapper.recognizer = recognizer
	wrapper.stream = stream
	return wrapper, nil
}

func (wrapper *SDKWrapper) Close() {
	wrapper.stream.CloseStream()
	<-wrapper.recognizer.StopContinuousRecognitionAsync()
	wrapper.stream.Close()
	wrapper.recognizer.Close()
}

func (wrapper *SDKWrapper) Write(buffer []byte) error {
	if atomic.LoadInt32(&wrapper.started) != 1 {
		return fmt.Errorf("Trying to write when recognizer is stopped")
	}
	return wrapper.stream.Write(buffer)
}

func (wrapper *SDKWrapper) StartContinuous(callback func(*SDKWrapperEvent)) error {
	if atomic.SwapInt32(&wrapper.started, 1) == 1 {
		return nil
	}
	wrapper.recognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		wrapperEvent := new(SDKWrapperEvent)
		wrapperEvent.EventType = Recognized
		wrapperEvent.Recognized = &event
		callback(wrapperEvent)
	})
	wrapper.recognizer.Recognizing(func(event speech.SpeechRecognitionEventArgs) {
		wrapperEvent := new(SDKWrapperEvent)
		wrapperEvent.EventType = Recognizing
		wrapperEvent.Recognizing = &event
		callback(wrapperEvent)
	})
	wrapper.recognizer.Canceled(func(event speech.SpeechRecognitionCanceledEventArgs) {
		wrapperEvent := new(SDKWrapperEvent)
		wrapperEvent.EventType = Cancellation
		wrapperEvent.Cancellation = &event
		callback(wrapperEvent)
	})
	return <-wrapper.recognizer.StartContinuousRecognitionAsync()
}

func (wrapper *SDKWrapper) StopContinuous() error {
	if atomic.SwapInt32(&wrapper.started, 0) == 0 {
		return nil
	}
	var empty = []byte{}
	wrapper.stream.Write(empty)
	wrapper.recognizer.Recognized(nil)
	wrapper.recognizer.Recognizing(nil)
	wrapper.recognizer.Canceled(nil)
	return <-wrapper.recognizer.StopContinuousRecognitionAsync()
}

func PumpFileContinuously(stop chan int, filename string, wrapper *SDKWrapper) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 3200)
	for {
		select {
		case <-stop:
			fmt.Println("Stopping pump...")
			return
		case <-time.After(100 * time.Millisecond):
		}
		n, err := reader.Read(buffer)
		if err == io.EOF {
			file.Seek(44, io.SeekStart)
			continue
		}
		if err != nil {
			fmt.Println("Error reading file: ", err)
			break
		}
		err = wrapper.Write(buffer[0:n])
		if err != nil {
			fmt.Println("Error writing to the stream")
		}
	}
}

func RecognizeContinuousUsingWrapper(subscription string, region string, file string) {
	/* If running this in a server, each worker thread should run something similar to this */
	wrapper, err := NewWrapper(subscription, region)
	if err != nil {
		fmt.Println("Got an error: ", err)
	}
	defer wrapper.Close()
	stop := make(chan int)
	go PumpFileContinuously(stop, file, wrapper)
	fmt.Println("Starting Continuous...")
	wrapper.StartContinuous(func(event *SDKWrapperEvent) {
		defer event.Close()
		switch event.EventType {
		case Recognized:
			fmt.Println("Got a recognized event")
		case Recognizing:
			fmt.Println("Got a recognizing event")
		case Cancellation:
			fmt.Println("Got a cancellation event")
		}
	})
	<-time.After(10 * time.Second)
	stop <- 1
	fmt.Println("Stopping Continuous...")
	wrapper.StopContinuous()
	fmt.Println("Exiting...")
}
