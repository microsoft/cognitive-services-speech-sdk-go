// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_recognizer.h>
import "C"

// SessionEventArgs represents the session event arguments.
type SessionEventArgs struct {
	handle C.SPXHANDLE
	// SessionID Session identifier (a GUID in string format).
	SessionID string
}

// Close releases the underlying resources.
func (event SessionEventArgs) Close() {
	C.recognizer_event_handle_release(event.handle)
}

// NewSessionEventArgsFromHandle creates the object from the handle (for internal use)
func NewSessionEventArgsFromHandle(handle common.SPXHandle) (*SessionEventArgs, error) {
	buffer := C.malloc(C.sizeof_char * 37)
	defer C.free(unsafe.Pointer(buffer))
	ret := uintptr(C.recognizer_session_event_get_session_id(uintptr2handle(handle), (*C.char)(buffer), 37))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	event := new(SessionEventArgs)
	event.handle = uintptr2handle(handle)
	event.SessionID = C.GoString((*C.char)(buffer))
	return event, nil
}

// SessionEventHandler is the type of the event handler that receives SessionEventArgs
type SessionEventHandler func(event SessionEventArgs)
