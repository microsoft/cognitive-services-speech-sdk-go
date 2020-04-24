// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_recognizer.h>
import "C"

// RecognitionEventArgs represents the recognition event arguments.
type RecognitionEventArgs struct {
	SessionEventArgs
	Offset uint64
}

// NewRecognitionEventArgsFromHandle creates the object from the handle (for internal use)
func NewRecognitionEventArgsFromHandle(handle common.SPXHandle) (*RecognitionEventArgs, error) {
	base, err := NewSessionEventArgsFromHandle(handle)
	if err != nil {
		return nil, err
	}
	var offset C.uint64_t
	ret := uintptr(C.recognizer_recognition_event_get_offset(uintptr2handle(handle), &offset))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	event := new(RecognitionEventArgs)
	event.SessionEventArgs = *base
	event.Offset = uint64(offset)
	return event, nil
}

// RecognitionEventHandler is the type of the event handler that receives RecognitionEventArgs
type RecognitionEventHandler func(event RecognitionEventArgs)
