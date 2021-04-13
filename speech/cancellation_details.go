// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_result.h>
//
import "C"

// CancellationDetails contains detailed information about why a result was canceled.
// Added in version 1.17.0
type CancellationDetails struct {
	Reason       common.CancellationReason
	ErrorCode    common.CancellationErrorCode
	ErrorDetails string
}

// NewCancellationDetailsFromSpeechSynthesisResult creates the object from the speech synthesis result.
func NewCancellationDetailsFromSpeechSynthesisResult(result *SpeechSynthesisResult) (*CancellationDetails, error) {
	cancellationDetails := new(CancellationDetails)
	/* Reason */
	var cReason C.Result_CancellationReason
	ret := uintptr(C.synth_result_get_reason_canceled(result.handle, &cReason))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	cancellationDetails.Reason = (common.CancellationReason)(cReason)
	/* ErrorCode */
	var cCode C.Result_CancellationErrorCode
	ret = uintptr(C.synth_result_get_canceled_error_code(result.handle, &cCode))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	cancellationDetails.ErrorCode = (common.CancellationErrorCode)(cCode)
	cancellationDetails.ErrorDetails = result.Properties.GetProperty(common.CancellationDetailsReasonDetailedText, "")
	return cancellationDetails, nil
}
