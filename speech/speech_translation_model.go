// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <speechapi_c_speech_translation_model.h>
import "C"

// SpeechTranslationModelInfo contains information about a speech translation model that is
// available for embedded (offline) speech translation.
type SpeechTranslationModelInfo struct {
	handle C.SPXHANDLE
}

func newSpeechTranslationModelFromHandle(handle C.SPXHANDLE) *SpeechTranslationModelInfo {
	model := new(SpeechTranslationModelInfo)
	model.handle = handle
	return model
}

// GetHandle gets the handle to the resource (for internal use).
func (model SpeechTranslationModelInfo) GetHandle() common.SPXHandle {
	return handle2uintptr(model.handle)
}

// Name returns the model name.
func (model SpeechTranslationModelInfo) Name() string {
	return C.GoString(C.speech_translation_model_get_name(model.handle))
}

// SourceLanguages returns the source languages supported by the model.
func (model SpeechTranslationModelInfo) SourceLanguages() []string {
	return splitLocales(C.GoString(C.speech_translation_model_get_source_languages(model.handle)))
}

// TargetLanguages returns the target languages supported by the model.
func (model SpeechTranslationModelInfo) TargetLanguages() []string {
	return splitLocales(C.GoString(C.speech_translation_model_get_target_languages(model.handle)))
}

// Path returns the model path.
func (model SpeechTranslationModelInfo) Path() string {
	return C.GoString(C.speech_translation_model_get_path(model.handle))
}

// Version returns the model version.
func (model SpeechTranslationModelInfo) Version() string {
	return C.GoString(C.speech_translation_model_get_version(model.handle))
}

// Close disposes the associated resources.
func (model SpeechTranslationModelInfo) Close() {
	C.speech_translation_model_handle_release(model.handle)
}
