// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_synthesizer.h>
import "C"

// SpeechSynthesisEventArgs represents the speech synthesis event arguments.
type SpeechSynthesisEventArgs struct {
	handle C.SPXHANDLE
	Result SpeechSynthesisResult
}

// Close releases the underlying resources
func (event SpeechSynthesisEventArgs) Close() {
	event.Result.Close()
	C.synthesizer_event_handle_release(event.handle)
}

// NewSpeechSynthesisEventArgsFromHandle creates the object from the handle (for internal use)
func NewSpeechSynthesisEventArgsFromHandle(handle common.SPXHandle) (*SpeechSynthesisEventArgs, error) {
	event := new(SpeechSynthesisEventArgs)
	event.handle = uintptr2handle(handle)
	var resultHandle C.SPXHANDLE
	ret := uintptr(C.synthesizer_synthesis_event_get_result(event.handle, &resultHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result, err := NewSpeechSynthesisResultFromHandle(handle2uintptr(resultHandle))
	if err != nil {
		return nil, err
	}
	event.Result = *result
	return event, nil
}

// SpeechSynthesisEventHandler is the type of the event handler that receives SpeechSynthesisEventArgs
type SpeechSynthesisEventHandler func(event SpeechSynthesisEventArgs)
