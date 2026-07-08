// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_embedded_speech_config.h>
import "C"

// EmbeddedSpeechConfig defines configurations for embedded (offline) speech recognition, synthesis and
// translation. An embedded speech config wraps a regular SpeechConfig, so the underlying SpeechConfig can
// be passed directly to recognizer and synthesizer factory functions (for example
// NewSpeechRecognizerFromConfig). For translation, use NewTranslationRecognizerFromEmbeddedConfig.
//
// Note: Embedded speech recognition, synthesis and translation require licensed models installed on the
// device and the corresponding native runtime extensions. This is a Limited Access feature.
type EmbeddedSpeechConfig struct {
	*SpeechConfig
}

// GetSpeechConfig returns the underlying SpeechConfig. Use it with the existing recognizer and synthesizer
// factory functions, for example NewSpeechRecognizerFromConfig or NewSpeechSynthesizerFromConfig.
func (config *EmbeddedSpeechConfig) GetSpeechConfig() *SpeechConfig {
	return config.SpeechConfig
}

func newEmbeddedSpeechConfigFromHandle(handle C.SPXHANDLE) (*EmbeddedSpeechConfig, error) {
	speechConfig, err := NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}
	config := new(EmbeddedSpeechConfig)
	config.SpeechConfig = speechConfig
	return config, nil
}

// NewEmbeddedSpeechConfigFromPath creates an instance of the embedded speech config with a path to offline
// speech recognition and/or synthesis models. The path can point to a single model directory or to a root
// directory that contains multiple models. This method can be called multiple times (through FromPaths) to
// add several model locations.
func NewEmbeddedSpeechConfigFromPath(path string) (*EmbeddedSpeechConfig, error) {
	return NewEmbeddedSpeechConfigFromPaths([]string{path})
}

// NewEmbeddedSpeechConfigFromPaths creates an instance of the embedded speech config with paths to offline
// speech recognition and/or synthesis models. Each path can point to a single model directory or to a root
// directory that contains multiple models.
func NewEmbeddedSpeechConfigFromPaths(paths []string) (*EmbeddedSpeechConfig, error) {
	if len(paths) == 0 {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	var handle C.SPXHANDLE
	ret := uintptr(C.embedded_speech_config_create(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	for _, path := range paths {
		p := C.CString(path)
		ret = uintptr(C.embedded_speech_config_add_path(handle, p))
		C.free(unsafe.Pointer(p))
		if ret != C.SPX_NOERROR {
			return nil, common.NewCarbonError(ret)
		}
	}
	return newEmbeddedSpeechConfigFromHandle(handle)
}

// GetSpeechRecognitionModels returns the list of embedded speech recognition models available in the
// configured model paths. The caller is responsible for calling Close on each returned model.
func (config *EmbeddedSpeechConfig) GetSpeechRecognitionModels() ([]*SpeechRecognitionModelInfo, error) {
	var numModels C.uint32_t
	ret := uintptr(C.embedded_speech_config_get_num_speech_reco_models(config.getHandle(), &numModels))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	models := make([]*SpeechRecognitionModelInfo, 0, int(numModels))
	for i := 0; i < int(numModels); i++ {
		var modelHandle C.SPXHANDLE
		ret = uintptr(C.embedded_speech_config_get_speech_reco_model(config.getHandle(), C.uint32_t(i), &modelHandle))
		if ret != C.SPX_NOERROR {
			for _, m := range models {
				m.Close()
			}
			return nil, common.NewCarbonError(ret)
		}
		models = append(models, newSpeechRecognitionModelFromHandle(modelHandle))
	}
	return models, nil
}

// SetSpeechRecognitionModel sets the model for embedded speech recognition.
// The name is the model name (see GetSpeechRecognitionModels) and license is the license text for the model.
func (config *EmbeddedSpeechConfig) SetSpeechRecognitionModel(name string, license string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	l := C.CString(license)
	defer C.free(unsafe.Pointer(l))
	ret := uintptr(C.embedded_speech_config_set_speech_recognition_model(config.getHandle(), n, l))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// GetSpeechRecognitionModelName returns the model name for embedded speech recognition.
func (config *EmbeddedSpeechConfig) GetSpeechRecognitionModelName() string {
	return config.GetProperty(common.SpeechServiceConnectionRecoModelName)
}

// SetSpeechSynthesisVoice sets the voice for embedded speech synthesis.
// The name is the voice name and license is the license text for the voice.
func (config *EmbeddedSpeechConfig) SetSpeechSynthesisVoice(name string, license string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	l := C.CString(license)
	defer C.free(unsafe.Pointer(l))
	ret := uintptr(C.embedded_speech_config_set_speech_synthesis_voice(config.getHandle(), n, l))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// GetSpeechSynthesisVoiceName returns the voice name for embedded speech synthesis.
func (config *EmbeddedSpeechConfig) GetSpeechSynthesisVoiceName() string {
	return config.GetProperty(common.SpeechServiceConnectionSynthOfflineVoice)
}

// SetKeywordRecognitionModel sets the model for keyword recognition.
// This is for customer specific models tailored to detect wake words and direct commands.
// The name is the model name and license is the license text for the model.
func (config *EmbeddedSpeechConfig) SetKeywordRecognitionModel(name string, license string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	l := C.CString(license)
	defer C.free(unsafe.Pointer(l))
	ret := uintptr(C.embedded_speech_config_set_keyword_recognition_model(config.getHandle(), n, l))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// GetKeywordRecognitionModelName returns the model name for keyword recognition.
func (config *EmbeddedSpeechConfig) GetKeywordRecognitionModelName() string {
	return config.GetProperty(common.KeywordRecognitionModelName)
}

// GetSpeechTranslationModels returns the list of embedded speech translation models available in the
// configured model paths. The caller is responsible for calling Close on each returned model.
func (config *EmbeddedSpeechConfig) GetSpeechTranslationModels() ([]*SpeechTranslationModelInfo, error) {
	var numModels C.uint32_t
	ret := uintptr(C.embedded_speech_config_get_num_speech_translation_models(config.getHandle(), &numModels))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	models := make([]*SpeechTranslationModelInfo, 0, int(numModels))
	for i := 0; i < int(numModels); i++ {
		var modelHandle C.SPXHANDLE
		ret = uintptr(C.embedded_speech_config_get_speech_translation_model(config.getHandle(), C.uint32_t(i), &modelHandle))
		if ret != C.SPX_NOERROR {
			for _, m := range models {
				m.Close()
			}
			return nil, common.NewCarbonError(ret)
		}
		models = append(models, newSpeechTranslationModelFromHandle(modelHandle))
	}
	return models, nil
}

// SetSpeechTranslationModel sets the model for embedded speech translation.
// The name is the model name (see GetSpeechTranslationModels) and license is the license text for the model.
func (config *EmbeddedSpeechConfig) SetSpeechTranslationModel(name string, license string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	l := C.CString(license)
	defer C.free(unsafe.Pointer(l))
	ret := uintptr(C.embedded_speech_config_set_speech_translation_model(config.getHandle(), n, l))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// GetSpeechTranslationModelName returns the model name for embedded speech translation.
func (config *EmbeddedSpeechConfig) GetSpeechTranslationModelName() string {
	return config.GetProperty(common.SpeechTranslationModelName)
}
