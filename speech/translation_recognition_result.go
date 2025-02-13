// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_recognizer.h>
// #include <speechapi_c_result.h>
// #include <speechapi_c_translation_result.h>
import "C"

// TranslationRecognitionResult represents the result of a translation recognition.
type TranslationRecognitionResult struct {
	SpeechRecognitionResult
	translations map[string]string
}

// NewTranslationRecognitionResultFromHandle creates a TranslationRecognitionResult from a handle.
func NewTranslationRecognitionResultFromHandle(handle common.SPXHandle) (*TranslationRecognitionResult, error) {
	result := new(TranslationRecognitionResult)
	result.translations = make(map[string]string)

	// Get base recognition result
	baseResult, err := NewSpeechRecognitionResultFromHandle(handle)
	if err != nil {
		return nil, err
	}
	result.SpeechRecognitionResult = *baseResult

	// Get translation count
	var count C.size_t
	ret := uintptr(C.translation_text_result_get_translation_count(uintptr2handle(handle), &count))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	// Get translations
	for i := C.size_t(0); i < count; i++ {
		var languageSize, textSize C.size_t
		ret = uintptr(C.translation_text_result_get_translation(uintptr2handle(handle), i, nil, nil, &languageSize, &textSize))
		if ret != C.SPX_NOERROR {
			return nil, common.NewCarbonError(ret)
		}

		language := make([]byte, languageSize)
		text := make([]byte, textSize)
		ret = uintptr(C.translation_text_result_get_translation(uintptr2handle(handle), i,
			(*C.char)(unsafe.Pointer(&language[0])),
			(*C.char)(unsafe.Pointer(&text[0])),
			&languageSize, &textSize))
		if ret != C.SPX_NOERROR {
			return nil, common.NewCarbonError(ret)
		}

		result.translations[string(language[:languageSize-1])] = string(text[:textSize-1])
	}

	return result, nil
}

// GetTranslations returns all available translations.
func (result TranslationRecognitionResult) GetTranslations() map[string]string {
	return result.translations
}

// GetTranslation returns the translation for the specified language.
func (result TranslationRecognitionResult) GetTranslation(language string) string {
	return result.translations[language]
}

// TranslationSynthesisResult represents the voice output of the translated text.
type TranslationSynthesisResult struct {
	Reason    common.ResultReason
	audioData []byte
}

// NewTranslationSynthesisResultFromHandle creates a TranslationSynthesisResult from a handle.
func NewTranslationSynthesisResultFromHandle(handle common.SPXHandle) (*TranslationSynthesisResult, error) {
	result := new(TranslationSynthesisResult)

	var reason C.Result_Reason
	ret := uintptr(C.result_get_reason(uintptr2handle(handle), &reason))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result.Reason = common.ResultReason(reason)

	var size C.size_t
	ret = uintptr(C.translation_synthesis_result_get_audio_data(uintptr2handle(handle), nil, &size))
	if ret == uintptr(C.SPXERR_BUFFER_TOO_SMALL) {
		result.audioData = make([]byte, size)
		ret = uintptr(C.translation_synthesis_result_get_audio_data(uintptr2handle(handle),
			(*C.uint8_t)(unsafe.Pointer(&result.audioData[0])), &size))
	}
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	return result, nil
}

// GetAudioData returns the voice output of the translated text.
func (result TranslationSynthesisResult) GetAudioData() []byte {
	return result.audioData
}

// TranslationRecognitionEventArgs represents the event arguments for a translation recognition event.
type TranslationRecognitionEventArgs struct {
	RecognitionEventArgs
	Result *TranslationRecognitionResult
}

// NewTranslationRecognitionEventArgsFromHandle creates a TranslationRecognitionEventArgs from a handle.
func NewTranslationRecognitionEventArgsFromHandle(handle common.SPXHandle) (*TranslationRecognitionEventArgs, error) {
	base, err := NewRecognitionEventArgsFromHandle(handle)
	if err != nil {
		return nil, err
	}

	event := new(TranslationRecognitionEventArgs)
	event.RecognitionEventArgs = *base
	event.handle = uintptr2handle(handle)

	var resultHandle C.SPXRESULTHANDLE
	ret := uintptr(C.recognizer_recognition_event_get_result(event.handle, &resultHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	result, err := NewTranslationRecognitionResultFromHandle(handle2uintptr(resultHandle))
	if err != nil {
		return nil, err
	}

	event.Result = result
	return event, nil
}

// TranslationRecognitionCanceledEventArgs represents the event arguments for a translation recognition canceled event.
type TranslationRecognitionCanceledEventArgs struct {
	TranslationRecognitionEventArgs
	ErrorDetails string
	Reason       common.CancellationReason
	ErrorCode    common.CancellationErrorCode
}

// NewTranslationRecognitionCanceledEventArgsFromHandle creates a TranslationRecognitionCanceledEventArgs from a handle.
func NewTranslationRecognitionCanceledEventArgsFromHandle(handle common.SPXHandle) (*TranslationRecognitionCanceledEventArgs, error) {
	var reason C.Result_CancellationReason
	var errorCode C.Result_CancellationErrorCode

	baseArgs, err := NewTranslationRecognitionEventArgsFromHandle(handle)
	if err != nil {
		return nil, err
	}
	event := new(TranslationRecognitionCanceledEventArgs)
	event.TranslationRecognitionEventArgs = *baseArgs

	ret := uintptr(C.result_get_reason_canceled(event.Result.handle, &reason))
	if ret != C.SPX_NOERROR {
		event.Close()
		return nil, common.NewCarbonError(ret)
	}

	/* ErrorCode */
	ret = uintptr(C.result_get_canceled_error_code(event.Result.handle, &errorCode))
	if ret != C.SPX_NOERROR {
		event.Close()
		return nil, common.NewCarbonError(ret)
	}

	event.ErrorDetails = event.Result.Properties.GetProperty(common.SpeechServiceResponseJSONErrorDetails, "")
	event.ErrorCode = (common.CancellationErrorCode)(errorCode)
	event.Reason = (common.CancellationReason)(reason)

	return event, nil
}

// TranslationSynthesisEventArgs represents the event arguments for a translation synthesis event.
type TranslationSynthesisEventArgs struct {
	SessionEventArgs
	Result *TranslationSynthesisResult
}

// NewTranslationSynthesisEventArgsFromHandle creates a TranslationSynthesisEventArgs from a handle.
func NewTranslationSynthesisEventArgsFromHandle(handle common.SPXHandle) (*TranslationSynthesisEventArgs, error) {
	var resultHandle C.SPXRESULTHANDLE
	ret := uintptr(C.recognizer_recognition_event_get_result(uintptr2handle(handle), &resultHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	result, err := NewTranslationSynthesisResultFromHandle(handle2uintptr(resultHandle))
	if err != nil {
		return nil, err
	}

	return &TranslationSynthesisEventArgs{Result: result}, nil
}

// Event handler types
type TranslationRecognitionEventHandler func(event TranslationRecognitionEventArgs)
type TranslationRecognitionCanceledEventHandler func(event TranslationRecognitionCanceledEventArgs)
type TranslationSynthesisEventHandler func(event TranslationSynthesisEventArgs)

// TranslationRecognitionOutcome represents the outcome of a translation recognition operation.
type TranslationRecognitionOutcome struct {
	Result *TranslationRecognitionResult
	common.OperationOutcome
}
