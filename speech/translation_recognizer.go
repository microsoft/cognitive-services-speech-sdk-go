// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

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
//
// /* Proxy functions forward declarations */
// void cgo_recognizer_session_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_session_stopped(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_speech_start_detected(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_speech_end_detected(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_recognized(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_recognizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_canceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_translation_synthesis(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
//
import "C"

// TranslationRecognizer is the class for translation recognizers.
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

// NewTranslationRecognizerFromConfig creates a translation recognizer from a speech translation config and audio config.
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

// NewTranslationRecognizerFromAutoDetectSourceLangConfig creates a translation recognizer from a speech translation config, auto detection source language config and audio config.
func NewTranslationRecognizerFromAutoDetectSourceLangConfig(config *SpeechTranslationConfig, langConfig *AutoDetectSourceLanguageConfig, audioConfig *audio.AudioConfig) (*TranslationRecognizer, error) {
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
	return newTranslationRecognizerFromHandle(handle)
}

// RecognizeOnceAsync starts translation recognition, and returns after a single utterance is recognized.
// The end of a single utterance is determined by listening for silence at the end or until a maximum
// of 15 seconds of audio is processed. The task returns the recognition text as result.
// Note: Since RecognizeOnceAsync() returns only a single utterance, it is suitable only for single
// shot recognition like command or query.
// For long-running multi-utterance recognition, use StartContinuousRecognitionAsync() instead.
func (recognizer TranslationRecognizer) RecognizeOnceAsync() chan TranslationRecognitionOutcome {
	outcome := make(chan TranslationRecognitionOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		ret := uintptr(C.recognizer_recognize_once(recognizer.handle, &handle))
		if ret != C.SPX_NOERROR {
			outcome <- TranslationRecognitionOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewTranslationRecognitionResultFromHandle(handle2uintptr(handle))
			outcome <- TranslationRecognitionOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// StartContinuousRecognitionAsync asynchronously initiates continuous translation recognition operation.
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

// StopContinuousRecognitionAsync asynchronously terminates ongoing continuous translation recognition operation.
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

// GetEndpointID gets the endpoint ID of a customized speech model that is used for translation recognition.
func (recognizer TranslationRecognizer) GetEndpointID() string {
	return recognizer.Properties.GetProperty(common.SpeechServiceConnectionEndpointID, "")
}

// SetAuthorizationToken sets the authorization token that will be used for connecting to the service.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token
// expires, the caller needs to refresh it by calling this setter with a new valid token.
// Otherwise, the recognizer will encounter errors during recognition.
func (recognizer TranslationRecognizer) SetAuthorizationToken(token string) error {
	return recognizer.Properties.SetProperty(common.SpeechServiceAuthorizationToken, token)
}

// AuthorizationToken is the authorization token.
func (recognizer TranslationRecognizer) AuthorizationToken() string {
	return recognizer.Properties.GetProperty(common.SpeechServiceAuthorizationToken, "")
}

// SessionStarted signals events indicating the start of a recognition session (operation).
func (recognizer TranslationRecognizer) SessionStarted(handler SessionEventHandler) {
	registerSessionStartedCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_session_started_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_session_started)),
			nil)
	} else {
		C.recognizer_session_started_set_callback(recognizer.handle, nil, nil)
	}
}

// SessionStopped signals events indicating the end of a recognition session (operation).
func (recognizer TranslationRecognizer) SessionStopped(handler SessionEventHandler) {
	registerSessionStoppedCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_session_stopped_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_session_stopped)),
			nil)
	} else {
		C.recognizer_session_stopped_set_callback(recognizer.handle, nil, nil)
	}
}

// SpeechStartDetected signals for events indicating the start of speech.
func (recognizer TranslationRecognizer) SpeechStartDetected(handler RecognitionEventHandler) {
	registerSpeechStartDetectedCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_speech_start_detected_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_speech_start_detected)),
			nil)
	} else {
		C.recognizer_speech_start_detected_set_callback(recognizer.handle, nil, nil)
	}
}

// SpeechEndDetected signals for events indicating the end of speech.
func (recognizer TranslationRecognizer) SpeechEndDetected(handler RecognitionEventHandler) {
	registerSpeechEndDetectedCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_speech_end_detected_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_speech_end_detected)),
			nil)
	} else {
		C.recognizer_speech_end_detected_set_callback(recognizer.handle, nil, nil)
	}
}

// Recognizing signals for events containing intermediate recognition results.
func (recognizer TranslationRecognizer) Recognizing(handler TranslationRecognitionEventHandler) {
	registerTranslationRecognizingCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_recognizing_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_recognizing)),
			nil)
	} else {
		C.recognizer_recognizing_set_callback(recognizer.handle, nil, nil)
	}
}

// Recognized signals for events containing final recognition results.
// (indicating a successful recognition attempt).
func (recognizer TranslationRecognizer) Recognized(handler TranslationRecognitionEventHandler) {
	registerTranslationRecognizedCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_recognized_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_recognized)),
			nil)
	} else {
		C.recognizer_recognized_set_callback(recognizer.handle, nil, nil)
	}
}

// Canceled signals for events containing canceled recognition results
// (indicating a recognition attempt that was canceled as a result or a direct cancellation request
// or, alternatively, a transport or protocol failure).
func (recognizer TranslationRecognizer) Canceled(handler TranslationRecognitionCanceledEventHandler) {
	registerTranslationCanceledCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_canceled_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_canceled)),
			nil)
	} else {
		C.recognizer_canceled_set_callback(recognizer.handle, nil, nil)
	}
}

// Synthesizing signals for events containing translation synthesis results.
func (recognizer TranslationRecognizer) Synthesizing(handler TranslationSynthesisEventHandler) {
	registerTranslationSynthesisCallback(handler, recognizer.handle)
	if handler != nil {
		C.translator_synthesizing_audio_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_translation_synthesis)),
			nil)
	} else {
		C.translator_synthesizing_audio_set_callback(recognizer.handle, nil, nil)
	}
}

// Close disposes the associated resources.
func (recognizer TranslationRecognizer) Close() {
	recognizer.SessionStarted(nil)
	recognizer.SessionStopped(nil)
	recognizer.SpeechStartDetected(nil)
	recognizer.SpeechEndDetected(nil)
	recognizer.Recognizing(nil)
	recognizer.Recognized(nil)
	recognizer.Canceled(nil)
	recognizer.Synthesizing(nil)
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
