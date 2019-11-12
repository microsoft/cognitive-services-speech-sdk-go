package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"strconv"
)

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

// SubscriptionKey is the subscription key that is used to create Speech Recognizer or Intent Recognizer or Translation
// Recognizer or Speech Synthesizer
func (config *SpeechConfig) SubscriptionKey() string {
	return config.GetProperty(common.SpeechServiceConnectionKey)
}

// Region is the region key that used to create Speech Recognizer or Intent Recognizer or Translation Recognizer or
// Speech Synthesizer.
func (config *SpeechConfig) Region() string {
	return config.GetProperty(common.SpeechServiceConnectionRegion)
}

// AuthorizationToken is the authorization token to connect to the service.
func (config *SpeechConfig) AuthorizationToken() string {
	return config.GetProperty(common.SpeechServiceAuthorizationToken)
}

// SetAuthorizationToken sets the authorization token to connect to the service.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token
// expires, the caller needs to refresh it by calling this setter with a new valid token.
// As configuration values are copied when creating a new recognizer, the new token value will not apply to
// recognizers that have already been created.
// For recognizers that have been created before, you need to set authorization token of the corresponding recognizer
// to refresh the token. Otherwise, the recognizers will encounter errors during recognition.
func (config *SpeechConfig) SetAuthorizationToken(authToken string) error {
	return config.SetProperty(common.SpeechServiceAuthorizationToken, authToken)
}

// SpeechRecognitionLanguage is the input language to the speech recognition.
// The language is specified in BCP-47 format.
func (config *SpeechConfig) SpeechRecognitionLanguage() string {
	return config.GetProperty(common.SpeechServiceConnectionRecoLanguage)
}

// SetSpeechRecognitionLanguage sets the input language to the speech recognizer.
func (config *SpeechConfig) SetSpeechRecognitionLanguage(language string) error {
	return config.SetProperty(common.SpeechServiceConnectionRecoLanguage, language);
}

// OutputFormat is result output format.
func (config *SpeechConfig) OutputFormat() common.OutputFormat {
	format := config.GetProperty(common.SpeechServiceResponseRequestDetailedResultTrueFalse)
	if format == "true" {
		return common.Detailed
	}
	return common.Simple
}

// SetOutputFormat sets output format.
func (config *SpeechConfig) SetOutputFormat(outputFormat common.OutputFormat) error {
	val := "false"
	if outputFormat == common.Detailed {
		val = "true"
	}
	return config.SetProperty(common.SpeechServiceResponseRequestDetailedResultTrueFalse, val)
}

// EndpointID is the endpoint ID
func (config *SpeechConfig) EndpointID() string {
	return config.GetProperty(common.SpeechServiceConnectionEndpointID)
}

// SetEndpointID sets the endpoint ID
func (config *SpeechConfig) SetEndpointID(endpointID string) error {
	return config.SetProperty(common.SpeechServiceConnectionEndpointID, endpointID)
}

// SpeechSynthesisLanguage is the language of the speech synthesizer.
// Added in version 1.4.0
func (config *SpeechConfig) SpeechSynthesisLanguage() string {
	return config.GetProperty(common.SpeechServiceConnectionSynthLanguage)
}

// SetSpeechSynthesisLanguage sets the language of the speech synthesizer.
// Added in version 1.4.0
func (config *SpeechConfig) SetSpeechSynthesisLanguage(language string) error {
	return config.SetProperty(common.SpeechServiceConnectionSynthLanguage, language)
}

// SpeechSynthesisVoiceName is the voice of the speech synthesizer.
// Added in version 1.4.0
func (config *SpeechConfig) SpeechSynthesisVoiceName() string {
	return config.GetProperty(common.SpeechServiceConnectionSynthVoice);
}

// SetSpeechSynthesisVoiceName sets the voice of the speech synthesizer.
// Added in version 1.4.0
func (config *SpeechConfig) SetSpeechSynthesisVoiceName(voiceName string) error {
	return config.SetProperty(common.SpeechServiceConnectionSynthVoice, voiceName)
}

// SpeechSynthesisOutputFormat is the speech synthesis output format.
// Added in version 1.4.0
func (config *SpeechConfig) SpeechSynthesisOutputFormat() string {
	return config.GetProperty(common.SpeechServiceConnectionSynthOutputFormat);
}

// SetSpeechSynthesisOutputFormat sets the speech synthesis output format (e.g. Riff16Khz16BitMonoPcm).
// Added in version 1.4.0
func (config *SpeechConfig) SetSpeechSynthesisOutputFormat(format common.SpeechSynthesisOutputFormat) error {
	ret := uintptr(C.speech_config_set_audio_output_format(config.handle, (C.Speech_Synthesis_Output_Format)(format)))
	if (ret != C.SPX_NOERROR) {
		return common.NewCarbonError(ret)
	}
	return nil
}

// SetProxy sets proxy configuration
// Added in version 1.1.0
//
// Note: Proxy functionality is not available on macOS. This function will have no effect on this platform.
func (config *SpeechConfig) SetProxy(hostname string, port uint64) error {
	res := config.SetProperty(common.SpeechServiceConnectionProxyHostName, hostname)
	if res != nil {
		return res
	}
	return config.SetProperty(common.SpeechServiceConnectionProxyPort, strconv.FormatUint(port, 10))
}

// SetProxyWithUsernameAndPassword sets proxy configuration with username and password
// Added in version 1.1.0
//
// Note: Proxy functionality is not available on macOS. This function will have no effect on this platform.
func (config *SpeechConfig) SetProxyWithUsernameAndPassword(hostname string, port uint64, username string, password string) error {
	res := config.SetProxy(hostname, port)
	if res != nil {
		return res
	}
	res = config.SetProperty(common.SpeechServiceConnectionProxyUserName, username)
	if res != nil {
		return res
	}
	return config.SetProperty(common.SpeechServiceConnectionProxyPassword, password)
}

// SetProperty sets a property value by ID.
func (config *SpeechConfig) SetProperty(id common.PropertyID, value string) error {
	v := C.CString(value)
	ret := uintptr(C.property_bag_set_string(config.propertyBagHandle, (C.int)(id), nil, v))
	C.free(unsafe.Pointer(v))
	if (ret != C.SPX_NOERROR) {
		return common.NewCarbonError(ret)
	}
	return nil
}

// GetProperty gets a property value by ID.
func (config *SpeechConfig) GetProperty(id common.PropertyID) string {
	emptyString := C.CString("")
	defer C.free(unsafe.Pointer(emptyString))
	value := C.property_bag_get_string(config.propertyBagHandle, (C.int)(id), nil, emptyString)
	goValue := C.GoString(value)
	C.property_bag_free_string(value)
	return goValue
}

// SetPropertyByString sets a property value by string.
func (config *SpeechConfig) SetPropertyByString(name string, value string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))
	ret := uintptr(C.property_bag_set_string(config.propertyBagHandle, -1, n, v))
	if (ret != C.SPX_NOERROR) {
		return common.NewCarbonError(ret)
	}
	return nil
}

// GetPropertyByString gets a property value by string.
func (config *SpeechConfig) GetPropertyByString(name string) string {
	emptyString := C.CString("")
	defer C.free(unsafe.Pointer(emptyString))
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	value := C.property_bag_get_string(config.propertyBagHandle, -1, n, emptyString)
	goValue := C.GoString(value)
	C.property_bag_free_string(value)
	return goValue
}

// SetServiceProperty sets a property value that will be passed to service using the specified channel.
// Added in version 1.5.0.
func (config *SpeechConfig) SetServiceProperty(name string, value string, channel common.ServicePropertyChannel) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))
	ret := uintptr(C.speech_config_set_service_property(config.handle, n, v, (C.SpeechConfig_ServicePropertyChannel)(channel)))
	if (ret != C.SPX_NOERROR) {
		return common.NewCarbonError(ret)
	}
	return nil
}

// SetProfanity sets profanity option.
// Added in version 1.5.0.
func (config *SpeechConfig) SetProfanity(profanity common.ProfanityOption) error {
	ret := uintptr(C.speech_config_set_profanity(config.handle, (C.SpeechConfig_ProfanityOption)(profanity)))
	if (ret != C.SPX_NOERROR) {
		return common.NewCarbonError(ret)
	}
	return nil
}

// EnableAudioLogging enables audio logging in service.
// Added in version 1.5.0.
func (config *SpeechConfig) EnableAudioLogging() error {
	return config.SetProperty(common.SpeechServiceConnectionEnableAudioLogging, "true")
}

// RequestWordLevelTimestamps includes word-level timestamps in response result.
// Added in version 1.5.0.
func (config *SpeechConfig) RequestWordLevelTimestamps() error {
	return config.SetProperty(common.SpeechServiceResponseRequestWordLevelTimestamps, "true")
}

// EnableDictation enables dictation mode. Only supported in speech continuous recognition.
// Added in version 1.5.0.
func (config *SpeechConfig) EnableDictation() error {
	return config.SetProperty(common.SpeechServiceConnectionRecoMode, "DICTATION")
}
