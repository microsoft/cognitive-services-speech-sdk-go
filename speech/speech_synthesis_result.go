// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"time"
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_result.h>
// #include <speechapi_c_synthesizer.h>
//
import "C"

// SpeechSynthesisResult contains detailed information about result of a synthesis operation.
type SpeechSynthesisResult struct {
	handle C.SPXHANDLE

	// ResultID specifies the result identifier.
	ResultID string

	// Reason specifies status of speech synthesis result.
	Reason common.ResultReason

	// AudioData presents the synthesized audio.
	AudioData []byte

	// AudioDuration presents the time duration of synthesized audio.
	AudioDuration time.Duration

	// Collection of additional synthesisResult properties.
	Properties *common.PropertyCollection
}

// Close releases the underlying resources
func (result *SpeechSynthesisResult) Close() {
	result.Properties.Close()
	if result.handle != C.SPXHANDLE_INVALID {
		C.synthesizer_result_handle_release(result.handle)
		result.handle = C.SPXHANDLE_INVALID
	}
}

// NewSpeechSynthesisResultFromHandle creates a SpeechSynthesisResult from a handle (for internal use)
func NewSpeechSynthesisResultFromHandle(handle common.SPXHandle) (*SpeechSynthesisResult, error) {

	result := new(SpeechSynthesisResult)
	result.handle = uintptr2handle(handle)
	/* AudioData length and duration */
	var cAudioLength C.uint32_t
	var cAudioDuration C.uint64_t
	ret := uintptr(C.synth_result_get_audio_length_duration(result.handle, &cAudioLength, &cAudioDuration))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.AudioDuration = time.Duration(cAudioDuration) * time.Millisecond
	// using max(1024, cAudioLength) as buffer size
	if cAudioLength < 1024 {
		cAudioLength = 1024
	}
	buffer := C.malloc(C.sizeof_char * (C.size_t)(cAudioLength))
	defer C.free(unsafe.Pointer(buffer))
	/* ResultID */
	ret = uintptr(C.synth_result_get_result_id(result.handle, (*C.char)(buffer), 1024))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.ResultID = C.GoString((*C.char)(buffer))
	/* Reason */
	var cReason C.Result_Reason
	ret = uintptr(C.synth_result_get_reason(result.handle, &cReason))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Reason = (common.ResultReason)(cReason)
	/* AudioData */
	var outSize C.uint32_t
	ret = uintptr(C.synth_result_get_audio_data(result.handle, (*C.uint8_t)(buffer), cAudioLength, &outSize))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.AudioData = C.GoBytes(buffer, (C.int)(outSize))
	/* Properties */
	var propBagHandle C.SPXHANDLE
	ret = uintptr(C.synth_result_get_property_bag(uintptr2handle(handle), &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return result, nil
}

// SpeechSynthesisOutcome is a wrapper type to be returned by operations returning SpeechSynthesisResult and error
type SpeechSynthesisOutcome struct {
	common.OperationOutcome

	// Result is the result of the operation
	Result *SpeechSynthesisResult
}

// Close releases the underlying resources
func (outcome SpeechSynthesisOutcome) Close() {
	if outcome.Result != nil {
		outcome.Result.Close()
	}
}
