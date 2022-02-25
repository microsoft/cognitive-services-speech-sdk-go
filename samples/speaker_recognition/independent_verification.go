// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker_recognition

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speaker"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/samples/helpers"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func GetNewVoiceProfileFromClient(client *speaker.VoiceProfileClient, expectedType common.VoiceProfileType) *speaker.VoiceProfile {
	future := client.CreateProfileAsync(expectedType, "en-US")
	outcome := <-future
	if outcome.Failed() {
		fmt.Println("Got an error creating profile: ", outcome.Error.Error())
		return nil
	}
	profile := outcome.profile
	id, err := profile.Id()
	if err != nil {
		fmt.Println("Unexpected error creating profile id: ", err)
		return nil
	}
	profileType, err := profile.Type();
	if err != nil {
		fmt.Println("Unexpected error getting profile type: ", err)
		return nil
	}
	if profileType != expectedType {
		fmt.Println("Profile type does not match expected type")
		return nil
	}
	return profile
}

func EnrollProfile(client *speaker.VoiceProfileClient, profile *speaker.VoiceProfile, audioConfig audio.AudioConfig) {
	enrollmentReason, currentReason := common.EnrollingVoiceProfile, common.EnrollingVoiceProfile
	var currentResult *VoiceProfileEnrollmentResult
	expectedEnrollmentCount := 1
	for currentReason == enrollmentReason {
		enrollFuture := client.EnrollProfileAsync(profile, audioConfig)
		enrollOutcome := <-enrollFuture
		if enrollOutcome.Failed() {
			fmt.Println("Got an error enrolling profile: ", enrollOutcome.Error.Error())
			return
		}
		currentResult = enrollOutcome.Result
		currentReason = currentResult.Reason
		if currentResult.EnrollmentsCount != expectedEnrollmentCount {
			fmt.Println("Unexpected enrollments for profile: ", currentResult.RemainingEnrollmentsCount)
		}
		expectedEnrollmentCount += 1
	}
	if currentReason != common.EnrolledVoiceProfile {
		fmt.Println("Unexpected result enrolling profile: ", currentResult)
	}
	expectedEnrollmentsLength := big.NewInt(0)
	if currentResult.RemainingEnrollmentsLength.Int64() != expectedEnrollmentsLength.Int64() {
		fmt.Println("Unexpected remaining enrollment length for profile: ", currentResult.RemainingEnrollmentsLength)
	}
}

func DeleteProfile(client *speaker.VoiceProfileClient, profile *speaker.VoiceProfile) {
	deleteFuture := client.DeleteProfileAsync(profile)
	deleteOutcome := <-deleteFuture
	if deleteOutcome.Failed() {
		fmt.Println("Got an error deleting profile: ", deleteOutcome.Error.Error())
		return
	}
	result := deleteOutcome.Result
	if result.Reason != common.DeletedVoiceProfile {
		fmt.Println("Unexpected result deleting profile: ", result)
	}
}

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
		return nil
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
