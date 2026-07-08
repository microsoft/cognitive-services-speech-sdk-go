// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"strings"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <speechapi_c_speech_recognition_model.h>
import "C"

// SpeechRecognitionModelInfo contains information about a speech recognition model that is
// available for embedded (offline) speech recognition.
type SpeechRecognitionModelInfo struct {
	handle C.SPXHANDLE
}

// splitLocales splits the pipe-delimited locales string returned by the native layer.
func splitLocales(locales string) []string {
	if locales == "" {
		return []string{}
	}
	return strings.Split(locales, "|")
}

func newSpeechRecognitionModelFromHandle(handle C.SPXHANDLE) *SpeechRecognitionModelInfo {
	model := new(SpeechRecognitionModelInfo)
	model.handle = handle
	return model
}

// GetHandle gets the handle to the resource (for internal use).
func (model SpeechRecognitionModelInfo) GetHandle() common.SPXHandle {
	return handle2uintptr(model.handle)
}

// Name returns the model name.
func (model SpeechRecognitionModelInfo) Name() string {
	return C.GoString(C.speech_recognition_model_get_name(model.handle))
}

// Locales returns the locales supported by the model.
func (model SpeechRecognitionModelInfo) Locales() []string {
	return splitLocales(C.GoString(C.speech_recognition_model_get_locales(model.handle)))
}

// Path returns the model path.
func (model SpeechRecognitionModelInfo) Path() string {
	return C.GoString(C.speech_recognition_model_get_path(model.handle))
}

// Version returns the model version.
func (model SpeechRecognitionModelInfo) Version() string {
	return C.GoString(C.speech_recognition_model_get_version(model.handle))
}

// Close disposes the associated resources.
func (model SpeechRecognitionModelInfo) Close() {
	C.speech_recognition_model_handle_release(model.handle)
}
