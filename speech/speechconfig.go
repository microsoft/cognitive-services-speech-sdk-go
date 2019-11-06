package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)
// #cgo CFLAGS: -I/home/gelecaro/carbon/current/include/c_api
// #cgo LDFLAGS: -L/home/gelecaro/carbon/current/lib/x64 -lMicrosoft.CognitiveServices.Speech.core
// #include <stdlib.h>
// #include <speechapi_c_speech_config.h>
// #include <speechapi_c_property_bag.h>
import "C"
import "unsafe"

// SpeechConfig is the class that defines configurations for speech / intent recognition, or speech synthesis.
type SpeechConfig struct {
	handle C.SPXHANDLE
	propertyBagHandle C.SPXPROPERTYBAGHANDLE
}

// NewSpeechConfigFromSubscription creates an instance of the speech config with specified subscription key and region.
func NewSpeechConfigFromSubscription(subscriptionKey string, region string) (*SpeechConfig, error) {
	var handle C.SPXHANDLE
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	sk := C.CString(subscriptionKey)
	defer C.free(unsafe.Pointer(sk))
	r := C.CString(region)
	defer C.free(unsafe.Pointer(r))
	ret := uintptr(C.speech_config_from_subscription(&handle, sk, r))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.speech_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(SpeechConfig)
	config.handle = handle
	config.propertyBagHandle = propBagHandle
	return config, nil
}

// NewSpeechConfigFromAuthorizationToken creates an instance of the speech config with specified authorization token and
// region.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token expires, the
// caller needs to refresh it by calling this setter with a new valid token.
// As configuration values are copied when creating a new recognizer, the new token value will not apply to recognizers
// that have already been created.
// For recognizers that have been created before, you need to set authorization token of the corresponding recognizer
// to refresh the token. Otherwise, the recognizers will encounter errors during recognition.
func NewSpeechConfigFromAuthorizationToken(authorizationToken string, region string) (*SpeechConfig, error) {
	var handle C.SPXHANDLE
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	authToken := C.CString(authorizationToken)
	defer C.free(unsafe.Pointer(authToken))
	r := C.CString(region)
	defer C.free(unsafe.Pointer(r))
	ret := uintptr(C.speech_config_from_authorization_token(&handle, authToken, r))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.speech_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(SpeechConfig)
	config.handle = handle
	config.propertyBagHandle = propBagHandle
	return config, nil
}

// NewSpeechConfigFromEndpointWithSubscription creates an instance of the speech config with specified endpoint
// and subscription.
// This method is intended only for users who use a non-standard service endpoint.
// Note: The query parameters specified in the endpoint URI are not changed, even if they are set by any other APIs.
// For example, if the recognition language is defined in URI as query parameter "language=de-DE", and also set by
// SetSpeechRecognitionLanguage("en-US"), the language setting in URI takes precedence, and the effective language
// is "de-DE".
/// Only the parameters that are not specified in the endpoint URI can be set by other APIs.
/// Note: To use an authorization token with endoint, use FromEndpoint,
/// and then call SetAuthorizationToken() on the created SpeechConfig instance.
func NewSpeechConfigFromEndpointWithSubscription(endpoint string, subscriptionKey string) (*SpeechConfig, error) {
	var handle C.SPXHANDLE
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	e := C.CString(endpoint)
	defer C.free(unsafe.Pointer(e))
	sk := C.CString(subscriptionKey)
	defer C.free(unsafe.Pointer(sk))
	ret := uintptr(C.speech_config_from_endpoint(&handle, e, sk))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.speech_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(SpeechConfig)
	config.handle = handle
	config.propertyBagHandle = propBagHandle
	return config, nil
}

// NewSpeechConfigFromEndpoint creates an instance of SpeechConfig with specified endpoint.
// This method is intended only for users who use a non-standard service endpoint.
// Note: The query parameters specified in the endpoint URI are not changed, even if they are set by any other APIs.
// For example, if the recognition language is defined in URI as query parameter "language=de-DE", and also set by
// SetSpeechRecognitionLanguage("en-US"), the language setting in URI takes precedence, and the effective language is
// "de-DE".
// Only the parameters that are not specified in the endpoint URI can be set by other APIs.
// Note: If the endpoint requires a subscription key for authentication, use NewSpeechConfigFromEndpointWithSubscription
// to pass the subscription key as parameter.
// To use an authorization token with FromEndpoint, use this method to create a SpeechConfig instance, and then
// call SetAuthorizationToken() on the created SpeechConfig instance.
// Note: Added in version 1.5.0.
func NewSpeechConfigFromEndpoint(endpoint string) (*SpeechConfig, error) {
	var handle C.SPXHANDLE
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	e := C.CString(endpoint)
	defer C.free(unsafe.Pointer(e))
	ret := uintptr(C.speech_config_from_endpoint(&handle, e, nil))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.speech_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(SpeechConfig)
	config.handle = handle
	config.propertyBagHandle = propBagHandle
	return config, nil
}

// NewSpeechConfigFromHostWithSubscription creates an instance of the speech config with specified host and subscription.
// This method is intended only for users who use a non-default service host. Standard resource path will be assumed.
// For services with a non-standard resource path or no path at all, use FromEndpoint instead.
// Note: Query parameters are not allowed in the host URI and must be set by other APIs.
// Note: To use an authorization token with host, use NewSpeechConfigFromHost,
// and then call SetAuthorizationToken() on the created SpeechConfig instance.
// Note: Added in version 1.8.0.
func NewSpeechConfigFromHostWithSubscription(host string, subscriptionKey string) (*SpeechConfig, error) {
	var handle C.SPXHANDLE
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	h := C.CString(host)
	defer C.free(unsafe.Pointer(h))
	sk := C.CString(subscriptionKey)
	defer C.free(unsafe.Pointer(sk))
	ret := uintptr(C.speech_config_from_host(&handle, h, sk))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.speech_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(SpeechConfig)
	config.handle = handle
	config.propertyBagHandle = propBagHandle
	return config, nil
}

// NewSpeechConfigFromHost Creates an instance of SpeechConfig with specified host.
// This method is intended only for users who use a non-default service host. Standard resource path will be assumed.
// For services with a non-standard resource path or no path at all, use FromEndpoint instead.
// Note: Query parameters are not allowed in the host URI and must be set by other APIs.
// Note: If the host requires a subscription key for authentication, use NewSpeechConfigFromHostWithSubscription to pass
// the subscription key as parameter.
// To use an authorization token with FromHost, use this method to create a SpeechConfig instance, and then
// call SetAuthorizationToken() on the created SpeechConfig instance.
// Note: Added in version 1.8.0.
func NewSpeechConfigFromHost(host string) (*SpeechConfig, error) {
	var handle C.SPXHANDLE
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	h := C.CString(host)
	defer C.free(unsafe.Pointer(h))
	ret := uintptr(C.speech_config_from_host(&handle, h, nil))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.speech_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(SpeechConfig)
	config.handle = handle
	config.propertyBagHandle = propBagHandle
	return config, nil
}

// Only getters
func (config *SpeechConfig) SubscriptionKey() string {
	return config.GetProperty(common.SpeechServiceConnectionKey)
}

func (config *SpeechConfig) Region() string {
	return config.GetProperty(common.SpeechServiceConnectionRegion)
}

// Getter-setter
func (config *SpeechConfig) AuthorizationToken() string {
	return config.GetProperty(common.SpeechServiceAuthorizationToken)
}

func (config *SpeechConfig) SetAuthorizationToken(authToken string) error {
	return config.SetProperty(common.SpeechServiceAuthorizationToken, authToken)
}

func (config *SpeechConfig) SpeechRecognitionLanguage() string {
	return config.GetProperty(common.SpeechServiceConnectionRecoLanguage)
}

func (config *SpeechConfig) SetSpeechRecognitionLanguage(language string) error {
	return config.SetProperty(common.SpeechServiceConnectionRecoLanguage, language);
}

func (config *SpeechConfig) OutputFormat() common.OutputFormat {
	format := config.GetProperty(common.SpeechServiceResponseRequestDetailedResultTrueFalse)
	if format == "true" {
		return common.Detailed
	}
	return common.Simple
}

func (config *SpeechConfig) SetOutputFormat(outputFormat common.OutputFormat) error {
	val := "false"
	if outputFormat == common.Detailed {
		val = "true"
	}
	return config.SetProperty(common.SpeechServiceResponseRequestDetailedResultTrueFalse, val)
}

func (config *SpeechConfig) EndpointId() string {
	return config.GetProperty(common.SpeechServiceConnectionEndpointID)
}

func (config *SpeechConfig) SetEndpointId(endpointId string) error {
	return config.SetProperty(common.SpeechServiceConnectionEndpointID, endpointId)
}

// func (config *SpeechConfig) SpeechSynthesisLanguage
// {}
// func (config *SpeechConfig) SetSpeechSynthesisLanguage
// {}
// func (config *SpeechConfig) SpeechSynthesisVoiceName
// {}
// func (config *SpeechConfig) SetSpeechSynthesisVoiceName
// {}
// func (config *SpeechConfig) SpeechSynthesisOutputFormat
// {}
// func (config *SpeechConfig) SetSpeechSynthesisOutputFormat
// {}
// func (config *SpeechConfig) SpeechSynthesisOutputFormat
// {}
// func (config *SpeechConfig) SetSpeechSynthesisOutputFormat
// {}

// // Member functions
// func (config *SpeechConfig) SetProxy
// {}

func (config *SpeechConfig) SetProperty(id common.PropertyID, value string) error {
	v := C.CString(value)
	ret := uintptr(C.property_bag_set_string(config.propertyBagHandle, (C.int)(id), nil, v))
	C.free(unsafe.Pointer(v))
	if (ret != C.SPX_NOERROR) {
		return common.NewCarbonError(ret)
	}
	return nil
}

func (config *SpeechConfig) GetProperty(id common.PropertyID) string {
	emptyString := C.CString("")
	defer C.free(unsafe.Pointer(emptyString))
	value := C.property_bag_get_string(config.propertyBagHandle, (C.int)(id), nil, emptyString)
	goValue := C.GoString(value)
	C.property_bag_free_string(value)
	return goValue
}

// func (config *SpeechConfig) SetServiceProperty
// {}
// func (config *SpeechConfig) SetProfanity
// {}
// func (config *SpeechConfig) EnableAudioLogging
// {}
// func (config *SpeechConfig) RequestWordLevelTimestamps
// {}
// func (config *SpeechConfig) EnableDictation
// {}
