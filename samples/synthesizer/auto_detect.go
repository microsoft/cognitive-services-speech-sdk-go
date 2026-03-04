// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package synthesizer

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func SynthesisFromAutoDetectSourceLangConfig(subscription string, region string, file string) {
	_ = file
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer config.Close()
	languageConfig, err := speech.NewAutoDetectSourceLanguageConfigFromOpenRange()
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer languageConfig.Close()
	speechSynthesizer, err := speech.NewSpeechSynthesizerFromAutoDetectSourceLangConfig(config, languageConfig, nil)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechSynthesizer.Close()
	task := speechSynthesizer.SpeakTextAsync("Hello, world.")
	var outcome speech.SpeechSynthesisOutcome
	select {
	case outcome = <-task:
	case <-time.After(30 * time.Second):
		fmt.Println("Timed out")
		return
	}
	defer outcome.Close()
	if outcome.Error != nil {
		fmt.Println("Got an error: ", outcome.Error)
		return
	}
	if outcome.Result.Reason != common.SynthesizingAudioCompleted {
		fmt.Println("Synthesis failed with reason:", outcome.Result.Reason)
		return
	}
	fmt.Println("Speech synthesized successfully with auto-detect source language config.")
}
