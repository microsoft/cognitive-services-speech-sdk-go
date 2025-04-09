// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_result.h>
// #include <speechapi_c_conversation_transcription_result.h>
import "C"

// ConversationTranscriptionResult contains detailed information about result of a conversation transcription operation.
type ConversationTranscriptionResult struct {
	SpeechRecognitionResult // Embedded for common fields
	SpeakerId              string
}

// NewConversationTranscriptionResultFromHandle creates a ConversationTranscriptionResult from a handle (for internal use)
func NewConversationTranscriptionResultFromHandle(handle common.SPXHandle) (*ConversationTranscriptionResult, error) {
	// Create base result first
	baseResult, err := NewSpeechRecognitionResultFromHandle(handle)
	if err != nil {
		return nil, err
	}

	result := &ConversationTranscriptionResult{
		SpeechRecognitionResult: *baseResult,
	}

	// Get speaker ID
	buffer := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(buffer))
	
	ret := uintptr(C.conversation_transcription_result_get_speaker_id(result.handle, (*C.char)(buffer), 1024))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.SpeakerId = C.GoString((*C.char)(buffer))

	return result, nil
}

// Close releases the underlying resources
func (result ConversationTranscriptionResult) Close() {
	// Only call the base Close since we don't have additional resources to clean up
	result.SpeechRecognitionResult.Close()
}