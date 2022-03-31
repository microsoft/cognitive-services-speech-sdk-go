// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker_recognition

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speaker"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func IndependentVerification(subscription string, region string, file string) {
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer config.Close()
	client, err := speaker.NewVoiceProfileClientFromConfig(config)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer client.Close()
	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer audioConfig.Close()
	<-time.After(10 * time.Second)
	expectedType := common.VoiceProfileType(3)
	
	profile := GetNewVoiceProfileFromClient(client, expectedType)
	if profile == nil {
		fmt.Println("Error creating profile")
		return
	}
	defer profile.Close()

	EnrollProfile(client, profile, audioConfig)

	model, err := speaker.NewSpeakerVerificationModelFromProfile(profile)
	if err != nil {
		fmt.Println("Error creating Verification model: ", err)
	}
	if model == nil {
		fmt.Println("Error creating Verification model: nil model")
		return
	}
	verifyAudioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer verifyAudioConfig.Close()
	speakerRecognizer, err := speaker.NewSpeakerRecognizerFromConfig(config, verifyAudioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	verifyFuture := speakerRecognizer.VerifyOnceAsync(model)
	verifyOutcome := <-verifyFuture
	if verifyOutcome.Failed() {
		fmt.Println("Got an error verifying profile: ", verifyOutcome.Error.Error())
		return
	}
	verifyResult := verifyOutcome.Result
	if verifyResult.Reason != common.RecognizedSpeaker {
		fmt.Println("Got an unexpected result verifying profile: ", verifyResult)
	}
	expectedID, _ := profile.Id()
	if verifyResult.ProfileID != expectedID {
		fmt.Println("Got an unexpected profile id verifying profile: ", verifyResult.ProfileID)
	}
	if verifyResult.Score < 1.0 {
		fmt.Println("Got an unexpected score verifying profile: ", verifyResult.Score)
	}

	DeleteProfile(client, profile)
}
