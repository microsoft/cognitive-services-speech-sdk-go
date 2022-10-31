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
// #include <speechapi_c_factory.h>
// #include <speechapi_c_grammar.h>
//
// /* Proxy functions forward declarations */
// void cgo_recognizer_session_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_session_stopped(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_speech_start_detected(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_speech_end_detected(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_recognized(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_recognizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_recognizer_canceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
//
import "C"

// SpeechRecognizer is the class for speech recognizers.
type SpeechRecognizer struct {
	Properties                 *common.PropertyCollection
	handle                     C.SPXHANDLE
	handleAsyncStartContinuous C.SPXASYNCHANDLE
	handleAsyncStopContinuous  C.SPXASYNCHANDLE
	handleAsyncStartKeyword    C.SPXASYNCHANDLE
	handleAsyncStopKeyword     C.SPXASYNCHANDLE
}

func newSpeechRecognizerFromHandle(handle C.SPXHANDLE) (*SpeechRecognizer, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.recognizer_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	recognizer := new(SpeechRecognizer)
	recognizer.handle = handle
	recognizer.handleAsyncStartContinuous = C.SPXHANDLE_INVALID
	recognizer.handleAsyncStopContinuous = C.SPXHANDLE_INVALID
	recognizer.handleAsyncStartKeyword = C.SPXHANDLE_INVALID
	recognizer.handleAsyncStopKeyword = C.SPXHANDLE_INVALID
	recognizer.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return recognizer, nil
}

// NewSpeechRecognizerFromConfig creates a speech recognizer from a speech config and audio config.
func NewSpeechRecognizerFromConfig(config *SpeechConfig, audioConfig *audio.AudioConfig) (*SpeechRecognizer, error) {
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
	ret := uintptr(C.recognizer_create_speech_recognizer_from_config(&handle, configHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSpeechRecognizerFromHandle(handle)
}

// NewSpeechRecognizerFomAutoDetectSourceLangConfig creates a speech recognizer from a speech config, auto detection source language config and audio config
func NewSpeechRecognizerFomAutoDetectSourceLangConfig(config *SpeechConfig, langConfig *AutoDetectSourceLanguageConfig, audioConfig *audio.AudioConfig) (*SpeechRecognizer, error) {
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
	ret := uintptr(C.recognizer_create_speech_recognizer_from_auto_detect_source_lang_config(&handle, configHandle, langConfigHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSpeechRecognizerFromHandle(handle)
}

// NewSpeechRecognizerFromSourceLanguageConfig creates a speech recognizer from a speech config, source language config and audio config
func NewSpeechRecognizerFromSourceLanguageConfig(config *SpeechConfig, sourceLanguageConfig *SourceLanguageConfig, audioConfig *audio.AudioConfig) (*SpeechRecognizer, error) {
	var handle C.SPXHANDLE
	if config == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	if sourceLanguageConfig == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	configHandle := config.getHandle()
	langConfigHandle := sourceLanguageConfig.getHandle()
	var audioHandle C.SPXHANDLE
	if audioConfig == nil {
		audioHandle = nil
	} else {
		audioHandle = uintptr2handle(audioConfig.GetHandle())
	}
	ret := uintptr(C.recognizer_create_speech_recognizer_from_source_lang_config(&handle, configHandle, langConfigHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSpeechRecognizerFromHandle(handle)
}

// NewSpeechRecognizerFromSourceLanguage creates a speech recognizer from a speech config, source language and audio config
func NewSpeechRecognizerFromSourceLanguage(config *SpeechConfig, sourceLanguage string, audioConfig *audio.AudioConfig) (*SpeechRecognizer, error) {
	languageConfig, err := NewSourceLanguageConfigFromLanguage(sourceLanguage)
	if err != nil {
		return nil, err
	}
	return NewSpeechRecognizerFromSourceLanguageConfig(config, languageConfig, audioConfig)
}

// RecognizeOnceAsync starts speech recognition, and returns after a single utterance is recognized.
// The end of a single utterance is determined by listening for silence at the end or until a maximum
// of 15 seconds of audio is processed.  The task returns the recognition text as result.
// Note: Since RecognizeOnceAsync() returns only a single utterance, it is suitable only for single
// shot recognition like command or query.
// For long-running multi-utterance recognition, use StartContinuousRecognitionAsync() instead.
func (recognizer SpeechRecognizer) RecognizeOnceAsync() chan SpeechRecognitionOutcome {
	outcome := make(chan SpeechRecognitionOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		ret := uintptr(C.recognizer_recognize_once(recognizer.handle, &handle))
		if ret != C.SPX_NOERROR {
			outcome <- SpeechRecognitionOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewSpeechRecognitionResultFromHandle(handle2uintptr(handle))
			outcome <- SpeechRecognitionOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

func releaseAsyncHandleIfValid(handle *C.SPXASYNCHANDLE) uintptr {
	ret := uintptr(C.SPX_NOERROR)
	if *handle != C.SPXHANDLE_INVALID && C.recognizer_async_handle_is_valid(*handle) {
		ret = uintptr(C.recognizer_async_handle_release(*handle))
		*handle = C.SPXHANDLE_INVALID
	}
	return ret
}

// StartContinuousRecognitionAsync asynchronously initiates continuous speech recognition operation.
func (recognizer SpeechRecognizer) StartContinuousRecognitionAsync() chan error {
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

// StopContinuousRecognitionAsync asynchronously terminates ongoing continuous speech recognition operation.
func (recognizer SpeechRecognizer) StopContinuousRecognitionAsync() chan error {
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

// StartKeywordRecognitionAsync asynchronously initiates keyword recognition operation.
func (recognizer SpeechRecognizer) StartKeywordRecognitionAsync(model KeywordRecognitionModel) chan error {
	outcome := make(chan error)
	modelHandle := uintptr2handle(model.GetHandle())
	go func() {
		ret := releaseAsyncHandleIfValid(&recognizer.handleAsyncStartKeyword)
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_start_keyword_recognition_async(recognizer.handle, modelHandle, &recognizer.handleAsyncStartKeyword))
		}
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_start_keyword_recognition_async_wait_for(recognizer.handleAsyncStartKeyword, math.MaxUint32))
		}
		releaseAsyncHandleIfValid(&recognizer.handleAsyncStartKeyword)
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
			return
		}
		outcome <- nil
	}()
	return outcome
}

// StopKeywordRecognitionAsync asynchronously terminates keyword recognition operation.
func (recognizer SpeechRecognizer) StopKeywordRecognitionAsync() chan error {
	outcome := make(chan error)
	go func() {
		ret := releaseAsyncHandleIfValid(&recognizer.handleAsyncStopKeyword)
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_stop_keyword_recognition_async(recognizer.handle, &recognizer.handleAsyncStopKeyword))
		}
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_stop_keyword_recognition_async_wait_for(recognizer.handleAsyncStopKeyword, math.MaxUint32))
		}
		releaseAsyncHandleIfValid(&recognizer.handleAsyncStopKeyword)
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
			return
		}
		outcome <- nil
	}()
	return outcome
}

// GetEndpointID gets the endpoint ID of a customized speech model that is used for speech recognition.
func (recognizer SpeechRecognizer) GetEndpointID() string {
	return recognizer.Properties.GetProperty(common.SpeechServiceConnectionEndpointID, "")
}

// SetAuthorizationToken sets the authorization token that will be used for connecting to the service.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token
// expires, the caller needs to refresh it by calling this setter with a new valid token.
// Otherwise, the recognizer will encounter errors during recognition.
func (recognizer SpeechRecognizer) SetAuthorizationToken(token string) error {
	return recognizer.Properties.SetProperty(common.SpeechServiceAuthorizationToken, token)
}

// AuthorizationToken is the authorization token.
func (recognizer SpeechRecognizer) AuthorizationToken() string {
	return recognizer.Properties.GetProperty(common.SpeechServiceAuthorizationToken, "")
}

// SessionStarted signals events indicating the start of a recognition session (operation).
func (recognizer SpeechRecognizer) SessionStarted(handler SessionEventHandler) {
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
func (recognizer SpeechRecognizer) SessionStopped(handler SessionEventHandler) {
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
func (recognizer SpeechRecognizer) SpeechStartDetected(handler RecognitionEventHandler) {
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

// SpeechEndDetected signals for  events indicating the end of speech.
func (recognizer SpeechRecognizer) SpeechEndDetected(handler RecognitionEventHandler) {
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
func (recognizer SpeechRecognizer) Recognizing(handler SpeechRecognitionEventHandler) {
	registerRecognizingCallback(handler, recognizer.handle)
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
func (recognizer SpeechRecognizer) Recognized(handler SpeechRecognitionEventHandler) {
	registerRecognizedCallback(handler, recognizer.handle)
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
func (recognizer SpeechRecognizer) Canceled(handler SpeechRecognitionCanceledEventHandler) {
	registerCanceledCallback(handler, recognizer.handle)
	if handler != nil {
		C.recognizer_canceled_set_callback(
			recognizer.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_recognizer_canceled)),
			nil)
	} else {
		C.recognizer_canceled_set_callback(recognizer.handle, nil, nil)
	}
}

// Close disposes the associated resources.
func (recognizer SpeechRecognizer) Close() {
	recognizer.SessionStarted(nil)
	recognizer.SessionStopped(nil)
	recognizer.SpeechStartDetected(nil)
	recognizer.SpeechEndDetected(nil)
	recognizer.Recognizing(nil)
	recognizer.Recognized(nil)
	recognizer.Canceled(nil)
	var asyncHandles = []*C.SPXASYNCHANDLE{
		&recognizer.handleAsyncStartContinuous,
		&recognizer.handleAsyncStopContinuous,
		&recognizer.handleAsyncStartKeyword,
		&recognizer.handleAsyncStopKeyword,
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

type grammarPhrase struct {
	handle C.SPXHANDLE
}

func grammarPhraseFromText(text string) (*grammarPhrase, error) {
	var handle C.SPXHANDLE
	txt := C.CString(text)
	defer C.free(unsafe.Pointer(txt))
	ret := uintptr(C.grammar_phrase_create_from_text(&handle, txt))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	phrase := new(grammarPhrase)
	phrase.handle = handle
	return phrase, nil
}

func (grammar *grammarPhrase) Close() {
	C.grammar_phrase_handle_release(grammar.handle)
}

type PhraseListGrammar struct {
	handle C.SPXHANDLE
}

// NewPhraseListGrammarFromRecognizer Creates a phrase list grammar for the specified recognizer.
func NewPhraseListGrammarFromRecognizer(recognizer *SpeechRecognizer) (*PhraseListGrammar, error) {
	var handle C.SPXHANDLE
	name := C.CString("")
	defer C.free(unsafe.Pointer(name))
	ret := uintptr(C.phrase_list_grammar_from_recognizer_by_name(&handle, recognizer.handle, name))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	grammar := new(PhraseListGrammar)
	grammar.handle = handle
	return grammar, nil
}

// Close releases the associated resources.
func (grammar *PhraseListGrammar) Close() {
	C.grammar_handle_release(grammar.handle)
}

// AddPhrase adds a simple phrase that may be spoken by the user.
func (grammar *PhraseListGrammar) AddPhrase(text string) error {
	phrase, err := grammarPhraseFromText(text)
	if err != nil {
		return err
	}
	defer phrase.Close()

	ret := uintptr(C.phrase_list_grammar_add_phrase(grammar.handle, phrase.handle))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// Clears all phrases from the phrase list grammar.
func (grammar *PhraseListGrammar) Clear() error {
	ret := uintptr(C.phrase_list_grammar_clear(grammar.handle))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}
