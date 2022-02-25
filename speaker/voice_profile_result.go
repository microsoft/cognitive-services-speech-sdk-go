// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
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

// VoiceProfileResult contains information about result from voice profile operations.
type VoiceProfileResult struct {
	handle C.SPXHANDLE

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
func (result VoiceProfileResult) Close() {
	result.Properties.Close()
	C.recognizer_result_handle_release(result.handle)
}

// newVoiceProfileResultFromHandle creates a VoiceProfileResult from a handle (for internal use)
func newVoiceProfileResultFromHandle (handle common.SPXHandle) (*VoiceProfileResult, error) {
	result := new(VoiceProfileResult)
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
	return result, nil
}

// VoiceProfileOutcome is a wrapper type to be returned by operations returning VoiceProfileResult and error
type VoiceProfileOutcome struct {
	common.OperationOutcome

	// Result is the result of the operation
	Result *VoiceProfileResult
}

// Close releases the underlying resources
func (outcome VoiceProfileOutcome) Close() {
	if outcome.Result != nil {
		outcome.Result.Close()
	}
}
