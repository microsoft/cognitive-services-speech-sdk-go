// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_keyword_recognition_model.h>
import "C"
import "unsafe"

// KeywordRecognitionModel represents the keyword recognition model used with StartKeywordRecognitionAsync methods.
type KeywordRecognitionModel struct {
	handle C.SPXHANDLE
}

// Close disposes the associated resources.
func (model KeywordRecognitionModel) Close() {
	C.keyword_recognition_model_handle_release(model.handle)
}

// GetHandle gets the handle to the resource (for internal use)
func (model KeywordRecognitionModel) GetHandle() common.SPXHandle {
	return handle2uintptr(model.handle)
}

/// NewKeywordRecognitionModelFromFile creates a keyword recognition model using the specified file.
func NewKeywordRecognitionModelFromFile(filename string) (*KeywordRecognitionModel, error) {
	var handle C.SPXHANDLE
	f := C.CString(filename)
	defer C.free(unsafe.Pointer(f))
	ret := uintptr(C.keyword_recognition_model_create_from_file(f, &handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	model := new(KeywordRecognitionModel)
	model.handle = handle
	return model, nil
}
