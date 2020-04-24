// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_result.h>
import "C"

// SpeechRecognitionCanceledEventArgs represents speech recognition canceled event arguments.
type SpeechRecognitionCanceledEventArgs struct {
	SpeechRecognitionEventArgs
	Reason       common.CancellationReason
	ErrorCode    common.CancellationErrorCode
	ErrorDetails string
}

// NewSpeechRecognitionCanceledEventArgsFromHandle creates the object from the handle (for internal use)
func NewSpeechRecognitionCanceledEventArgsFromHandle(handle common.SPXHandle) (*SpeechRecognitionCanceledEventArgs, error) {
	baseArgs, err := NewSpeechRecognitionEventArgsFromHandle(handle)
	if err != nil {
		return nil, err
	}
	event := new(SpeechRecognitionCanceledEventArgs)
	event.SpeechRecognitionEventArgs = *baseArgs
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

// SpeechRecognitionCanceledEventHandler is the type of the event handler that receives SpeechRecognitionCanceledEventArgs
type SpeechRecognitionCanceledEventHandler func(event SpeechRecognitionCanceledEventArgs)
