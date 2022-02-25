// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	"strings"
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_result.h>
// #include <speechapi_c_recognizer.h>
//
import "C"

// VoiceProfilePhraseResult contains activation phrases needed to successfully enroll a voice profile.
type VoiceProfilePhraseResult struct {
	handle C.SPXHANDLE

	// Activation phrases for voice profile enrollment 
	Phrases []string

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
func (result VoiceProfilePhraseResult) Close() {
	result.Properties.Close()
	C.recognizer_result_handle_release(result.handle)
}

// newVoiceProfilePhraseResultFromHandle creates a VoiceProfilePhraseResult from a handle (for internal use)
func newVoiceProfilePhraseResultFromHandle (handle common.SPXHandle) (*VoiceProfilePhraseResult, error) {
	result := new(VoiceProfilePhraseResult)
	result.handle = uintptr2handle(handle)
	buffer := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(buffer))
	/* ResultID */
	ret := uintptr(C.result_get_result_id(result.handle, (*C.char)(buffer), 1024))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.ResultID = C.GoString((*C.char)(buffer))
	/* Reason */
	var cReason C.Result_Reason
	ret = uintptr(C.result_get_reason(result.handle, &cReason))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Reason = (common.ResultReason)(cReason)
	/* Properties */
	var propBagHandle C.SPXHANDLE
	ret = uintptr(C.result_get_property_bag(result.handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	result.ErrorDetails = result.Properties.GetProperty(common.CancellationDetailsReasonDetailedText, "")

	/* Phrases */
	phrasesString := result.Properties.GetPropertyByString("speakerrecognition.phrases", "")
	if len(phrasesString) > 0 {
		result.Phrases = strings.Split(phrasesString, "|")
	}
	return result, nil
}

// VoiceProfilePhraseOutcome is a wrapper type to be returned by operations returning VoiceProfilePhraseResult and error
type VoiceProfilePhraseOutcome struct {
	common.OperationOutcome

	// Result is the result of the operation
	Result *VoiceProfilePhraseResult
}

// Close releases the underlying resources
func (outcome VoiceProfilePhraseOutcome) Close() {
	if outcome.Result != nil {
		outcome.Result.Close()
	}
}
