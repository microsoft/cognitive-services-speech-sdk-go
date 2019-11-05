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

// func NewSpeechConfigFromAuthorizationToken(authorizationToken string, region string) (SpeechSpeechConfig, error) {

// }
// func NewSpeechConfigFromEndpoint(endpoint string, subscriptionKey string) (SpeechConfig, error) {

// }
// func NewSpeechConfigFromHost(host string, subscriptionKey string) (SpeechConfig, error) {

// }
// /* end */

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
