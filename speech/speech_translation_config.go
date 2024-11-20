package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_speech_config.h>
// #include <speechapi_c_speech_translation_config.h>
import "C"

type SpeechTranslationConfig struct {
	handle     C.SPXHANDLE
	properties *common.PropertyCollection
}

func (config SpeechTranslationConfig) GetHandle() common.SPXHandle {
	return handle2uintptr(config.handle)
}

func NewSpeechTranslationConfigFromHandle(handle common.SPXHandle) (*SpeechTranslationConfig, error) {
	var cHandle = uintptr2handle(handle)
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	ret := uintptr(C.speech_config_get_property_bag(cHandle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(SpeechTranslationConfig)
	config.handle = cHandle
	config.properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	err := config.properties.SetPropertyByString("SPEECHSDK-SPEECH-CONFIG-SYSTEM-LANGUAGE", "Go")
	if err != nil {
		config.Close()
		return nil, err
	}
	return config, nil
}

func NewSpeechTranslationConfigFromSubscription(subscriptionKey string, region string) (*SpeechTranslationConfig, error) {
	var handle C.SPXHANDLE
	sk := C.CString(subscriptionKey)
	defer C.free(unsafe.Pointer(sk))
	r := C.CString(region)
	defer C.free(unsafe.Pointer(r))
	ret := uintptr(C.speech_translation_config_from_subscription(&handle, sk, r))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return NewSpeechTranslationConfigFromHandle(handle2uintptr(handle))
}

func NewSpeechTranslationConfigFromEndpointWithSubscription(endpoint string, subscriptionKey string) (*SpeechTranslationConfig, error) {
	var handle C.SPXHANDLE
	e := C.CString(endpoint)
	defer C.free(unsafe.Pointer(e))
	sk := C.CString(subscriptionKey)
	defer C.free(unsafe.Pointer(sk))
	ret := uintptr(C.speech_translation_config_from_endpoint(&handle, e, sk))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return NewSpeechTranslationConfigFromHandle(handle2uintptr(handle))
}

func (config *SpeechTranslationConfig) SubscriptionKey() string {
	return config.GetProperty(common.SpeechServiceConnectionKey)
}

func (config *SpeechTranslationConfig) Region() string {
	return config.GetProperty(common.SpeechServiceConnectionRegion)
}

func (config *SpeechTranslationConfig) AddTargetLanguage(language string) error {
	l := C.CString(language)
	defer C.free(unsafe.Pointer(l))
	ret := uintptr(C.speech_translation_config_add_target_language(config.handle, l))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

func (config *SpeechTranslationConfig) SetSpeechRecognitionLanguage(language string) error {
	return config.SetProperty(common.SpeechServiceConnectionRecoLanguage, language)
}

func (config *SpeechTranslationConfig) SpeechRecognitionLanguage() string {
	return config.GetProperty(common.SpeechServiceConnectionRecoLanguage)
}

func (config *SpeechTranslationConfig) SetTranslationVoiceName(voiceName string) error {
	return config.SetProperty(common.SpeechServiceConnectionTranslationVoice, voiceName)
}

func (config *SpeechTranslationConfig) TranslationVoiceName() string {
	return config.GetProperty(common.SpeechServiceConnectionTranslationVoice)
}

func (config *SpeechTranslationConfig) SetProperty(id common.PropertyID, value string) error {
	return config.properties.SetProperty(id, value)
}

func (config *SpeechTranslationConfig) GetProperty(id common.PropertyID) string {
	return config.properties.GetProperty(id, "")
}

func (config *SpeechTranslationConfig) Close() {
	config.properties.Close()
	C.speech_config_release(config.handle)
}

func (config *SpeechTranslationConfig) getHandle() C.SPXHANDLE {
	return config.handle
}
