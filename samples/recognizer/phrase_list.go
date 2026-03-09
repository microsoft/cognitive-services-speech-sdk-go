// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package recognizer

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func RecognizeOnceWithPhraseList(subscription string, region string, file string) {
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

	// Create a phrase list grammar to improve recognition accuracy
	phraseListGrammar, err := speech.NewPhraseListGrammarFromRecognizer(speechRecognizer)
	if err != nil {
		fmt.Println("Got an error creating phrase list grammar: ", err)
		return
	}
	defer phraseListGrammar.Close()

	// Add phrases that are likely to be spoken
	phraseListGrammar.AddPhrase("Contoso")
	phraseListGrammar.AddPhrase("Cognito")
	phraseListGrammar.AddPhrase("peloozoid")

	// Set the weight to bias recognition towards the phrases in the list.
	// Valid range is 0.0 to 2.0. Default is 1.0. Higher values increase bias.
	err = phraseListGrammar.SetWeight(2.0)
	if err != nil {
		fmt.Println("Got an error setting phrase list weight: ", err)
		return
	}

	fmt.Println("Phrase list grammar configured with custom weight.")

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
