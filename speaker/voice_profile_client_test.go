// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	// "fmt"
	"os"
	"testing"
	//"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
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

func createClient(t *testing.T) *VoiceProfileClient {
	subscription := os.Getenv("SPEAKER_RECOGNITION_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEAKER_RECOGNITION_SUBSCRIPTION_REGION")
	return createClientFromSubscriptionRegion(t, subscription, region)
}

func TestNewVoiceProfileClient(t *testing.T) {
	client := createClient(t)
	if client == nil {
		t.Error("Unexpected error: nil voice profile client")
	}
	defer client.Close()
}

func TestVoiceProfileClientCreateProfile(t *testing.T) {
	client := createClient(t)
	if client == nil {
		t.Error("Unexpected error: nil voice profile client")
	}
	defer client.Close()
	expectedType := common.VoiceProfileType(2)
	future := client.CreateProfileAsync(expectedType, "en-US")
	outcome := <-future
	if outcome.Failed() {
		t.Error("Got an error: ", outcome.Error.Error())
		return
	}
	profile := outcome.profile
	defer profile.Close()
	id, err := profile.Id()
	if err != nil {
		t.Error("Unexpected error creating profile id: ", err)
	}
	profileType, err := profile.Type();
	if err != nil {
		t.Error("Unexpected error getting profile type: ", err)
	}
	if profileType != expectedType {
		t.Error("Profile type does not match expected type")
	}
	t.Log("Profile id: ", id)
}