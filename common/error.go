// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

import (
	"fmt"
)

// #include <speechapi_c_error.h>
import "C"

type CarbonError struct {
	Code    int
	Message string
}

var errorString = map[int]string{
	0x000: "SPX_NOERROR",
	0xfff: "SPXERR_NOT_IMPL",
	0x001: "SPXERR_UNINITIALIZED",
	0x002: "SPXERR_ALREADY_INITIALIZED",
	0x003: "SPXERR_UNHANDLED_EXCEPTION",
	0x004: "SPXERR_NOT_FOUND",
	0x005: "SPXERR_INVALID_ARG",
	0x006: "SPXERR_TIMEOUT",
	0x007: "SPXERR_ALREADY_IN_PROGRESS",
	0x008: "SPXERR_FILE_OPEN_FAILED",
	0x009: "SPXERR_UNEXPECTED_EOF",
	0x00a: "SPXERR_INVALID_HEADER",
	0x00b: "SPXERR_AUDIO_IS_PUMPING",
	0x00c: "SPXERR_UNSUPPORTED_FORMAT",
	0x00d: "SPXERR_ABORT",
	0x00e: "SPXERR_MIC_NOT_AVAILABLE",
	0x00f: "SPXERR_INVALID_STATE",
	0x010: "SPXERR_UUID_CREATE_FAILED",
	0x011: "SPXERR_SETFORMAT_UNEXPECTED_STATE_TRANSITION",
	0x012: "SPXERR_PROCESS_AUDIO_INVALID_STATE",
	0x013: "SPXERR_START_RECOGNIZING_INVALID_STATE_TRANSITION",
	0x014: "SPXERR_UNEXPECTED_CREATE_OBJECT_FAILURE",
	0x015: "SPXERR_MIC_ERROR",
	0x016: "SPXERR_NO_AUDIO_INPUT",
	0x017: "SPXERR_UNEXPECTED_USP_SITE_FAILURE",
	0x018: "SPXERR_UNEXPECTED_UNIDEC_SITE_FAILURE",
	0x019: "SPXERR_BUFFER_TOO_SMALL",
	0x01A: "SPXERR_OUT_OF_MEMORY",
	0x01B: "SPXERR_RUNTIME_ERROR",
	0x01C: "SPXERR_INVALID_URL",
	0x01D: "SPXERR_INVALID_REGION",
	0x01E: "SPXERR_SWITCH_MODE_NOT_ALLOWED",
	0x01F: "SPXERR_CHANGE_CONNECTION_STATUS_NOT_ALLOWED",
	0x020: "SPXERR_EXPLICIT_CONNECTION_NOT_SUPPORTED_BY_RECOGNIZER",
	0x021: "SPXERR_INVALID_HANDLE",
	0x022: "SPXERR_INVALID_RECOGNIZER",
	0x023: "SPXERR_OUT_OF_RANGE",
	0x024: "SPXERR_EXTENSION_LIBRARY_NOT_FOUND",
	0x025: "SPXERR_UNEXPECTED_TTS_ENGINE_SITE_FAILURE",
	0x026: "SPXERR_UNEXPECTED_AUDIO_OUTPUT_FAILURE",
	0x027: "SPXERR_GSTREAMER_INTERNAL_ERROR",
	0x028: "SPXERR_CONTAINER_FORMAT_NOT_SUPPORTED_ERROR",
	0x029: "SPXERR_GSTREAMER_NOT_FOUND_ERROR",
	0x02A: "SPXERR_INVALID_LANGUAGE",
	0x02B: "SPXERR_UNSUPPORTED_API_ERROR",
	0x02C: "SPXERR_RINGBUFFER_DATA_UNAVAILABLE",
	0x030: "SPXERR_UNEXPECTED_CONVERSATION_SITE_FAILURE",
	0x031: "SPXERR_UNEXPECTED_CONVERSATION_TRANSLATOR_SITE_FAILURE",
	0x032: "SPXERR_CANCELED",
}

func NewCarbonError(errorHandle uintptr) CarbonError {
	var carbonError CarbonError
	carbonError.Code = getErrorCode(SPXHandle(errorHandle))
	carbonError.Message = getErrorMessage(SPXHandle(errorHandle))
	// When the message is empty, construct the error message using the errorHandle value directly.
	if carbonError.Message == "" {
		codeAsHexString := fmt.Sprintf("0x%0x", carbonError.Code)
		carbonError.Message = "Exception with an error code: " + codeAsHexString + " (" + errorString[carbonError.Code] + ")"
	}
	return carbonError
}

func (e CarbonError) Error() string {
	return e.Message
}

func getErrorCode(errorHandle SPXHandle) int {
	ret := int(C.error_get_error_code(uintptr2handle(errorHandle)))
	// A 0 means there was no corresponding event stored.
	// So this must be a SPX_* error and not a stored exception.
	// Return the HR as the error.
	if ret == 0 {
		return int(errorHandle)
	}
	return ret
}

func getErrorMessage(errorHandle SPXHandle) string {
	message := ""
	ret := C.error_get_message(uintptr2handle(errorHandle))
	if ret != nil {
		message = C.GoString(ret)
	}
	return message
}
