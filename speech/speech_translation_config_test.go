package speech

import (
	"testing"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

func TestNewSpeechTranslationConfigFromSubscription(t *testing.T) {
	subscription := "test"
	region := "region"
	config, err := NewSpeechTranslationConfigFromSubscription(subscription, region)
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

func TestNewSpeechTranslationConfigFromEndpointWithSubscription(t *testing.T) {
	subscription := "test"
	config, err := NewSpeechTranslationConfigFromEndpointWithSubscription("endpoint", subscription)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.SubscriptionKey() != subscription {
		t.Error("Subscription not properly set")
	}
}

func TestAddTargetLanguage(t *testing.T) {
	subscription := "test"
	region := "region"
	config, err := NewSpeechTranslationConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Unexpected error")
	}
	language := "en"
	err = config.AddTargetLanguage(language)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.GetProperty(common.SpeechServiceConnectionTranslationToLanguages) != language {
		t.Error("Property value not valid")
	}
}

func TestSetSpeechRecognitionLanguage(t *testing.T) {
	subscription := "test"
	region := "region"
	config, err := NewSpeechTranslationConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Unexpected error")
	}
	language := "en-US"
	err = config.SetSpeechRecognitionLanguage(language)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.SpeechRecognitionLanguage() != language {
		t.Error("Property value not valid")
	}
}

func TestSetTranslationVoiceName(t *testing.T) {
	subscription := "test"
	region := "region"
	config, err := NewSpeechTranslationConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Unexpected error")
	}
	voiceName := "voiceName"
	err = config.SetTranslationVoiceName(voiceName)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.TranslationVoiceName() != voiceName {
		t.Error("Property value not valid")
	}
}
