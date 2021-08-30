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
// #include <speechapi_c_recognizer.h>
//
import "C"

// SpeechRecognitionResult contains detailed information about result of a recognition operation.
type SpeechRecognitionResult struct {
	handle C.SPXHANDLE

	// ResultID specifies the result identifier.
	ResultID string

	// Reason specifies status of speech recognition result.
	Reason common.ResultReason

	// Text presents the recognized text in the result.
	Text string

	// Duration of the recognized speech.
	Duration time.Duration

	// Offset of the recognized speech in ticks.
	Offset time.Duration

	// Collection of additional RecognitionResult properties.
	Properties *common.PropertyCollection
}

// Close releases the underlying resources
func (result SpeechRecognitionResult) Close() {
	result.Properties.Close()
	C.recognizer_result_handle_release(result.handle)
}

// NewSpeechRecognitionResultFromHandle creates a SpeechRecognitionResult from a handle (for internal use)
func NewSpeechRecognitionResultFromHandle(handle common.SPXHandle) (*SpeechRecognitionResult, error) {
	buffer := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(buffer))
	result := new(SpeechRecognitionResult)
	result.handle = uintptr2handle(handle)
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
	/* Text */
	ret = uintptr(C.result_get_text(result.handle, (*C.char)(buffer), 1024))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Text = C.GoString((*C.char)(buffer))
	/* Duration */
	var cDuration C.uint64_t
	ret = uintptr(C.result_get_duration(result.handle, &cDuration))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Duration = time.Nanosecond * time.Duration(100*cDuration)
	/* Offset */
	var cOffset C.uint64_t
	ret = uintptr(C.result_get_offset(result.handle, &cOffset))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Offset = time.Nanosecond * time.Duration(100*cOffset)
	/* Properties */
	var propBagHandle C.SPXHANDLE
	ret = uintptr(C.result_get_property_bag(uintptr2handle(handle), &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return result, nil
}

// SpeechRecognitionOutcome is a wrapper type to be returned by operations returning SpeechRecognitionResult and error
type SpeechRecognitionOutcome struct {
	common.OperationOutcome

	// Result is the result of the operation
	Result *SpeechRecognitionResult
}

// Close releases the underlying resources
func (outcome SpeechRecognitionOutcome) Close() {
	if outcome.Result != nil {
		outcome.Result.Close()
	}
}
