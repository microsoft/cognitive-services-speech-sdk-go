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
// #include <speechapi_c_synthesizer.h>
//
import "C"

// SynthesisVoicesResult contains information about result from voices list of speech synthesizers.
type SynthesisVoicesResult struct {
	handle C.SPXHANDLE

	// Voices specifies all voices retrieved
	Voices []*VoiceInfo

	// ResultID specifies the result identifier.
	ResultID string

	// Reason specifies status of speech synthesis result.
	Reason common.ResultReason

	// ErrorDetails presents error details.
	ErrorDetails string

	// Collection of additional properties.
	Properties *common.PropertyCollection
}

// Close releases the underlying resources
func (result SynthesisVoicesResult) Close() {
	for _, voice := range result.Voices {
		voice.Close()
	}
	result.Properties.Close()
	C.synthesizer_result_handle_release(result.handle)
}

// NewSynthesisVoicesResultFromHandle creates a SynthesisVoicesResult from a handle (for internal use)
func NewSynthesisVoicesResultFromHandle(handle common.SPXHandle) (*SynthesisVoicesResult, error) {
	result := new(SynthesisVoicesResult)
	result.handle = uintptr2handle(handle)
	buffer := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(buffer))
	/* ResultID */
	ret := uintptr(C.synthesis_voices_result_get_result_id(result.handle, (*C.char)(buffer), 1024))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.ResultID = C.GoString((*C.char)(buffer))
	/* Reason */
	var cReason C.Result_Reason
	ret = uintptr(C.synthesis_voices_result_get_reason(result.handle, &cReason))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Reason = (common.ResultReason)(cReason)
	/* Properties */
	var propBagHandle C.SPXHANDLE
	ret = uintptr(C.synthesis_voices_result_get_property_bag(result.handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	result.ErrorDetails = result.Properties.GetProperty(common.CancellationDetailsReasonDetailedText, "")
	/* Voices */
	var voiceNum C.uint32_t
	ret = uintptr(C.synthesis_voices_result_get_voice_num(result.handle, &voiceNum))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	voices := make([]*VoiceInfo, voiceNum)
	var voice *VoiceInfo
	var hVoice C.SPXRESULTHANDLE
	var err error
	for i := 0; i < int(voiceNum); i++ {
		ret = uintptr(C.synthesis_voices_result_get_voice_info(result.handle, (C.uint32_t)(i), &hVoice))
		if ret != C.SPX_NOERROR {
			return nil, common.NewCarbonError(ret)
		}
		voice, err = NewVoiceInfoFromHandle(handle2uintptr(hVoice))
		if err != nil {
			return nil, err
		}
		voices[i] = voice
	}
	result.Voices = voices
	return result, nil
}

// SpeechSynthesisVoicesOutcome is a wrapper type to be returned by operations returning SynthesisVoicesResult and error
type SpeechSynthesisVoicesOutcome struct {
	common.OperationOutcome

	// Result is the result of the operation
	Result *SynthesisVoicesResult
}

// Close releases the underlying resources
func (outcome SpeechSynthesisVoicesOutcome) Close() {
	if outcome.Result != nil {
		outcome.Result.Close()
	}
}
