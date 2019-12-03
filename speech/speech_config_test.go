//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package speech

import (
	"testing"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

func TestFromSubscription(t *testing.T) {
	subscription := "test"
	region := "region"
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.SubscriptionKey() != subscription {
		t.Error("Subscription not properly set")
	}
	if config.Region() != region {
		t.Error("Region not properly set")
	}
}

func TestFromAuthorizationToken(t *testing.T) {
	auth := "test"
	region := "region"
	config, err := NewSpeechConfigFromAuthorizationToken(auth, region)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.AuthorizationToken() != auth {
		t.Error("Authorization Token not properly set")
	}
	if config.Region() != region {
		t.Error("Region not properly set")
	}
}

func TestPropertiesByID(t *testing.T) {
	subscription := "test"
	region := "region"
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Unexpected error")
	}
	value := "value1"
	err = config.SetProperty(common.SpeechServiceConnectionKey, value)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.GetProperty(common.SpeechServiceConnectionKey) != value {
		t.Error("Propery value not valid")
	}
}

func TestPropertiesByString(t *testing.T) {
	subscription := "test"
	region := "region"
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Unexpected error")
	}
	value := "value1"
	err = config.SetPropertyByString("key1", value)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.GetPropertyByString("key1") != value {
		t.Error("Propery value not valid")
	}

}