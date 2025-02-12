// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"strings"
	"unsafe"

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
	speechConfig, err := NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}

	config.SpeechConfig = *speechConfig

	return config, nil
}

// NewSpeechTranslationConfigFromAuthorizationToken creates a speech translation config instance with specified authorization token and region.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token expires, the
// caller needs to refresh it by calling this setter with a new valid token.
// As configuration values are copied when creating a new recognizer, the new token value will not apply to recognizers
// that have already been created.
// For recognizers that have been created before, you need to set authorization token of the corresponding recognizer
// to refresh the token. Otherwise, the recognizers will encounter errors during recognition.
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
	speechConfig, err := NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}

	config.SpeechConfig = *speechConfig

	return config, nil
}

// NewSpeechTranslationConfigFromEndpointWithSubscription creates a speech translation config instance with specified endpoint and subscription.
// This method is intended only for users who use a non-standard service endpoint.
// Note: The query parameters specified in the endpoint URI are not changed, even if they are set by any other APIs.
// For example, if the recognition language is defined in URI as query parameter "language=de-DE", and also set by
// SetSpeechRecognitionLanguage("en-US"), the language setting in URI takes precedence, and the effective language
// is "de-DE".
// / Only the parameters that are not specified in the endpoint URI can be set by other APIs.
// / Note: To use an authorization token with endoint, use FromEndpoint,
// / and then call SetAuthorizationToken() on the created SpeechConfig instance.
func NewSpeechTranslationConfigFromEndpointWithSubscription(endpoint string, subscription string) (*SpeechTranslationConfig, error) {
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
	speechConfig, err := NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}

	config.SpeechConfig = *speechConfig

	return config, nil
}

// NewSpeechTranslationConfigFromEndpoint creates a speech translation config instance with specified endpoint and subscription.
// This method is intended only for users who use a non-standard service endpoint.
// Note: The query parameters specified in the endpoint URI are not changed, even if they are set by any other APIs.
// For example, if the recognition language is defined in URI as query parameter "language=de-DE", and also set by
// SetSpeechRecognitionLanguage("en-US"), the language setting in URI takes precedence, and the effective language
// is "de-DE".
// / Only the parameters that are not specified in the endpoint URI can be set by other APIs.
func NewSpeechTranslationConfigFromEndpoint(endpoint string) (*SpeechTranslationConfig, error) {
	var handle C.SPXHANDLE
	endpointCStr := C.CString(endpoint)
	defer C.free(unsafe.Pointer(endpointCStr))

	ret := uintptr(C.speech_translation_config_from_endpoint(&handle, endpointCStr, nil))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	config := new(SpeechTranslationConfig)
	speechConfig, err := NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}

	config.SpeechConfig = *speechConfig

	return config, nil
}

// NewSpeechTranslationConfigFromHostWithSubscription creates a speech translation config instance with specified host and subscription.
// This method is intended only for users who use a non-default service host. Standard resource path will be assumed.
// For services with a non-standard resource path or no path at all, use FromEndpoint instead.
// Note: Query parameters are not allowed in the host URI and must be set by other APIs.
// Note: To use an authorization token with host, use NewSpeechConfigFromHost,
// and then call SetAuthorizationToken() on the created SpeechConfig instance.
func NewSpeechTranslationConfigFromHostWithSubscription(host string, subscription string) (*SpeechTranslationConfig, error) {
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
	speechConfig, err := NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}

	config.SpeechConfig = *speechConfig

	return config, nil
}

// NewSpeechTranslationConfigFromHost creates a speech translation config instance with specified host and subscription.
// This method is intended only for users who use a non-default service host. Standard resource path will be assumed.
// For services with a non-standard resource path or no path at all, use FromEndpoint instead.
// Note: Query parameters are not allowed in the host URI and must be set by other APIs.
// Note: If the host requires a subscription key for authentication, use NewSpeechConfigFromHostWithSubscription to pass
// the subscription key as parameter.
// To use an authorization token with FromHost, use this method to create a SpeechConfig instance, and then
// call SetAuthorizationToken() on the created SpeechConfig instance.
func NewSpeechTranslationConfigFromHost(host string) (*SpeechTranslationConfig, error) {
	var handle C.SPXHANDLE
	hostCStr := C.CString(host)
	defer C.free(unsafe.Pointer(hostCStr))

	ret := uintptr(C.speech_translation_config_from_host(&handle, hostCStr, nil))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	config := new(SpeechTranslationConfig)
	speechConfig, err := NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}

	config.SpeechConfig = *speechConfig

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
	languages := config.properties.GetProperty(common.SpeechServiceConnectionTranslationToLanguages, "")
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
	return config.GetProperty(common.SpeechServiceConnectionTranslationVoice)
}
