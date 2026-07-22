// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package translation

import (
	"fmt"
	"os"
	"sync"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// LiveInterpreterFromFile demonstrates real-time speech-to-speech translation using the Live
// Interpreter API with personal voice. Live Interpreter continuously identifies the language being
// spoken (no source language required) and delivers low-latency translated speech in a natural voice
// that preserves the speaker's style and tone.
//
// Prerequisites:
//   - Personal voice access must be granted for your Speech resource (a Limited Access feature).
//     Apply at https://aka.ms/customneural and select "Personal Voice" for Question 20.
//   - The Speech resource must be in a region that supports Live Interpreter. See the Speech service
//     regions table for current regional availability.
//
// Live Interpreter requires the universal v2 endpoint together with open-range language
// identification, so the config is created with FromEndpoint (not FromSubscription).
func LiveInterpreterFromFile(subscription string, region string, file string) {
	// When you use multilingual translation with language identification you must use the v2 endpoint
	// and create the SpeechTranslationConfig with FromEndpoint.
	//   Regional form (used here):  wss://<region>.stt.speech.microsoft.com/speech/universal/v2
	//   Resource form (from docs):  wss://<YourResourceName>.cognitiveservices.azure.com/stt/speech/universal/v2
	endpoint := fmt.Sprintf("wss://%s.stt.speech.microsoft.com/speech/universal/v2", region)

	config, err := speech.NewSpeechTranslationConfigFromEndpointWithSubscription(endpoint, subscription)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer config.Close()

	// Set the translation target language and enable personal voice.
	targetLanguage := "zh-Hans"
	if err := config.AddTargetLanguage(targetLanguage); err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	if err := config.SetVoiceName("personal-voice"); err != nil {
		fmt.Println("Got an error: ", err)
		return
	}

	// You don't need to define any candidate languages to detect: open range enables continuous
	// language identification across the full set of supported languages.
	autoDetectSourceLanguageConfig, err := speech.NewAutoDetectSourceLanguageConfigFromOpenRange()
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer autoDetectSourceLanguageConfig.Close()

	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer audioConfig.Close()

	recognizer, err := speech.NewTranslationRecognizerFromAutoDetectSourceLangConfig(config, autoDetectSourceLanguageConfig, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer recognizer.Close()

	// stop is closed exactly once when translation completes (cancellation or session stop).
	stop := make(chan struct{})
	var stopOnce sync.Once
	signalStop := func() { stopOnce.Do(func() { close(stop) }) }

	// Index of the output audio files produced by the synthesizer.
	audioIndex := 0

	recognizer.Recognizing(func(event speech.TranslationRecognitionEventArgs) {
		defer event.Close()
		detected := event.Result.Properties.GetProperty(common.SpeechServiceConnectionAutoDetectSourceLanguageResult, "")
		fmt.Printf("RECOGNIZING in '%s': Text=%s\n", detected, event.Result.Text)
		for lang, translation := range event.Result.GetTranslations() {
			fmt.Printf("    TRANSLATING into '%s': %s\n", lang, translation)
		}
	})

	recognizer.Recognized(func(event speech.TranslationRecognitionEventArgs) {
		defer event.Close()
		switch event.Result.Reason {
		case common.TranslatedSpeech:
			detected := event.Result.Properties.GetProperty(common.SpeechServiceConnectionAutoDetectSourceLanguageResult, "")
			fmt.Printf("RECOGNIZED in '%s': Text=%s\n", detected, event.Result.Text)
			for lang, translation := range event.Result.GetTranslations() {
				fmt.Printf("    TRANSLATED into '%s': %s\n", lang, translation)
			}
		case common.RecognizedSpeech:
			fmt.Printf("RECOGNIZED: Text=%s\n", event.Result.Text)
			fmt.Println("    Speech not translated.")
		case common.NoMatch:
			fmt.Println("NOMATCH: Speech could not be recognized.")
		}
	})

	recognizer.Synthesizing(func(event speech.TranslationSynthesisEventArgs) {
		defer event.Close()
		audioData := event.Result.GetAudioData()
		if len(audioData) == 0 {
			return
		}
		audioIndex++
		outputFile := fmt.Sprintf("live-interpreter-%d.wav", audioIndex)
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating file: ", err)
			return
		}
		defer f.Close()
		if _, err := f.Write(audioData); err != nil {
			fmt.Println("Error writing audio: ", err)
			return
		}
		fmt.Printf("Audio synthesized: %d byte(s) -> %s\n", len(audioData), outputFile)
	})

	recognizer.Canceled(func(event speech.TranslationRecognitionCanceledEventArgs) {
		defer event.Close()
		fmt.Printf("CANCELED: Reason=%d\n", event.Reason)
		if event.Reason == common.Error {
			fmt.Printf("CANCELED: ErrorCode=%d\n", event.ErrorCode)
			fmt.Printf("CANCELED: ErrorDetails=%s\n", event.ErrorDetails)
			fmt.Println("CANCELED: Did you set the resource key/region and enable personal voice access?")
		}
		signalStop()
	})

	recognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session stopped.")
		signalStop()
	})

	fmt.Println("Start translation...")
	if err := <-recognizer.StartContinuousRecognitionAsync(); err != nil {
		fmt.Println("Got an error: ", err)
		return
	}

	// Wait for the translation to complete (end of stream cancellation or session stop).
	<-stop

	if err := <-recognizer.StopContinuousRecognitionAsync(); err != nil {
		fmt.Println("Got an error: ", err)
	}
}
