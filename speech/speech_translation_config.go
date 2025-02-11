// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"strings"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_speech_config.h>
// #include <speechapi_c_speech_translation_config.h>
//
import "C"

// SpeechTranslationConfig defines configurations for translation with speech input.
type SpeechTranslationConfig struct {
	SpeechConfig
}

// NewSpeechTranslationConfigFromSubscription creates a speech translation config instance with specified subscription key and region.
func NewSpeechTranslationConfigFromSubscription(subscription string, region string) (*SpeechTranslationConfig, error) {
	var handle C.SPXHANDLE
	subscriptionCStr := C.CString(subscription)
	defer C.free(unsafe.Pointer(subscriptionCStr))
	regionCStr := C.CString(region)
	defer C.free(unsafe.Pointer(regionCStr))

	ret := uintptr(C.speech_translation_config_from_subscription(&handle, subscriptionCStr, regionCStr))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	config := new(SpeechTranslationConfig)
	config.handle = handle
	return config, nil
}

// NewSpeechTranslationConfigFromAuthorizationToken creates a speech translation config instance with specified authorization token and region.
func NewSpeechTranslationConfigFromAuthorizationToken(authToken string, region string) (*SpeechTranslationConfig, error) {
	var handle C.SPXHANDLE
	authTokenCStr := C.CString(authToken)
	defer C.free(unsafe.Pointer(authTokenCStr))
	regionCStr := C.CString(region)
	defer C.free(unsafe.Pointer(regionCStr))

	ret := uintptr(C.speech_translation_config_from_authorization_token(&handle, authTokenCStr, regionCStr))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	config := new(SpeechTranslationConfig)
	config.handle = handle
	return config, nil
}

// NewSpeechTranslationConfigFromEndpoint creates a speech translation config instance with specified endpoint and subscription.
func NewSpeechTranslationConfigFromEndpoint(endpoint string, subscription string) (*SpeechTranslationConfig, error) {
	var handle C.SPXHANDLE
	endpointCStr := C.CString(endpoint)
	defer C.free(unsafe.Pointer(endpointCStr))
	subscriptionCStr := C.CString(subscription)
	defer C.free(unsafe.Pointer(subscriptionCStr))

	ret := uintptr(C.speech_translation_config_from_endpoint(&handle, endpointCStr, subscriptionCStr))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	config := new(SpeechTranslationConfig)
	config.handle = handle
	return config, nil
}

// NewSpeechTranslationConfigFromHost creates a speech translation config instance with specified host and subscription.
func NewSpeechTranslationConfigFromHost(host string, subscription string) (*SpeechTranslationConfig, error) {
	var handle C.SPXHANDLE
	hostCStr := C.CString(host)
	defer C.free(unsafe.Pointer(hostCStr))
	subscriptionCStr := C.CString(subscription)
	defer C.free(unsafe.Pointer(subscriptionCStr))

	ret := uintptr(C.speech_translation_config_from_host(&handle, hostCStr, subscriptionCStr))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	config := new(SpeechTranslationConfig)
	config.handle = handle
	return config, nil
}

// AddTargetLanguage adds a target language for translation.
func (config *SpeechTranslationConfig) AddTargetLanguage(language string) error {
	languageCStr := C.CString(language)
	defer C.free(unsafe.Pointer(languageCStr))

	ret := uintptr(C.speech_translation_config_add_target_language(config.handle, languageCStr))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// RemoveTargetLanguage removes a target language for translation.
func (config *SpeechTranslationConfig) RemoveTargetLanguage(language string) error {
	languageCStr := C.CString(language)
	defer C.free(unsafe.Pointer(languageCStr))

	ret := uintptr(C.speech_translation_config_remove_target_language(config.handle, languageCStr))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// GetTargetLanguages gets target languages for translation.
func (config *SpeechTranslationConfig) GetTargetLanguages() []string {
	languages := config.GetProperty(common.SpeechServiceConnectionTranslationToLanguages, "")
	if languages == "" {
		return []string{}
	}
	return strings.Split(languages, ",")
}

// SetCustomModelCategoryId sets a Category Id that will be passed to service.
// Category Id is used to find the custom model.
func (config *SpeechTranslationConfig) SetCustomModelCategoryId(categoryId string) error {
	categoryIdCStr := C.CString(categoryId)
	defer C.free(unsafe.Pointer(categoryIdCStr))

	ret := uintptr(C.speech_translation_config_set_custom_model_category_id(config.handle, categoryIdCStr))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// SetVoiceName sets output voice name.
func (config *SpeechTranslationConfig) SetVoiceName(voice string) {
	config.SetProperty(common.SpeechServiceConnectionTranslationVoice, voice)
}

// GetVoiceName gets output voice name.
func (config *SpeechTranslationConfig) GetVoiceName() string {
	return config.GetProperty(common.SpeechServiceConnectionTranslationVoice, "")
}
