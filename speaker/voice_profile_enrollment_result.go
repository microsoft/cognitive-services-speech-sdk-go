// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	"math/big"
	"unsafe"
	"strconv"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_result.h>
// #include <speechapi_c_recognizer.h>
//
import "C"

// VoiceProfileEnrollmentResult contains information about result from voice profile operations.
type VoiceProfileEnrollmentResult struct {
	handle C.SPXHANDLE

	// ResultID specifies the result identifier.
	ResultID string

	// ProfileID specifies the profile ID of the profile being enrolled.
	ProfileID string

	// EnrollmentsCount specifies the number of successful enrollments for the profile
	EnrollmentsCount int

	// RemainingEnrollmentsCount specifies the number of successful enrollments remaining until profile is enrolled
	RemainingEnrollmentsCount int

	// EnrollmentsLength specifies in hundreds of nanoseconds the audio length registered enrolling the profile
	EnrollmentsLength big.Int

	// RemainingEnrollmentsLength specifies the amount of pure speech (which is the amount of audio after removing silence and non-speech segments) needed to complete profile enrollment in hundred nanoseconds.
	RemainingEnrollmentsLength big.Int

	// CreatedTime specifies the created time of the voice profile.
	CreatedTime string

	// LastUpdatedDateTime specifies the last updated time of the voice profile.
	LastUpdatedTime string

	// Reason specifies status of speech synthesis result.
	Reason common.ResultReason

	// ErrorDetails presents error details.
	ErrorDetails string

	// Collection of additional properties.
	Properties *common.PropertyCollection
}

// Close releases the underlying resources
func (result VoiceProfileEnrollmentResult) Close() {
	result.Properties.Close()
	C.recognizer_result_handle_release(result.handle)
}

// NewVoiceProfileEnrollmentResultFromHandle creates a VoiceProfileEnrollmentResult from a handle (for internal use)
func NewVoiceProfileEnrollmentResultFromHandle (handle common.SPXHandle) (*VoiceProfileEnrollmentResult, error) {
	result := new(VoiceProfileEnrollmentResult)
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

	/* ProfileID */
	result.ProfileID = result.Properties.GetPropertyByString("enrollment.profileId", "")
	
	/* EnrollmentsCount */
	value := result.Properties.GetPropertyByString("enrollment.enrollmentsCount", "")
	if value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			result.EnrollmentsCount = intVal
		}
	}
	
	/* RemainingEnrollmentsCount */
	value = result.Properties.GetPropertyByString("enrollment.remainingEnrollmentsCount", "")
	if value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			result.RemainingEnrollmentsCount = intVal
		}
	}

	/* EnrollmentsLength */
	value = result.Properties.GetPropertyByString("enrollment.enrollmentsLengthInSec", "")
	if value != "" {
		bigIntVal := new(big.Int)
		if bigIntVal, ok := bigIntVal.SetString(value, 10); ok {
			result.EnrollmentsLength = *bigIntVal
		}
	}
	
	/* RemainingEnrollmentsLength */
	value = result.Properties.GetPropertyByString("enrollment.remainingEnrollmentsLengthInSec", "")
	if value != "" {
		bigIntVal := new(big.Int)
		if bigIntVal, ok := bigIntVal.SetString(value, 10); ok {
			result.RemainingEnrollmentsLength = *bigIntVal
		}
	}

	/* CreatedTime */
	result.CreatedTime = result.Properties.GetPropertyByString("enrollment.createdDateTime", "")

	/* LastUpdatedTime */
	result.LastUpdatedTime = result.Properties.GetPropertyByString("enrollment.lastUpdatedDateTime", "")

	return result, nil
}

// VoiceProfileEnrollmentOutcome is a wrapper type to be returned by operations returning VoiceProfileEnrollmentResult and error
type VoiceProfileEnrollmentOutcome struct {
	common.OperationOutcome

	// Result is the result of the operation
	Result *VoiceProfileEnrollmentResult
}

// Close releases the underlying resources
func (outcome VoiceProfileEnrollmentOutcome) Close() {
	if outcome.Result != nil {
		outcome.Result.Close()
	}
}
