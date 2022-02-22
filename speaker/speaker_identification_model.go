// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_speaker_recognition.h>
import "C"

// SpeakerIdentificationModel is the class that defines a identification model to be used in speaker identification scenarios.
type SpeakerIdentificationModel struct {
	handle     C.SPXHANDLE
}

// NewSpeakerIdentificationModelFromHandle creates a SpeakerIdentificationModel instance from a valid handle. This is for internal use only.
func NewSpeakerIdentificationModelFromHandle(handle common.SPXHandle) (*SpeakerIdentificationModel, error) {
	model := new(SpeakerIdentificationModel)
	model.handle = uintptr2handle(handle)
	return model, nil
}

// NewSpeakerIdentificationModelFromProfile creates an instance of the identification model using the given voice profiles.
func NewSpeakerIdentificationModelFromProfiles(profiles []*VoiceProfile) (*SpeakerIdentificationModel, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.speaker_identification_model_create(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}

	for _, profile := range profiles {
		profileHandle := profile.GetHandle()
		ret := uintptr(C.speaker_identification_model_add_profile(handle, uintptr2handle(profileHandle)))
		if ret != C.SPX_NOERROR {
			C.speaker_identification_model_release_handle(handle)
			return nil, common.NewCarbonError(ret)
		}
	}
	
	return NewSpeakerIdentificationModelFromHandle(handle2uintptr(handle))
}

// Close disposes the associated resources.
func (model *SpeakerIdentificationModel) Close() {
	C.speaker_identification_model_release_handle(model.handle)
}

func (model *SpeakerIdentificationModel) GetHandle() common.SPXHandle {
	return handle2uintptr(model.handle)
}
