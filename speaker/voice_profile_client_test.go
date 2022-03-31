// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	// "fmt"
	"math/big"
	"os"
	"testing"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func createClientFromSubscriptionRegion(t *testing.T, subscription string, region string) *VoiceProfileClient {
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	defer config.Close()
	client, err := NewVoiceProfileClientFromConfig(config)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	return client
}

func createAudioConfigFromFileInput(t *testing.T, file string) *audio.AudioConfig {
	audioConfig, err := audio.NewAudioConfigFromWavFileInput(file)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	return audioConfig
}

func createClient(t *testing.T) *VoiceProfileClient {
	subscription := os.Getenv("SPEAKER_RECOGNITION_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEAKER_RECOGNITION_SUBSCRIPTION_REGION")
	return createClientFromSubscriptionRegion(t, subscription, region)
}

func createSpeakerRecognizerFromFile(t *testing.T, file string) *SpeakerRecognizer {
	subscription := os.Getenv("SPEAKER_RECOGNITION_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEAKER_RECOGNITION_SUBSCRIPTION_REGION")
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	defer config.Close()
	audioConfig := createAudioConfigFromFileInput(t, file)
	reco, err := NewSpeakerRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	return reco
}

func TestNewVoiceProfileClient(t *testing.T) {
	client := createClient(t)
	if client == nil {
		t.Error("Unexpected error: nil voice profile client")
	}
	defer client.Close()
}

func GetNewVoiceProfileFromClient(t *testing.T, client *VoiceProfileClient, expectedType common.VoiceProfileType) *VoiceProfile {
	/* Test profile creation */
	future := client.CreateProfileAsync(expectedType, "en-US")
	outcome := <-future
	if outcome.Failed() {
		t.Error("Got an error creating profile: ", outcome.Error.Error())
		return nil
	}
	profile := outcome.Profile
	id, err := profile.Id()
	if err != nil {
		t.Error("Unexpected error creating profile id: ", err)
		return nil
	}
	profileType, err := profile.Type();
	if err != nil {
		t.Error("Unexpected error getting profile type: ", err)
		return nil
	}
	if profileType != expectedType {
		t.Error("Profile type does not match expected type")
		return nil
	}
	t.Log("Profile id: ", id)
	return profile
}

func EnrollProfile(t *testing.T, client *VoiceProfileClient, profile *VoiceProfile, file string) {
	/* Test profile enrollment */
	audioConfig := createAudioConfigFromFileInput(t, file)
	defer audioConfig.Close()
	enrollmentReason, currentReason := common.EnrollingVoiceProfile, common.EnrollingVoiceProfile
	var currentResult *VoiceProfileEnrollmentResult
	expectedEnrollmentCount := 1
	for currentReason == enrollmentReason {
		enrollFuture := client.EnrollProfileAsync(profile, audioConfig)
		enrollOutcome := <-enrollFuture
		if enrollOutcome.Failed() {
			t.Error("Got an error enrolling profile: ", enrollOutcome.Error.Error())
			return
		}
		currentResult = enrollOutcome.Result
		currentReason = currentResult.Reason
		if currentResult.EnrollmentsCount != expectedEnrollmentCount {
			t.Error("Unexpected enrollments for profile: ", currentResult.RemainingEnrollmentsCount)
		}
		expectedEnrollmentCount += 1
	}
	if currentReason != common.EnrolledVoiceProfile {
		t.Error("Unexpected result enrolling profile: ", currentResult)
	}
	expectedEnrollmentsLength := big.NewInt(0)
	if currentResult.RemainingEnrollmentsLength.Int64() != expectedEnrollmentsLength.Int64() {
		t.Error("Unexpected remaining enrollment length for profile: ", currentResult.RemainingEnrollmentsLength)
	}
}

func DeleteProfile(t *testing.T, client *VoiceProfileClient, profile *VoiceProfile) {
	/* Test profile deletion */
	deleteFuture := client.DeleteProfileAsync(profile)
	deleteOutcome := <-deleteFuture
	if deleteOutcome.Failed() {
		t.Error("Got an error deleting profile: ", deleteOutcome.Error.Error())
		return
	}
	result := deleteOutcome.Result
	if result.Reason != common.DeletedVoiceProfile {
		t.Error("Unexpected result deleting profile: ", result)
	}
}

func TestVoiceProfileClientIdentification(t *testing.T) {
	client := createClient(t)
	if client == nil {
		t.Error("Unexpected error: nil voice profile client")
	}
	defer client.Close()
	expectedType := common.VoiceProfileType(1)
	
	profile := GetNewVoiceProfileFromClient(t, client, expectedType)
	if profile == nil {
		t.Error("Error creating profile")
		return
	}
	defer profile.Close()

	/* Test profile reset */
	resetFuture := client.ResetProfileAsync(profile)
	resetOutcome := <-resetFuture
	if resetOutcome.Failed() {
		t.Error("Got an error resetting profile: ", resetOutcome.Error.Error())
	}
	result := resetOutcome.Result
	if result.Reason != common.ResetVoiceProfile {
		t.Error("Unexpected result resetting profile: ", result)
	}
    
	EnrollProfile(t, client, profile, "../test_files/TalkForAFewSeconds16.wav")

	/* Test identification */
	profiles := []*VoiceProfile{profile}
	model, err := NewSpeakerIdentificationModelFromProfiles(profiles)
	if err != nil {
		t.Error("Error creating Identification model: ", err)
	}
	if model == nil {
		t.Error("Error creating Identification model: nil model")
		return
	}
	speakerRecognizer := createSpeakerRecognizerFromFile(t, "../test_files/TalkForAFewSeconds16.wav")
	identifyFuture := speakerRecognizer.IdentifyOnceAsync(model)
	identifyOutcome := <-identifyFuture
	if identifyOutcome.Failed() {
		t.Error("Got an error identifying profile: ", identifyOutcome.Error.Error())
		return
	}
	identifyResult := identifyOutcome.Result
	if identifyResult.Reason != common.RecognizedSpeakers {
		t.Error("Got an unexpected result identifying profile: ", identifyResult)
	}
	expectedID, _ := profile.Id()
	if identifyResult.ProfileID != expectedID {
		t.Error("Got an unexpected profile id identifying profile: ", identifyResult.ProfileID)
	}
	if identifyResult.Score < 1.0 {
		t.Error("Got an unexpected score identifying profile: ", identifyResult.Score)
	}

	DeleteProfile(t, client, profile)
}

func TestVoiceProfileClientIndependentVerification(t *testing.T) {
	client := createClient(t)
	if client == nil {
		t.Error("Unexpected error: nil voice profile client")
	}
	defer client.Close()
	expectedType := common.VoiceProfileType(3)
	
	profile := GetNewVoiceProfileFromClient(t, client, expectedType)
	if profile == nil {
		t.Error("Error creating profile")
		return
	}
	defer profile.Close()

	EnrollProfile(t, client, profile, "../test_files/TalkForAFewSeconds16.wav")

	/* Test enrollment result */
	enrollFuture := client.RetrieveEnrollmentResultAsync(profile)
	enrollOutcome := <-enrollFuture
	if enrollOutcome.Failed() {
		t.Error("Got an error enrolling profile: ", enrollOutcome.Error.Error())
		return
	}
	enrollResult := enrollOutcome.Result
	enrollReason := enrollResult.Reason
	if enrollReason != common.EnrolledVoiceProfile {
		t.Error("Unexpected result enrolling profile: ", enrollResult)
	}
	expectedEnrollmentsLength := big.NewInt(0)
	if enrollResult.RemainingEnrollmentsLength.Int64() != expectedEnrollmentsLength.Int64() {
		t.Error("Unexpected remaining enrollment length for profile: ", enrollResult.RemainingEnrollmentsLength)
	}

	/* Test verification */
	model, err := NewSpeakerVerificationModelFromProfile(profile)
	if err != nil {
		t.Error("Error creating Verification model: ", err)
	}
	if model == nil {
		t.Error("Error creating Verification model: nil model")
		return
	}
	speakerRecognizer := createSpeakerRecognizerFromFile(t, "../test_files/TalkForAFewSeconds16.wav")
	verifyFuture := speakerRecognizer.VerifyOnceAsync(model)
	verifyOutcome := <-verifyFuture
	if verifyOutcome.Failed() {
		t.Error("Got an error verifying profile: ", verifyOutcome.Error.Error())
		return
	}
	verifyResult := verifyOutcome.Result
	if verifyResult.Reason != common.RecognizedSpeaker {
		t.Error("Got an unexpected result verifying profile: ", verifyResult)
	}
	expectedID, _ := profile.Id()
	if verifyResult.ProfileID != expectedID {
		t.Error("Got an unexpected profile id verifying profile: ", verifyResult.ProfileID)
	}
	if verifyResult.Score < 1.0 {
		t.Error("Got an unexpected score verifying profile: ", verifyResult.Score)
	}

	DeleteProfile(t, client, profile)
}

func TestVoiceProfileClientDependentVerification(t *testing.T) {
	client := createClient(t)
	if client == nil {
		t.Error("Unexpected error: nil voice profile client")
	}
	defer client.Close()
	expectedType := common.VoiceProfileType(2)
	
	profile := GetNewVoiceProfileFromClient(t, client, expectedType)
	if profile == nil {
		t.Error("Error creating profile")
		return
	}
	defer profile.Close()

	EnrollProfile(t, client, profile, "../test_files/myVoiceIsMyPassportVerifyMe01.wav")

	/* Test verification */
	model, err := NewSpeakerVerificationModelFromProfile(profile)
	if err != nil {
		t.Error("Error creating Verification model: ", err)
	}
	if model == nil {
		t.Error("Error creating Verification model: nil model")
		return
	}
	speakerRecognizer := createSpeakerRecognizerFromFile(t, "../test_files/myVoiceIsMyPassportVerifyMe01.wav")
	verifyFuture := speakerRecognizer.VerifyOnceAsync(model)
	verifyOutcome := <-verifyFuture
	if verifyOutcome.Failed() {
		t.Error("Got an error verifying profile: ", verifyOutcome.Error.Error())
		return
	}
	verifyResult := verifyOutcome.Result
	if verifyResult.Reason != common.RecognizedSpeaker {
		t.Error("Got an unexpected result verifying profile: ", verifyResult)
	}
	expectedID, _ := profile.Id()
	if verifyResult.ProfileID != expectedID {
		t.Error("Got an unexpected profile id verifying profile: ", verifyResult.ProfileID)
	}
	if verifyResult.Score < 1.0 {
		t.Error("Got an unexpected score verifying profile: ", verifyResult.Score)
	}

	DeleteProfile(t, client, profile)
}

func TestGetActivationPhrases(t *testing.T) {
	client := createClient(t)
	if client == nil {
		t.Error("Unexpected error: nil voice profile client")
	}
	defer client.Close()
	expectedType := common.VoiceProfileType(2)
	future := client.GetActivationPhrasesAsync(expectedType, "en-US")
	outcome := <-future
	if outcome.Failed() {
		t.Error("Got an error getting activation phrases: ", outcome.Error.Error())
		return
	}
	result := outcome.Result
	defer result.Close()
	phrases := result.Phrases
	if len(phrases) < 1 {
		t.Error("Unexpected error getting phrases, no phrases received")
	}
	for _, phrase := range phrases {
		t.Log("Phrase received: ", phrase)
	}
}

func TestGetAllProfiles(t *testing.T) {
	client := createClient(t)
	if client == nil {
		t.Error("Unexpected error: nil voice profile client")
	}
	defer client.Close()
	expectedType := common.VoiceProfileType(2)
	
	profile := GetNewVoiceProfileFromClient(t, client, expectedType)
	if profile == nil {
		t.Error("Error creating profile")
		return
	}
	defer profile.Close()

	expectedID, _ := profile.Id()
	profileType, _ := profile.Type()

	future := client.GetAllProfilesAsync(profileType)
	outcome := <-future
	if outcome.Failed() {
		t.Error("Error getting all profiles: ", outcome.Error.Error())
		return
	}
	profiles := outcome.Profiles
	if len(profiles) < 1 {
		t.Error("Unexpected error getting profiles, no profiles received")
	}
	profileFound := false
	for _, p := range profiles {
		id, _ := p.Id()
		t.Log("Profile id in list: ", id)
		
		if id == expectedID {
			profileFound = true
		} else {
			p.Close() // Not closing all unused profiles may produce memory issues
		}
	}

	if !profileFound {
		t.Error("Unexpected error getting profiles, added profile not found")
	}

	DeleteProfile(t, client, profile)
}