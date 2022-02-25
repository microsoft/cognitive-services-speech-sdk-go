// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// #include <stdlib.h>
// #include <speechapi_c_factory.h>
// #include <speechapi_c_speaker_recognition.h>
//
import "C"

// SpeakerRecognizer is the class for speaker recognizers.
type SpeakerRecognizer struct {
	Properties                 *common.PropertyCollection
	handle                     C.SPXHANDLE
}

func newSpeakerRecognizerFromHandle(handle C.SPXHANDLE) (*SpeakerRecognizer, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.speaker_recognizer_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	recognizer := new(SpeakerRecognizer)
	recognizer.handle = handle
	recognizer.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return recognizer, nil
}

// NewSpeakerRecognizerFromConfig creates a speaker recognizer from a speech config and audio config.
func NewSpeakerRecognizerFromConfig(config *speech.SpeechConfig, audioConfig *audio.AudioConfig) (*SpeakerRecognizer, error) {
	var handle C.SPXHANDLE
	if config == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	configHandle := config.GetHandle()
	var audioHandle C.SPXHANDLE
	if audioConfig == nil {
		audioHandle = nil
	} else {
		audioHandle = uintptr2handle(audioConfig.GetHandle())
	}
	ret := uintptr(C.recognizer_create_speaker_recognizer_from_config(&handle, uintptr2handle(configHandle), audioHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSpeakerRecognizerFromHandle(handle)
}

// VerifyOnceAsync starts speaker verification, and returns a score indicates whether the profile in the model is verified or not
func (recognizer SpeakerRecognizer) VerifyOnceAsync(model *SpeakerVerificationModel) chan SpeakerRecognitionOutcome {
	outcome := make(chan SpeakerRecognitionOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		modelHandle := uintptr2handle(model.GetHandle())
		ret := uintptr(C.speaker_recognizer_verify(recognizer.handle, modelHandle, &handle))
		if ret != C.SPX_NOERROR {
			outcome <- SpeakerRecognitionOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewSpeakerRecognitionResultFromHandle(handle2uintptr(handle))
			outcome <- SpeakerRecognitionOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// IdentifyOnceAsync starts speaker verification, and returns a score indicates whether the profile in the model is verified or not
func (recognizer SpeakerRecognizer) IdentifyOnceAsync(model *SpeakerIdentificationModel) chan SpeakerRecognitionOutcome {
	outcome := make(chan SpeakerRecognitionOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		modelHandle := uintptr2handle(model.GetHandle())
		ret := uintptr(C.speaker_recognizer_identify(recognizer.handle, modelHandle, &handle))
		if ret != C.SPX_NOERROR {
			outcome <- SpeakerRecognitionOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewSpeakerRecognitionResultFromHandle(handle2uintptr(handle))
			outcome <- SpeakerRecognitionOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// SetAuthorizationToken sets the authorization token that will be used for connecting to the service.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token
// expires, the caller needs to refresh it by calling this setter with a new valid token.
// Otherwise, the recognizer will encounter errors during recognition.
func (recognizer SpeakerRecognizer) SetAuthorizationToken(token string) error {
	return recognizer.Properties.SetProperty(common.SpeechServiceAuthorizationToken, token)
}

// AuthorizationToken is the authorization token.
func (recognizer SpeakerRecognizer) AuthorizationToken() string {
	return recognizer.Properties.GetProperty(common.SpeechServiceAuthorizationToken, "")
}

// Close disposes the associated resources.
func (recognizer SpeakerRecognizer) Close() {
	recognizer.Properties.Close()
	C.speaker_recognizer_release_handle(recognizer.handle)
}
