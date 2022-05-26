// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	"strings"
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// #include <stdlib.h>
// #include <string.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_factory.h>
// #include <speechapi_c_property_bag.h>
// #include <speechapi_c_speaker_recognition.h>
//
// SPXHR get_profiles_json_proxy(SPXVOICEPROFILECLIENTHANDLE hVoiceProfileClient, int type, char* buffer, size_t* pcch)
// {
//     char* b = NULL;
//     size_t s = 0;
//     SPXHR hr = get_profiles_json(hVoiceProfileClient, type, &b, &s);
//     *pcch = s;
//     if ((buffer != NULL) && (b != NULL))
//     {
//         memcpy(buffer, b, s);
//     }
//     if (b != NULL)
//     {
//	       property_bag_free_string(b);
//     }
//     return hr;
// }
//
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

	Profile *VoiceProfile
}

// Close releases the underlying resources
func (outcome CreateProfileOutcome) Close() {
	if outcome.Profile != nil {
		outcome.Profile.Close()
	}
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
			outcome <- CreateProfileOutcome{Profile: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			newProfile, err := newVoiceProfileFromHandle(handle2uintptr(profileHandle))
			if err != nil {
				outcome <- CreateProfileOutcome{Profile: nil, OperationOutcome: common.OperationOutcome{err}}
			} else {
				outcome <- CreateProfileOutcome{Profile: newProfile, OperationOutcome: common.OperationOutcome{nil}}
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
			result, err := newVoiceProfileResultFromHandle(handle2uintptr(handle))
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
			result, err := newVoiceProfileResultFromHandle(handle2uintptr(handle))
			outcome <- VoiceProfileOutcome{Result: result, OperationOutcome: common.OperationOutcome{err}}
		}
	}()
	return outcome
}

// GetActivationPhrasesAsync returns a result containing a list of activation phrases required for voice profile enrollment.
func (client VoiceProfileClient) GetActivationPhrasesAsync(profileType common.VoiceProfileType, locale string) chan VoiceProfilePhraseOutcome {
	outcome := make(chan VoiceProfilePhraseOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		loc := C.CString(locale)
		defer C.free(unsafe.Pointer(loc))
		ret := uintptr(C.get_activation_phrases(client.handle, loc, (C.int)(profileType), &handle))
		if ret != C.SPX_NOERROR {
			outcome <- VoiceProfilePhraseOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			newResult, err := newVoiceProfilePhraseResultFromHandle(handle2uintptr(handle))
			if err != nil {
				outcome <- VoiceProfilePhraseOutcome{Result: nil, OperationOutcome: common.OperationOutcome{err}}
			} else {
				outcome <- VoiceProfilePhraseOutcome{Result: newResult, OperationOutcome: common.OperationOutcome{nil}}
			}
		}
	}()
	return outcome
}

// EnrollProfileAsync sends audio for voice profile enrollment returns a result detailing enrollment status for the given profile
func (client VoiceProfileClient) EnrollProfileAsync(profile *VoiceProfile, audioConfig *audio.AudioConfig) chan VoiceProfileEnrollmentOutcome {
	outcome := make(chan VoiceProfileEnrollmentOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		profileHandle := profile.GetHandle()
		var audioHandle C.SPXHANDLE
		if audioConfig == nil {
			audioHandle = nil
		} else {
			audioHandle = uintptr2handle(audioConfig.GetHandle())
		}
		ret := uintptr(C.enroll_voice_profile(client.handle, uintptr2handle(profileHandle), audioHandle, &handle))
		if ret != C.SPX_NOERROR {
			outcome <- VoiceProfileEnrollmentOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			newResult, err := newVoiceProfileEnrollmentResultFromHandle(handle2uintptr(handle))
			if err != nil {
				outcome <- VoiceProfileEnrollmentOutcome{Result: nil, OperationOutcome: common.OperationOutcome{err}}
			} else {
				outcome <- VoiceProfileEnrollmentOutcome{Result: newResult, OperationOutcome: common.OperationOutcome{nil}}
			}
		}
	}()
	return outcome
}

// RetrieveEnrollmentResultAsync returns a result detailing enrollment status for the given profile
func (client VoiceProfileClient) RetrieveEnrollmentResultAsync(profile *VoiceProfile) chan VoiceProfileEnrollmentOutcome {
	outcome := make(chan VoiceProfileEnrollmentOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		id, err := profile.Id()
		if err != nil {
			outcome <- VoiceProfileEnrollmentOutcome{Result: nil, OperationOutcome: common.OperationOutcome{err}}
		}
		cId := C.CString(id)
		defer C.free(unsafe.Pointer(cId))
		profileType, err := profile.Type()
		if err != nil {
			outcome <- VoiceProfileEnrollmentOutcome{Result: nil, OperationOutcome: common.OperationOutcome{err}}
		}
		ret := uintptr(C.retrieve_enrollment_result(client.handle, cId, (C.int)(profileType), &handle))
		if ret != C.SPX_NOERROR {
			outcome <- VoiceProfileEnrollmentOutcome{Result: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(ret)}}
		} else {
			newResult, err := newVoiceProfileEnrollmentResultFromHandle(handle2uintptr(handle))
			if err != nil {
				outcome <- VoiceProfileEnrollmentOutcome{Result: nil, OperationOutcome: common.OperationOutcome{err}}
			} else {
				outcome <- VoiceProfileEnrollmentOutcome{Result: newResult, OperationOutcome: common.OperationOutcome{nil}}
			}
		}
	}()
	return outcome
}

type GetAllProfilesOutcome struct {
	common.OperationOutcome

	Profiles []*VoiceProfile
}

// Close releases the underlying resources
func (outcome GetAllProfilesOutcome) Close() {
	for _, profile := range outcome.Profiles {
		if profile != nil {
			profile.Close()
		}
	}
}

// GetAllProfilesAsync attempts to create a new voice profile on the service.
func (client VoiceProfileClient) GetAllProfilesAsync(profileType common.VoiceProfileType) chan GetAllProfilesOutcome {
	outcome := make(chan GetAllProfilesOutcome)
	go func() {
		var size C.size_t
		ret := uintptr(C.get_profiles_json_proxy(client.handle, (C.int)(profileType), nil, &size))
		if ret != C.SPX_NOERROR {
			outcome <- GetAllProfilesOutcome{Profiles: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))}}
		}
		rawProfileJson := C.malloc(C.sizeof_char * (size))
		defer C.free(unsafe.Pointer(rawProfileJson))
		ret = uintptr(C.get_profiles_json_proxy(client.handle, (C.int)(profileType), (*C.char)(rawProfileJson), &size))
		if ret != C.SPX_NOERROR {
			outcome <- GetAllProfilesOutcome{Profiles: nil, OperationOutcome: common.OperationOutcome{common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))}}
		} else {
			goProfilesJson := C.GoString((*C.char)(rawProfileJson))
			splitProfileIds := strings.Split(goProfilesJson, "|")
			profileList := make([]*VoiceProfile, len(splitProfileIds))
			for index, id := range splitProfileIds {
				profile, err := NewVoiceProfileFromIdAndType(id, profileType)
				if err != nil {
					outcome <- GetAllProfilesOutcome{Profiles: nil, OperationOutcome: common.OperationOutcome{err}}
				}
				profileList[index] = profile
			}
			outcome <- GetAllProfilesOutcome{Profiles: profileList, OperationOutcome: common.OperationOutcome{nil}}
		}
	}()
	return outcome
}
