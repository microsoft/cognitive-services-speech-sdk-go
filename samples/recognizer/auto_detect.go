// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package recognizer

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func RecognizeOnceFromAutoDetectSourceLangConfig(subscription string, region string, file string) {
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
	languageConfig, err := speech.NewAutoDetectSourceLanguageConfigFromLanguages([]string{"en-US", "de-DE"})
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer languageConfig.Close()
	speechRecognizer, err := speech.NewSpeechRecognizerFromAutoDetectSourceLangConfig(config, languageConfig, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechRecognizer.Close()
	task := speechRecognizer.RecognizeOnceAsync()
	var outcome speech.SpeechRecognitionOutcome
	select {
	case outcome = <-task:
	case <-time.After(10 * time.Second):
		fmt.Println("Timed out")
		return
	}
	defer outcome.Close()
	if outcome.Error != nil {
		fmt.Println("Got an error: ", outcome.Error)
		return
	}
	fmt.Println("Got a recognition!")
	fmt.Println(outcome.Result.Text)
}
