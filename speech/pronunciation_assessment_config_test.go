package speech

import (
	"testing"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

func TestPronunciationAssessmentConfig(t *testing.T) {
	config, err := NewPronunciationAssessmentConfig(
		"",
		common.PronunciationAssessmentGradingSystemHundredMark,
		common.PronunciationAssessmentGranularityPhoneme,
		false,
	)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.GetReferenceText() != "" {
		t.Error("Reference text not valid")
	}
	t.Log(config.String())
}

func TestPronunciationAssessmentFromJson(t *testing.T) {
	config, err := NewPronunciationAssessmentConfigFromJson(`{"dimension":"Comprehensive","referenceText":"","gradingSystem":"HundredMark","granularity":"Phoneme"}`)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.GetReferenceText() != "" {
		t.Error("Reference text not valid")
	}
	t.Log(config.String())
}

func TestPronunciationAssessmentPropertiesByID(t *testing.T) {
	config, err := NewPronunciationAssessmentConfig(
		"",
		common.PronunciationAssessmentGradingSystemFivePoint,
		common.PronunciationAssessmentGranularityWord,
		false,
	)
	if err != nil {
		t.Error("Unexpected error")
	}
	value := "value1"
	err = config.SetProperty(common.PronunciationAssessmentReferenceText, value)
	if err != nil {
		t.Error("Unexpected error")
	}
	if config.GetProperty(common.PronunciationAssessmentReferenceText) != value {
		t.Error("Propery value not valid")
	}
	t.Log(config.String())
}

func TestPronunciationAssessmentPropertiesByString(t *testing.T) {
	config, err := NewPronunciationAssessmentConfig(
		"",
		common.PronunciationAssessmentGradingSystemHundredMark,
		common.PronunciationAssessmentGranularityFullText,
		false,
	)
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
	t.Log(config.String())
}
