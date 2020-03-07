//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_recognizer.h>
// #include <speechapi_c_factory.h>
import "C"

type SpeechRecognizer struct {
	Properties common.PropertyCollection
	handle     C.SPXHANDLE
}

func newSpeechRecognizerFromHandle(handle C.SPXHANDLE) (*SpeechRecognizer, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.recognizer_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	recognizer := new(SpeechRecognizer)
	recognizer.handle = handle
	recognizer.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return recognizer, nil
}

/// Create a speech recognizer from a speech config and audio config.
func NewSpeechRecognizerFromConfig(config *SpeechConfig, audioConfig *audio.AudioConfig) (*SpeechRecognizer, error) {
	var handle C.SPXHANDLE
	if config == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	configHandle := config.getHandle()
	var audioHandle C.SPXHANDLE
	if audioConfig == nil {
		audioHandle = nil
	} else {
		audioHandle = uintptr2handle(audioConfig.GetHandle())
	}
	ret := uintptr(C.recognizer_create_speech_recognizer_from_config(&handle, configHandle, audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSpeechRecognizerFromHandle(handle)
}
