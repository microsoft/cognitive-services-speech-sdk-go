// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_recognizer.h>
import "C"

// SpeechRecognitionEventArgs represents the speech recognition event arguments.
type SpeechRecognitionEventArgs struct {
	RecognitionEventArgs
	handle C.SPXHANDLE
	Result SpeechRecognitionResult
}

// Close releases the underlying resources
func (event SpeechRecognitionEventArgs) Close() {
	event.RecognitionEventArgs.Close()
	event.Result.Close()
}

// NewSpeechRecognitionEventArgsFromHandle creates the object from the handle (for internal use)
func NewSpeechRecognitionEventArgsFromHandle(handle common.SPXHandle) (*SpeechRecognitionEventArgs, error) {
	base, err := NewRecognitionEventArgsFromHandle(handle)
	if err != nil {
		return nil, err
	}
	event := new(SpeechRecognitionEventArgs)
	event.RecognitionEventArgs = *base
	event.handle = uintptr2handle(handle)
	var resultHandle C.SPXHANDLE
	ret := uintptr(C.recognizer_recognition_event_get_result(event.handle, &resultHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result, err := NewSpeechRecognitionResultFromHandle(handle2uintptr(resultHandle))
	if err != nil {
		return nil, err
	}
	event.Result = *result
	return event, nil
}

// SpeechRecognitionEventHandler is the type of the event handler that receives SpeechRecognitionEventArgs
type SpeechRecognitionEventHandler func(event SpeechRecognitionEventArgs)
