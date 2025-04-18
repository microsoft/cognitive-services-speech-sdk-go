// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_result.h>
// #include <speechapi_c_recognizer.h>
import "C"

// ConversationTranscriptionEventArgs is used for conversation transcription events.
type ConversationTranscriptionEventArgs struct {
	RecognitionEventArgs      // Inherit from RecognitionEventArgs for consistency
	handle C.SPXHANDLE
	Result ConversationTranscriptionResult // Direct field instead of pointer
}

// NewConversationTranscriptionEventArgsFromHandle creates a ConversationTranscriptionEventArgs from an event handle
func NewConversationTranscriptionEventArgsFromHandle(handle common.SPXHandle) (*ConversationTranscriptionEventArgs, error) {
	// Create the base RecognitionEventArgs first
	base, err := NewRecognitionEventArgsFromHandle(handle)
	if err != nil {
		return nil, err
	}
	
	event := new(ConversationTranscriptionEventArgs)
	event.RecognitionEventArgs = *base
	event.handle = uintptr2handle(handle)
	
	// Get the result handle
	var resultHandle C.SPXHANDLE
	ret := uintptr(C.recognizer_recognition_event_get_result(event.handle, &resultHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	
	// Create the result
	result, err := NewConversationTranscriptionResultFromHandle(handle2uintptr(resultHandle))
	if err != nil {
		return nil, err
	}
	
	event.Result = *result
	return event, nil
}

// Close releases the underlying resources
func (event ConversationTranscriptionEventArgs) Close() {
	event.RecognitionEventArgs.Close()
	event.Result.Close()
}

// ConversationTranscriptionEventHandler is the type of the event handler that receives ConversationTranscriptionEventArgs
// type ConversationTranscriptionEventHandler func(event ConversationTranscriptionEventArgs)

// ConversationTranscriptionCanceledEventArgs is used for conversation transcription canceled events.
type ConversationTranscriptionCanceledEventArgs struct {
	ConversationTranscriptionEventArgs
	Reason       common.CancellationReason    // Direct field instead of nested object
	ErrorCode    common.CancellationErrorCode // Direct field instead of nested object
	ErrorDetails string                       // Direct field instead of nested object
}

// NewConversationTranscriptionCanceledEventArgsFromHandle creates a ConversationTranscriptionCanceledEventArgs from an event handle
func NewConversationTranscriptionCanceledEventArgsFromHandle(handle common.SPXHandle) (*ConversationTranscriptionCanceledEventArgs, error) {
	baseArgs, err := NewConversationTranscriptionEventArgsFromHandle(handle)
	if err != nil {
		return nil, err
	}

	event := new(ConversationTranscriptionCanceledEventArgs)
	event.ConversationTranscriptionEventArgs = *baseArgs
	
	/* Reason */
	var cReason C.Result_CancellationReason
	ret := uintptr(C.result_get_reason_canceled(event.Result.handle, &cReason))
	if ret != C.SPX_NOERROR {
		event.Close()
		return nil, common.NewCarbonError(ret)
	}
	event.Reason = (common.CancellationReason)(cReason)
	
	/* ErrorCode */
	var cCode C.Result_CancellationErrorCode
	ret = uintptr(C.result_get_canceled_error_code(event.Result.handle, &cCode))
	if ret != C.SPX_NOERROR {
		event.Close()
		return nil, common.NewCarbonError(ret)
	}
	event.ErrorCode = (common.CancellationErrorCode)(cCode)
	event.ErrorDetails = event.Result.Properties.GetProperty(common.SpeechServiceResponseJSONErrorDetails, "")
	
	return event, nil
}

// Close releases the associated resources.
func (event ConversationTranscriptionCanceledEventArgs) Close() {
	event.ConversationTranscriptionEventArgs.Close()
}

// ConversationTranscriptionCanceledEventHandler is the type of the event handler that receives ConversationTranscriptionCanceledEventArgs 
//type ConversationTranscriptionCanceledEventHandler func(event ConversationTranscriptionCanceledEventArgs)