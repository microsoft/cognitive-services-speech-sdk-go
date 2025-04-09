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
//
// /* Proxy functions forward declarations */
// void cgo_conversation_transcriber_session_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_conversation_transcriber_session_stopped(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_conversation_transcriber_speech_start_detected(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_conversation_transcriber_speech_end_detected(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_conversation_transcriber_transcribing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_conversation_transcriber_transcribed(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_conversation_transcriber_canceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
//
import "C"

// ConversationTranscriber is the class for conversation transcribers.
type ConversationTranscriber struct {
	Properties                 *common.PropertyCollection
	handle                     C.SPXHANDLE
	handleAsyncStartTranscribing C.SPXASYNCHANDLE
	handleAsyncStopTranscribing  C.SPXASYNCHANDLE
}

func newConversationTranscriberFromHandle(handle C.SPXHANDLE) (*ConversationTranscriber, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.recognizer_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	transcriber := new(ConversationTranscriber)
	transcriber.handle = handle
	transcriber.handleAsyncStartTranscribing = C.SPXHANDLE_INVALID
	transcriber.handleAsyncStopTranscribing = C.SPXHANDLE_INVALID
	transcriber.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	
	return transcriber, nil
}

// NewConversationTranscriberFromConfig creates a conversation transcriber from a speech config and audio config.
func NewConversationTranscriberFromConfig(config *SpeechConfig, audioConfig *audio.AudioConfig) (*ConversationTranscriber, error) {
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
	
	ret := uintptr(C.recognizer_create_conversation_transcriber_from_config(&handle, configHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	
	return newConversationTranscriberFromHandle(handle)
}

// NewConversationTranscriberFromAutoDetectSourceLangConfig creates a conversation transcriber with auto language detection
func NewConversationTranscriberFromAutoDetectSourceLangConfig(config *SpeechConfig, langConfig *AutoDetectSourceLanguageConfig, audioConfig *audio.AudioConfig) (*ConversationTranscriber, error) {
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
	
	ret := uintptr(C.recognizer_create_conversation_transcriber_from_auto_detect_source_lang_config(&handle, configHandle, langConfigHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	
	return newConversationTranscriberFromHandle(handle)
}

// NewConversationTranscriberFromSourceLanguageConfig creates a conversation transcriber with a specific source language
func NewConversationTranscriberFromSourceLanguageConfig(config *SpeechConfig, sourceLanguageConfig *SourceLanguageConfig, audioConfig *audio.AudioConfig) (*ConversationTranscriber, error) {
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
	
	ret := uintptr(C.recognizer_create_conversation_transcriber_from_source_lang_config(&handle, configHandle, langConfigHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	
	return newConversationTranscriberFromHandle(handle)
}

// StartTranscribingAsync asynchronously initiates continuous conversation transcription.
func (transcriber ConversationTranscriber) StartTranscribingAsync() chan error {
	outcome := make(chan error)
	
	go func() {
		// Close any unfinished previous attempt
		ret := releaseAsyncHandleIfValid(&transcriber.handleAsyncStartTranscribing)
		
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_start_continuous_recognition_async(transcriber.handle, &transcriber.handleAsyncStartTranscribing))
		}
		
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_start_continuous_recognition_async_wait_for(transcriber.handleAsyncStartTranscribing, math.MaxUint32))
		}
		
		releaseAsyncHandleIfValid(&transcriber.handleAsyncStartTranscribing)
		
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
			return
		}
		
		outcome <- nil
	}()
	
	return outcome
}

// StopTranscribingAsync asynchronously terminates ongoing continuous conversation transcription.
func (transcriber ConversationTranscriber) StopTranscribingAsync() chan error {
	outcome := make(chan error)
	
	go func() {
		ret := releaseAsyncHandleIfValid(&transcriber.handleAsyncStopTranscribing)
		
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_stop_continuous_recognition_async(transcriber.handle, &transcriber.handleAsyncStopTranscribing))
		}
		
		if ret == C.SPX_NOERROR {
			ret = uintptr(C.recognizer_stop_continuous_recognition_async_wait_for(transcriber.handleAsyncStopTranscribing, math.MaxUint32))
		}
		
		releaseAsyncHandleIfValid(&transcriber.handleAsyncStopTranscribing)
		
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
			return
		}
		
		outcome <- nil
	}()
	
	return outcome
}

// SetAuthorizationToken sets the authorization token that will be used for connecting to the service.
func (transcriber ConversationTranscriber) SetAuthorizationToken(token string) error {
	return transcriber.Properties.SetProperty(common.SpeechServiceAuthorizationToken, token)
}

// AuthorizationToken is the authorization token.
func (transcriber ConversationTranscriber) AuthorizationToken() string {
	return transcriber.Properties.GetProperty(common.SpeechServiceAuthorizationToken, "")
}

// SessionStarted signals events indicating the start of a recognition session (operation).
func (transcriber ConversationTranscriber) SessionStarted(handler SessionEventHandler) {
	registerSessionStartedCallback(handler, transcriber.handle)
	if handler != nil {
		C.recognizer_session_started_set_callback(
			transcriber.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_conversation_transcriber_session_started)),
			nil)
	} else {
		C.recognizer_session_started_set_callback(transcriber.handle, nil, nil)
	}
}

// SessionStopped signals events indicating the end of a recognition session (operation).
func (transcriber ConversationTranscriber) SessionStopped(handler SessionEventHandler) {
	registerSessionStoppedCallback(handler, transcriber.handle)
	if handler != nil {
		C.recognizer_session_stopped_set_callback(
			transcriber.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_conversation_transcriber_session_stopped)),
			nil)
	} else {
		C.recognizer_session_stopped_set_callback(transcriber.handle, nil, nil)
	}
}

// SpeechStartDetected signals for events indicating the start of speech.
func (transcriber ConversationTranscriber) SpeechStartDetected(handler RecognitionEventHandler) {
	registerSpeechStartDetectedCallback(handler, transcriber.handle)
	if handler != nil {
		C.recognizer_speech_start_detected_set_callback(
			transcriber.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_conversation_transcriber_speech_start_detected)),
			nil)
	} else {
		C.recognizer_speech_start_detected_set_callback(transcriber.handle, nil, nil)
	}
}

// SpeechEndDetected signals for events indicating the end of speech.
func (transcriber ConversationTranscriber) SpeechEndDetected(handler RecognitionEventHandler) {
	registerSpeechEndDetectedCallback(handler, transcriber.handle)
	if handler != nil {
		C.recognizer_speech_end_detected_set_callback(
			transcriber.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_conversation_transcriber_speech_end_detected)),
			nil)
	} else {
		C.recognizer_speech_end_detected_set_callback(transcriber.handle, nil, nil)
	}
}

// Transcribing signals for events containing intermediate transcription results.
func (transcriber ConversationTranscriber) Transcribing(handler ConversationTranscriptionEventHandler) {
	registerConversationTranscribingCallback(handler, transcriber.handle)
	if handler != nil {
		C.recognizer_recognizing_set_callback(
			transcriber.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_conversation_transcriber_transcribing)),
			nil)
	} else {
		C.recognizer_recognizing_set_callback(transcriber.handle, nil, nil)
	}
}

// Transcribed signals for events containing final transcription results.
func (transcriber ConversationTranscriber) Transcribed(handler ConversationTranscriptionEventHandler) {
	registerConversationTranscribedCallback(handler, transcriber.handle)
	if handler != nil {
		C.recognizer_recognized_set_callback(
			transcriber.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_conversation_transcriber_transcribed)),
			nil)
	} else {
		C.recognizer_recognized_set_callback(transcriber.handle, nil, nil)
	}
}

// Canceled signals for events containing canceled transcription results.
func (transcriber ConversationTranscriber) Canceled(handler ConversationTranscriptionCanceledEventHandler) {
	registerConversationCanceledCallback(handler, transcriber.handle)
	if handler != nil {
		C.recognizer_canceled_set_callback(
			transcriber.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_conversation_transcriber_canceled)),
			nil)
	} else {
		C.recognizer_canceled_set_callback(transcriber.handle, nil, nil)
	}
}

// Close disposes the associated resources.
func (transcriber ConversationTranscriber) Close() {
	transcriber.SessionStarted(nil)
	transcriber.SessionStopped(nil)
	transcriber.SpeechStartDetected(nil)
	transcriber.SpeechEndDetected(nil)
	transcriber.Transcribing(nil)
	transcriber.Transcribed(nil)
	transcriber.Canceled(nil)
	
	var asyncHandles = []*C.SPXASYNCHANDLE{
		&transcriber.handleAsyncStartTranscribing,
		&transcriber.handleAsyncStopTranscribing,
	}
	
	for i := 0; i < len(asyncHandles); i++ {
		handle := asyncHandles[i]
		releaseAsyncHandleIfValid(handle)
	}
	
	transcriber.Properties.Close()
	
	if transcriber.handle != C.SPXHANDLE_INVALID {
		C.recognizer_handle_release(transcriber.handle)
		transcriber.handle = C.SPXHANDLE_INVALID
	}
}