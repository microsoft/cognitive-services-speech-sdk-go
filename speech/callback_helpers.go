//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package speech

import (
	"sync"
)

// #include <speechapi_c_common.h>
import "C"

var mu sync.Mutex
var sessionStartedCallbacks = make(map[C.SPXHANDLE]SessionEventHandler)

func registerSessionStartedCallback(handler SessionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	sessionStartedCallbacks[handle] = handler
}

func getSessionStartedCallback(handle C.SPXHANDLE) SessionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return sessionStartedCallbacks[handle]
}

//export recognizerFireEventSessionStarted
func recognizerFireEventSessionStarted(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getSessionStartedCallback(handle)
	if handler == nil {
		return
	}
	event, err := NewSessionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}

var sessionStoppedCallbacks = make(map[C.SPXHANDLE]SessionEventHandler)

func registerSessionStoppedCallback(handler SessionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	sessionStoppedCallbacks[handle] = handler
}

func getSessionStoppedCallback(handle C.SPXHANDLE) SessionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return sessionStoppedCallbacks[handle]
}

//export recognizerFireEventSessionStopped
func recognizerFireEventSessionStopped(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getSessionStoppedCallback(handle)
	if handler == nil {
		return
	}
	event, err := NewSessionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}

var speechStartDetectedCallbacks = make(map[C.SPXHANDLE]RecognitionEventHandler)

func registerSpeechStartDetectedCallback(handler RecognitionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	speechStartDetectedCallbacks[handle] = handler
}

func getSpeechStartDetectedCallback(handle C.SPXHANDLE) RecognitionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return speechStartDetectedCallbacks[handle]
}

//export recognizerFireEventSpeechStartDetected
func recognizerFireEventSpeechStartDetected(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getSpeechStartDetectedCallback(handle)
	if handler == nil {
		return
	}
	event, err := NewRecognitionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}

var speechEndDetectedCallbacks = make(map[C.SPXHANDLE]RecognitionEventHandler)

func registerSpeechEndDetectedCallback(handler RecognitionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	speechEndDetectedCallbacks[handle] = handler
}

func getSpeechEndDetectedCallback(handle C.SPXHANDLE) RecognitionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return speechEndDetectedCallbacks[handle]
}

//export recognizerFireEventSpeechEndDetected
func recognizerFireEventSpeechEndDetected(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getSpeechEndDetectedCallback(handle)
	if handler == nil {
		return
	}
	event, err := NewRecognitionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}

var recognizedCallbacks = make(map[C.SPXHANDLE]SpeechRecognitionEventHandler)

func registerRecognizedCallback(handler SpeechRecognitionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	recognizedCallbacks[handle] = handler
}

func getRecognizedCallback(handle C.SPXHANDLE) SpeechRecognitionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return recognizedCallbacks[handle]
}

//export recognizerFireEventRecognized
func recognizerFireEventRecognized(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getRecognizedCallback(handle)
	if handler == nil {
		return
	}
	event, err := NewSpeechRecognitionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}

var recognizingCallbacks = make(map[C.SPXHANDLE]SpeechRecognitionEventHandler)

func registerRecognizingCallback(handler SpeechRecognitionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	recognizingCallbacks[handle] = handler
}

func getRecognizingCallback(handle C.SPXHANDLE) SpeechRecognitionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return recognizingCallbacks[handle]
}

//export recognizerFireEventRecognizing
func recognizerFireEventRecognizing(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getRecognizingCallback(handle)
	if handler == nil {
		return
	}
	event, err := NewSpeechRecognitionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}

var canceledCallbacks = make(map[C.SPXHANDLE]SpeechRecognitionCanceledEventHandler)

func registerCanceledCallback(handler SpeechRecognitionCanceledEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	canceledCallbacks[handle] = handler
}

func getCanceledCallback(handle C.SPXHANDLE) SpeechRecognitionCanceledEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return canceledCallbacks[handle]
}

//export recognizerFireEventCanceled
func recognizerFireEventCanceled(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getCanceledCallback(handle)
	if handler == nil {
		return
	}
	event, err := NewSpeechRecognitionCanceledEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil {
		return
	}
	handler(*event)
}
