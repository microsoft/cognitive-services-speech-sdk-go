// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	// "fmt"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_speaker_recognition.h>
import "C"

// SpeakerVerificationModel is the class that defines a verification model to be used in speaker verification scenarios.
type SpeakerVerificationModel struct {
	handle     C.SPXHANDLE
}

// NewSpeakerVerificationModelFromHandle creates a SpeakerVerificationModel instance from a valid handle. This is for internal use only.
func NewSpeakerVerificationModelFromHandle(handle common.SPXHandle) (*SpeakerVerificationModel, error) {
	model := new(SpeakerVerificationModel)
	model.handle = uintptr2handle(handle)
	return model, nil
}

// NewSpeakerVerificationModelFromProfile creates an instance of the verification model using the given voice profile.
func NewSpeakerVerificationModelFromProfile(profile VoiceProfile) (*SpeakerVerificationModel, error) {
	var handle C.SPXHANDLE
	profileHandle := profile.GetHandle()
	ret := uintptr(C.speaker_verification_model_create(&handle, uintptr2handle(profileHandle)))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	
	return NewSpeakerVerificationModelFromHandle(handle2uintptr(handle))
}

// Close disposes the associated resources.
func (model *SpeakerVerificationModel) Close() {
	C.speaker_verification_model_release_handle(model.handle)
}

func (model *SpeakerVerificationModel) GetHandle() common.SPXHandle {
	return handle2uintptr(model.handle)
}
