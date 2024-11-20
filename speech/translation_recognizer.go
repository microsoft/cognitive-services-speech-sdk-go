package speech

import (
	"math"
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_recognizer.h>
// #include <speechapi_c_translation_recognizer.h>
// #include <speechapi_c_factory.h>
// #include <speechapi_c_grammar.h>
//
// /* Proxy functions forward declarations */
// void cgo_recognizer_recognized(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_recognizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_canceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_translator_synthesizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
//
import "C"

type TranslationRecognizer struct {
	Properties                 *common.PropertyCollection
	handle                     C.SPXHANDLE
	handleAsyncStartContinuous C.SPXASYNCHANDLE
	handleAsyncStopContinuous  C.SPXASYNCHANDLE
}

func newTranslationRecognizerFromHandle(handle C.SPXHANDLE) (*TranslationRecognizer, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.recognizer_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	recognizer := new(TranslationRecognizer)
	recognizer.handle = handle
	recognizer.handleAsyncStartContinuous = C.SPXHANDLE_INVALID
	recognizer.handleAsyncStopContinuous = C.SPXHANDLE_INVALID
	recognizer.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return recognizer, nil
}

func NewTranslationRecognizerFromConfig(config *SpeechTranslationConfig, audioConfig *audio.AudioConfig) (*TranslationRecognizer, error) {
	var handle C.SPXHANDLE
	if config == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	configHandle := config.getHandle()
	var audioHandle C.SPXHANDLE
	if audioConfig == nil {
		audioHandle = nil
	} else {
		audioHandle = uintptr2handle(audioConfig.GetHandle())
	}
	ret := uintptr(C.recognizer_create_translation_recognizer_from_config(&handle, configHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newTranslationRecognizerFromHandle(handle)
}

func NewTranslationRecognizerFomAutoDetectSourceLangConfig(config *SpeechConfig, langConfig *AutoDetectSourceLanguageConfig, audioConfig *audio.AudioConfig) (*SpeechRecognizer, error) {
	var handle C.SPXHANDLE
	if config == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	if langConfig == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	configHandle := config.getHandle()
	langConfigHandle := langConfig.getHandle()
	var audioHandle C.SPXHANDLE
	if audioConfig == nil {
		audioHandle = nil
	} else {
		audioHandle = uintptr2handle(audioConfig.GetHandle())
	}
	ret := uintptr(C.recognizer_create_translation_recognizer_from_auto_detect_source_lang_config(&handle, configHandle, langConfigHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSpeechRecognizerFromHandle(handle)
}

func (recognizer TranslationRecognizer) StartContinuousRecognitionAsync() chan error {
	outcome := make(chan error)
	go func() {
		// Close any unfinished previous attempt
		ret := releaseAsyncHandleIfValid(&recognizer.handleAsyncStartContinuous)
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_start_continuous_recognition_async(recognizer.handle, &recognizer.handleAsyncStartContinuous))
		}
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_start_continuous_recognition_async_wait_for(recognizer.handleAsyncStartContinuous, math.MaxUint32))
		}
		releaseAsyncHandleIfValid(&recognizer.handleAsyncStartContinuous)
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
			return
		}
		outcome <- nil
	}()
	return outcome
}

func (recognizer TranslationRecognizer) StopContinuousRecognitionAsync() chan error {
	outcome := make(chan error)
	go func() {
		ret := releaseAsyncHandleIfValid(&recognizer.handleAsyncStopContinuous)
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_stop_continuous_recognition_async(recognizer.handle, &recognizer.handleAsyncStopContinuous))
		}
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_stop_continuous_recognition_async_wait_for(recognizer.handleAsyncStopContinuous, math.MaxUint32))
		}
		releaseAsyncHandleIfValid(&recognizer.handleAsyncStopContinuous)
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
			return
		}
		outcome <- nil
	}()
	return outcome
}

func (recognizer TranslationRecognizer) Recognizing(handler SpeechRecognitionEventHandler) {
	registerRecognizingCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_recognizing_set_callback(
			recognizer.handle,
			(C.PRECOGNITION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_recognizing)),
			nil)
	} else {
		C.recognizer_recognizing_set_callback(recognizer.handle, nil, nil)
	}
}

func (recognizer TranslationRecognizer) Recognized(handler SpeechRecognitionEventHandler) {
	registerRecognizedCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_recognized_set_callback(
			recognizer.handle,
			(C.PRECOGNITION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_recognized)),
			nil)
	} else {
		C.recognizer_recognized_set_callback(recognizer.handle, nil, nil)
	}
}

func (recognizer TranslationRecognizer) Canceled(handler SpeechRecognitionCanceledEventHandler) {
	registerCanceledCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_canceled_set_callback(
			recognizer.handle,
			(C.PRECOGNITION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_canceled)),
			nil)
	} else {
		C.recognizer_canceled_set_callback(recognizer.handle, nil, nil)
	}
}

func (recognizer TranslationRecognizer) Synthesizing(handler TranslationSynthesisEventHandler) {
	registerTranslationSynthesizingCallback(handler, recognizer.handle)
	if handler != nil {
		C.translator_synthesizing_audio_set_callback(
			recognizer.handle,
			(C.PTRANSLATIONSYNTHESIS_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_translator_synthesizing)),
			nil)
	} else {
		C.translator_synthesizing_audio_set_callback(recognizer.handle, nil, nil)
	}
}

func (recognizer TranslationRecognizer) Close() {
	recognizer.Recognizing(nil)
	recognizer.Recognized(nil)
	recognizer.Canceled(nil)
	var asyncHandles = []*C.SPXASYNCHANDLE{
		&recognizer.handleAsyncStartContinuous,
		&recognizer.handleAsyncStopContinuous,
	}
	for i := 0; i < len(asyncHandles); i++ {
		handle := asyncHandles[i]
		releaseAsyncHandleIfValid(handle)
	}
	recognizer.Properties.Close()
	if recognizer.handle != C.SPXHANDLE_INVALID {
		C.recognizer_handle_release(recognizer.handle)
		recognizer.handle = C.SPXHANDLE_INVALID
	}
}
