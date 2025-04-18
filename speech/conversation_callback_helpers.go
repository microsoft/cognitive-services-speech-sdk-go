// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

// #include <speechapi_c_common.h>
// #include <speechapi_c_recognizer.h>
import "C"

// ConversationTranscriptionEventHandler is the callback type for conversation transcription events.
type ConversationTranscriptionEventHandler func(event ConversationTranscriptionEventArgs)

// ConversationTranscriptionCanceledEventHandler is the callback type for conversation transcription canceled events.
type ConversationTranscriptionCanceledEventHandler func(event ConversationTranscriptionCanceledEventArgs)

var (
	conversationTranscribingCallbacks = make(map[C.SPXHANDLE]ConversationTranscriptionEventHandler)
	conversationTranscribedCallbacks  = make(map[C.SPXHANDLE]ConversationTranscriptionEventHandler)
	conversationCanceledCallbacks     = make(map[C.SPXHANDLE]ConversationTranscriptionCanceledEventHandler)
)

func registerConversationTranscribingCallback(handler ConversationTranscriptionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	conversationTranscribingCallbacks[handle] = handler
}

func getConversationTranscribingCallback(handle C.SPXHANDLE) ConversationTranscriptionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return conversationTranscribingCallbacks[handle]
}

//export conversationTranscriberFireEventTranscribing
func conversationTranscriberFireEventTranscribing(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getConversationTranscribingCallback(handle)
	event, err := NewConversationTranscriptionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil || handler == nil {
		C.recognizer_event_handle_release(eventHandle)
		return
	}
	handler(*event)
}

func registerConversationTranscribedCallback(handler ConversationTranscriptionEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	conversationTranscribedCallbacks[handle] = handler
}

func getConversationTranscribedCallback(handle C.SPXHANDLE) ConversationTranscriptionEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return conversationTranscribedCallbacks[handle]
}

//export conversationTranscriberFireEventTranscribed
func conversationTranscriberFireEventTranscribed(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getConversationTranscribedCallback(handle)
	event, err := NewConversationTranscriptionEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil || handler == nil {
		C.recognizer_event_handle_release(eventHandle)
		return
	}
	handler(*event)
}

func registerConversationCanceledCallback(handler ConversationTranscriptionCanceledEventHandler, handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	conversationCanceledCallbacks[handle] = handler
}

func getConversationCanceledCallback(handle C.SPXHANDLE) ConversationTranscriptionCanceledEventHandler {
	mu.Lock()
	defer mu.Unlock()
	return conversationCanceledCallbacks[handle]
}

//export conversationTranscriberFireEventCanceled
func conversationTranscriberFireEventCanceled(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	handler := getConversationCanceledCallback(handle)
	event, err := NewConversationTranscriptionCanceledEventArgsFromHandle(handle2uintptr(eventHandle))
	if err != nil || handler == nil {
		C.recognizer_event_handle_release(eventHandle)
		return
	}
	handler(*event)
}