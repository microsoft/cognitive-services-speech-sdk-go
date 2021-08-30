// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_synthesizer.h>
// #include <speechapi_c_factory.h>
//
// /* Proxy functions forward declarations */
// void cgo_synthesizer_synthesis_started(SPXSYNTHHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_synthesizer_synthesizing(SPXSYNTHHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_synthesizer_synthesis_completed(SPXSYNTHHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_synthesizer_synthesis_canceled(SPXSYNTHHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_synthesizer_word_boundary(SPXSYNTHHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_synthesizer_viseme_received(SPXSYNTHHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_synthesizer_bookmark_reached(SPXSYNTHHANDLE handle, SPXEVENTHANDLE event, void* context);
//
import "C"

// SpeechSynthesizer is the class for speech synthesizer.
type SpeechSynthesizer struct {
	Properties *common.PropertyCollection
	handle     C.SPXHANDLE
}

func newSpeechSynthesizerFromHandle(handle C.SPXHANDLE) (*SpeechSynthesizer, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.synthesizer_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	synthesizer := new(SpeechSynthesizer)
	synthesizer.handle = handle
	synthesizer.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return synthesizer, nil
}

// NewSpeechSynthesizerFromConfig creates a speech synthesizer from a speech config and audio config.
func NewSpeechSynthesizerFromConfig(config *SpeechConfig, audioConfig *audio.AudioConfig) (*SpeechSynthesizer, error) {
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
	ret := uintptr(C.synthesizer_create_speech_synthesizer_from_config(&handle, configHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSpeechSynthesizerFromHandle(handle)
}

// NewSpeechSynthesizerFomAutoDetectSourceLangConfig creates a speech synthesizer from a speech config, auto detection source language config and audio config
func NewSpeechSynthesizerFomAutoDetectSourceLangConfig(config *SpeechConfig, langConfig *AutoDetectSourceLanguageConfig, audioConfig *audio.AudioConfig) (*SpeechSynthesizer, error) {
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
	ret := uintptr(C.synthesizer_create_speech_synthesizer_from_auto_detect_source_lang_config(&handle, configHandle, langConfigHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSpeechSynthesizerFromHandle(handle)
}

// SpeakTextAsync executes the speech synthesis on plain text, asynchronously.
func (synthesizer SpeechSynthesizer) SpeakTextAsync(text string) chan SpeechSynthesisOutcome {
	outcome := make(chan SpeechSynthesisOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		cText := C.CString(text)
		defer C.free(unsafe.Pointer(cText))
		length := len(text)
		ret := uintptr(C.synthesizer_speak_text(synthesizer.handle, cText, (C.uint32_t)(length), &handle))
		if ret != C.SPX_NOERROR {
			outcome <- SpeechSynthesisOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewSpeechSynthesisResultFromHandle(handle2uintptr(handle))
			outcome <- SpeechSynthesisOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// SpeakSsmlAsync executes the speech synthesis on SSML, asynchronously.
func (synthesizer SpeechSynthesizer) SpeakSsmlAsync(ssml string) chan SpeechSynthesisOutcome {
	outcome := make(chan SpeechSynthesisOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		cText := C.CString(ssml)
		defer C.free(unsafe.Pointer(cText))
		length := len(ssml)
		ret := uintptr(C.synthesizer_speak_ssml(synthesizer.handle, cText, (C.uint32_t)(length), &handle))
		if ret != C.SPX_NOERROR {
			outcome <- SpeechSynthesisOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewSpeechSynthesisResultFromHandle(handle2uintptr(handle))
			outcome <- SpeechSynthesisOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// StartSpeakingTextAsync starts the speech synthesis on plain text, asynchronously.
// It returns when the synthesis request is started to process (the result reason is SynthesizingAudioStarted).
func (synthesizer SpeechSynthesizer) StartSpeakingTextAsync(text string) chan SpeechSynthesisOutcome {
	outcome := make(chan SpeechSynthesisOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		cText := C.CString(text)
		defer C.free(unsafe.Pointer(cText))
		length := len(text)
		ret := uintptr(C.synthesizer_start_speaking_text(synthesizer.handle, cText, (C.uint32_t)(length), &handle))
		if ret != C.SPX_NOERROR {
			outcome <- SpeechSynthesisOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewSpeechSynthesisResultFromHandle(handle2uintptr(handle))
			outcome <- SpeechSynthesisOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// StartSpeakingSsmlAsync starts the speech synthesis on SSML, asynchronously.
// It returns when the synthesis request is started to process (the result reason is SynthesizingAudioStarted).
func (synthesizer SpeechSynthesizer) StartSpeakingSsmlAsync(ssml string) chan SpeechSynthesisOutcome {
	outcome := make(chan SpeechSynthesisOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		cText := C.CString(ssml)
		defer C.free(unsafe.Pointer(cText))
		length := len(ssml)
		ret := uintptr(C.synthesizer_start_speaking_ssml(synthesizer.handle, cText, (C.uint32_t)(length), &handle))
		if ret != C.SPX_NOERROR {
			outcome <- SpeechSynthesisOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewSpeechSynthesisResultFromHandle(handle2uintptr(handle))
			outcome <- SpeechSynthesisOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// StopSpeakingAsync stops the speech synthesis, asynchronously.
// It stops audio speech synthesis and discards any unread data in audio.PullAudioOutputStream.
func (synthesizer SpeechSynthesizer) StopSpeakingAsync() chan error {
	outcome := make(chan error)
	go func() {
		ret := uintptr(C.synthesizer_stop_speaking(synthesizer.handle))
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
		} else {
			outcome <- nil
		}
	}()
	return outcome
}

// GetVoicesAsync gets the available voices, asynchronously.
// The parameter locale specifies the locale of voices, in BCP-47 format; or leave it empty to get all available voices.
func (synthesizer SpeechSynthesizer) GetVoicesAsync(locale string) chan SpeechSynthesisVoicesOutcome {
	outcome := make(chan SpeechSynthesisVoicesOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		cLocale := C.CString(locale)
		defer C.free(unsafe.Pointer(cLocale))
		ret := uintptr(C.synthesizer_get_voices_list(synthesizer.handle, cLocale, &handle))
		if ret != C.SPX_NOERROR {
			outcome <- SpeechSynthesisVoicesOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewSynthesisVoicesResultFromHandle(handle2uintptr(handle))
			outcome <- SpeechSynthesisVoicesOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// SetAuthorizationToken sets the authorization token that will be used for connecting to the service.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token
// expires, the caller needs to refresh it by calling this setter with a new valid token.
// Otherwise, the synthesizer will encounter errors during synthesizing.
func (synthesizer SpeechSynthesizer) SetAuthorizationToken(token string) error {
	return synthesizer.Properties.SetProperty(common.SpeechServiceAuthorizationToken, token)
}

// AuthorizationToken is the authorization token.
func (synthesizer SpeechSynthesizer) AuthorizationToken() string {
	return synthesizer.Properties.GetProperty(common.SpeechServiceAuthorizationToken, "")
}

// SynthesisStarted signals events indicating the start of a synthesis
func (synthesizer SpeechSynthesizer) SynthesisStarted(handler SpeechSynthesisEventHandler) {
	registerSynthesisStartedCallback(handler, synthesizer.handle)
	if handler != nil {
		C.synthesizer_started_set_callback(
			synthesizer.handle,
			(C.PSYNTHESIS_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_synthesizer_synthesis_started)),
			nil)
	} else {
		C.synthesizer_started_set_callback(synthesizer.handle, nil, nil)
	}
}

// Synthesizing signals events indicating audio chunk is received while the synthesis is on going.
func (synthesizer SpeechSynthesizer) Synthesizing(handler SpeechSynthesisEventHandler) {
	registerSynthesizingCallback(handler, synthesizer.handle)
	if handler != nil {
		C.synthesizer_synthesizing_set_callback(
			synthesizer.handle,
			(C.PSYNTHESIS_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_synthesizer_synthesizing)),
			nil)
	} else {
		C.synthesizer_synthesizing_set_callback(synthesizer.handle, nil, nil)
	}
}

// SynthesisCompleted signals events indicating synthesis is completed.
func (synthesizer SpeechSynthesizer) SynthesisCompleted(handler SpeechSynthesisEventHandler) {
	registerSynthesisCompletedCallback(handler, synthesizer.handle)
	if handler != nil {
		C.synthesizer_completed_set_callback(
			synthesizer.handle,
			(C.PSYNTHESIS_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_synthesizer_synthesis_completed)),
			nil)
	} else {
		C.synthesizer_completed_set_callback(synthesizer.handle, nil, nil)
	}
}

// SynthesisCanceled signals that a speech synthesis result is received when the synthesis is canceled.
func (synthesizer SpeechSynthesizer) SynthesisCanceled(handler SpeechSynthesisEventHandler) {
	registerSynthesisCanceledCallback(handler, synthesizer.handle)
	if handler != nil {
		C.synthesizer_canceled_set_callback(
			synthesizer.handle,
			(C.PSYNTHESIS_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_synthesizer_synthesis_canceled)),
			nil)
	} else {
		C.synthesizer_canceled_set_callback(synthesizer.handle, nil, nil)
	}
}

// WordBoundary signals that a word boundary event is received.
func (synthesizer SpeechSynthesizer) WordBoundary(handler SpeechSynthesisWordBoundaryEventHandler) {
	registerSynthesisWordBoundaryCallback(handler, synthesizer.handle)
	if handler != nil {
		C.synthesizer_word_boundary_set_callback(
			synthesizer.handle,
			(C.PSYNTHESIS_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_synthesizer_word_boundary)),
			nil)
	} else {
		C.synthesizer_word_boundary_set_callback(synthesizer.handle, nil, nil)
	}
}

// VisemeReceived signals that a viseme event is received.
func (synthesizer SpeechSynthesizer) VisemeReceived(handler SpeechSynthesisVisemeEventHandler) {
	registerSynthesisVisemeReceivedCallback(handler, synthesizer.handle)
	if handler != nil {
		C.synthesizer_viseme_received_set_callback(
			synthesizer.handle,
			(C.PSYNTHESIS_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_synthesizer_viseme_received)),
			nil)
	} else {
		C.synthesizer_viseme_received_set_callback(synthesizer.handle, nil, nil)
	}
}

// BookmarkReached signals that a viseme event is received.
func (synthesizer SpeechSynthesizer) BookmarkReached(handler SpeechSynthesisBookmarkEventHandler) {
	registerSynthesisBookmarkReachedCallback(handler, synthesizer.handle)
	if handler != nil {
		C.synthesizer_bookmark_reached_set_callback(
			synthesizer.handle,
			(C.PSYNTHESIS_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_synthesizer_bookmark_reached)),
			nil)
	} else {
		C.synthesizer_bookmark_reached_set_callback(synthesizer.handle, nil, nil)
	}
}

// Close disposes the associated resources.
func (synthesizer *SpeechSynthesizer) Close() {
	synthesizer.SynthesisStarted(nil)
	synthesizer.Synthesizing(nil)
	synthesizer.SynthesisCompleted(nil)
	synthesizer.SynthesisCanceled(nil)
	synthesizer.WordBoundary(nil)
	synthesizer.VisemeReceived(nil)
	synthesizer.BookmarkReached(nil)
	synthesizer.Properties.Close()
	if synthesizer.handle != C.SPXHANDLE_INVALID {
		C.synthesizer_handle_release(synthesizer.handle)
		synthesizer.handle = C.SPXHANDLE_INVALID
	}
}
