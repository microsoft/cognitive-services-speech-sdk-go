package dialog

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"sync"
)

// #include <speechapi_c_common.h>
import "C"

var mu sync.Mutex
var sessionStartedCallbacks = make(map[C.SPXHANDLE]speech.SessionEventHandler)

func registerSessionStartedCallback(handler speech.SessionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	sessionStartedCallbacks[handle] = handler
}

func getSessionStartedCallback(handle C.SPXHANDLE) speech.SessionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return sessionStartedCallbacks[handle]
}

//export dialogFireEventSessionStarted
func dialogFireEventSessionStarted(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getSessionStartedCallback(handle)
	if handler == nil {
		return
	}
	event, err := speech.NewSessionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}

var sessionStoppedCallbacks = make(map[C.SPXHANDLE]speech.SessionEventHandler)

func registerSessionStoppedCallback(handler speech.SessionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	sessionStoppedCallbacks[handle] = handler;
}

func getSessionStoppedCallback(handle C.SPXHANDLE) speech.SessionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return sessionStoppedCallbacks[handle]
}

//export dialogFireEventSessionStopped
func dialogFireEventSessionStopped(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getSessionStoppedCallback(handle)
	if handler == nil {
		return
	}
	event, err := speech.NewSessionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}