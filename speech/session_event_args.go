package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_recognizer.h>
import "C"
import "unsafe"

// SessionEventArgs represents the session event arguments.
type SessionEventArgs struct {
	sessionID string
}

// SessionID Session identifier (a GUID in string format).
func (event SessionEventArgs) SessionID() string {
	return event.sessionID
}

// NewSessionEventArgsFromHandle creates the object from the handle (for internal)
func NewSessionEventArgsFromHandle(handle SPXHandle) (*SessionEventArgs, error) {
	buffer := C.malloc(C.sizeof_char * 37)
	defer C.free(unsafe.Pointer(buffer))
	ret := uintptr(C.recognizer_session_event_get_session_id(uintptr2handle(handle), (*C.char)(buffer), 37))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	event := new(SessionEventArgs)
	event.sessionID = string(C.GoBytes(buffer, 37))
	return event, nil
}

// SessionEventHandler is the type of the event handler that receives SessionEventArgs
type SessionEventHandler func (event SessionEventArgs)