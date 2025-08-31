package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_pronunciation_assessment_config.h>
import "C"

type PronunciationAssessmentConfig struct {
	handle     C.SPXHANDLE
	properties *common.PropertyCollection
}

// NewPronunciationAssessmentConfigFromHandle creates a PronunciationAssessmentConfig instance from a valid handle. This is for internal use only.
func NewPronunciationAssessmentConfigFromHandle(handle common.SPXHandle) (*PronunciationAssessmentConfig, error) {
	var cHandle = uintptr2handle(handle)
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	ret := uintptr(C.pronunciation_assessment_config_get_property_bag(cHandle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(PronunciationAssessmentConfig)
	config.handle = cHandle
	config.properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return config, nil
}

func NewPronunciationAssessmentConfig(referenceText string, gradingSystem common.PronunciationAssessment_GradingSystem, granularity common.PronunciationAssessment_Granularity, enableMiscue bool) (*PronunciationAssessmentConfig, error) {
	var handle C.SPXHANDLE
	rt := C.CString(referenceText)
	defer C.free(unsafe.Pointer(rt))
	ret := uintptr(C.create_pronunciation_assessment_config(&handle, rt, (C.Pronunciation_Assessment_Grading_System)(gradingSystem), (C.Pronunciation_Assessment_Granularity)(granularity), (C.bool)(enableMiscue)))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return NewPronunciationAssessmentConfigFromHandle(handle2uintptr(handle))
}

func NewPronunciationAssessmentConfigFromJson(config string) (*PronunciationAssessmentConfig, error) {
	var handle C.SPXHANDLE
	cfg := C.CString(config)
	defer C.free(unsafe.Pointer(cfg))
	ret := uintptr(C.create_pronunciation_assessment_config_from_json(&handle, cfg))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return NewPronunciationAssessmentConfigFromHandle(handle2uintptr(handle))
}

func (config *PronunciationAssessmentConfig) GetReferenceText() string {
	return config.GetProperty(common.PronunciationAssessmentReferenceText)
}

func (config *PronunciationAssessmentConfig) SetReferenceText(referenceText string) error {
	return config.SetProperty(common.PronunciationAssessmentReferenceText, referenceText)
}

func (config *PronunciationAssessmentConfig) SetPhonemeAlphabet(phonemeAlphabet string) error {
	return config.SetProperty(common.PronunciationAssessmentPhonemeAlphabet, phonemeAlphabet)
}

func (config *PronunciationAssessmentConfig) SetNBestPhonemeCount(count string) error {
	return config.SetProperty(common.PronunciationAssessmentNBestPhonemeCount, count)
}

func (config *PronunciationAssessmentConfig) EnableProsodyAssessment() error {
	return config.SetProperty(common.PronunciationAssessmentEnableProsodyAssessment, "true")
}

func (config *PronunciationAssessmentConfig) EnableContentAssessmentWithTopic(contentTopic string) error {
	return config.SetProperty(common.PronunciationAssessmentContentTopic, contentTopic)
}

func (config *PronunciationAssessmentConfig) ApplyTo(recognizer *SpeechRecognizer) error {
	ret := uintptr(C.pronunciation_assessment_config_apply_to_recognizer(config.handle, recognizer.handle))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// SetProperty sets a property value by ID.
func (config *PronunciationAssessmentConfig) SetProperty(id common.PropertyID, value string) error {
	return config.properties.SetProperty(id, value)
}

// GetProperty gets a property value by ID.
func (config *PronunciationAssessmentConfig) GetProperty(id common.PropertyID) string {
	return config.properties.GetProperty(id, "")
}

// SetPropertyByString sets a property value by string.
func (config *PronunciationAssessmentConfig) SetPropertyByString(name string, value string) error {
	return config.properties.SetPropertyByString(name, value)
}

// GetPropertyByString gets a property value by string.
func (config *PronunciationAssessmentConfig) GetPropertyByString(name string) string {
	return config.properties.GetPropertyByString(name, "")
}

func (config *PronunciationAssessmentConfig) String() string {
	jsonCch := C.pronunciation_assessment_config_to_json(config.handle)
	return C.GoString(jsonCch)
}

func (config *PronunciationAssessmentConfig) Close() {
	config.properties.Close()
	C.pronunciation_assessment_config_release(config.handle)
}
