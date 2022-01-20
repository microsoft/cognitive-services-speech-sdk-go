// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	"unsafe"

	// "github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_factory.h>
// #include <speechapi_c_speaker_recognition.h>
import "C"

// VoiceProfileClient connects to a speaker recognition backend.
type VoiceProfileClient struct {
	Properties *common.PropertyCollection
	handle     C.SPXHANDLE
}

func newVoiceProfileClientFromHandle(handle C.SPXHANDLE) (*VoiceProfileClient, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.voice_profile_client_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	client := new(VoiceProfileClient)
	client.handle = handle
	client.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return client, nil
}

// NewVoiceProfileClientFromConfig creates a voice profile service client from a speech config.
// Users should use this function to create a voice profile client.
func NewVoiceProfileClientFromConfig(config *speech.SpeechConfig) (*VoiceProfileClient, error) {
	var handle C.SPXHANDLE
	if config == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	configHandle := config.GetHandle()
	ret := uintptr(C.create_voice_profile_client_from_config(&handle, (uintptr2handle(configHandle))))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newVoiceProfileClientFromHandle(handle)
}

// Close performs cleanup of resources.
func (client VoiceProfileClient) Close() {
	client.Properties.Close()
	C.voice_profile_client_release_handle(client.handle)
}

type CreateProfileOutcome struct {
	common.OperationOutcome

	profile *VoiceProfile
}

// CreateProfileAsync attempts to create a new voice profile on the service.
func (client VoiceProfileClient) CreateProfileAsync(profileType common.VoiceProfileType, locale string) chan CreateProfileOutcome {
	outcome := make(chan CreateProfileOutcome)
	go func() {
		var profileHandle C.SPXHANDLE
		loc := C.CString(locale)
		defer C.free(unsafe.Pointer(loc))
		ret := uintptr(C.create_voice_profile(client.handle, (C.int)(profileType), loc, &profileHandle))
		if ret != C.SPX_NOERROR {
			outcome <- CreateProfileOutcome{profile: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			newProfile, err := NewVoiceProfileFromHandle(handle2uintptr(profileHandle))
			if err != nil {
				outcome <- CreateProfileOutcome{profile: nil, OperationOutcome: common.OperationOutcome{err}}
			} else {
				outcome <- CreateProfileOutcome{profile: newProfile, OperationOutcome: common.OperationOutcome{nil}}
			}
		}
	}()
	return outcome
}

// DeleteProfileAsync sends a profile delete request to the service.
func (client VoiceProfileClient) DeleteProfileAsync(profile *VoiceProfile) <-chan VoiceProfileOutcome {
	outcome := make(chan VoiceProfileOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		profileHandle := profile.GetHandle()
		ret := uintptr(C.delete_voice_profile(client.handle, uintptr2handle(profileHandle), &handle))
		if ret != C.SPX_NOERROR {
			outcome <- VoiceProfileOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewVoiceProfileResultFromHandle(handle2uintptr(handle))
			outcome <- VoiceProfileOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// ResetProfileAsync sends a profile reset request to the service.
func (client VoiceProfileClient) ResetProfileAsync(profile *VoiceProfile) <-chan VoiceProfileOutcome {
	outcome := make(chan VoiceProfileOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		profileHandle := profile.GetHandle()
		ret := uintptr(C.reset_voice_profile(client.handle, uintptr2handle(profileHandle), &handle))
		if ret != C.SPX_NOERROR {
			outcome <- VoiceProfileOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			result, err := NewVoiceProfileResultFromHandle(handle2uintptr(handle))
			outcome <- VoiceProfileOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

