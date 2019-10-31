package speech

import (
	"fmt"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)
// #cgo CFLAGS: -I/home/gelecaro/carbon/current/include/c_api
// #cgo LDFLAGS: -L/home/gelecaro/carbon/current/lib/x64 -lMicrosoft.CognitiveServices.Speech.core
// #include <stdlib.h>
// #include <speechapi_c_speech_config.h>
import "C"
import "unsafe"

type SpeechConfig struct {
	handle C.SPXHANDLE
}

/* Factory methods */
func NewSpeechConfigFromSubscription(subscriptionKey string, region string) (*SpeechConfig, error) {
	var handle C.SPXHANDLE
	sk := C.CString(subscriptionKey)
	r := C.CString(region)
	ret := uintptr(C.speech_config_from_subscription(&handle, sk, r))
	C.free(unsafe.Pointer(sk))
	C.free(unsafe.Pointer(r))
	if ret == C.SPX_NOERROR {
		config := new(SpeechConfig)
		config.handle = handle
		return config, nil
	}
	return nil, common.NewCarbonError(ret)
}

// func NewSpeechConfigFromAuthorizationToken(authorizationToken string, region string) (SpeechSpeechConfig, error) {

// }
// func NewSpeechConfigFromEndpoint(endpoint string, subscriptionKey string) (SpeechConfig, error) {

// }
// func NewSpeechConfigFromHost(host string, subscriptionKey string) (SpeechConfig, error) {

// }
// /* end */

// // Only getters
func (config *SpeechConfig) SubscriptionKey() string {
	return fmt.Sprintf("%d", common.SpeechServiceAuthorization_Token)
}

// func (config *SpeechConfig) Region() string
// {

// }

// // Getter-setter
// func (config *SpeechConfig) AuthorizationToken() string
// {}
// func (config *SpeechConfig) SetAuthorizationToken()
// {}
// func (config *SpeechConfig) SpeechRecognitionLanguage
// {}
// func (config *SpeechConfig) SetSpeechRecognitionLanguage
// {}
// func (config *SpeechConfig) OutputFormat
// {}
// func (config *SpeechConfig) SetOutputFormat
// {}
// func (config *SpeechConfig) EndpointId
// {}
// func (config *SpeechConfig) SetEndpointId
// {}
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
// func (config *SpeechConfig) SetProperty
// {}
// func (config *SpeechConfig) GetProperty
// {}
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
