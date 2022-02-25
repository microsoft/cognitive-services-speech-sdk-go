// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_speaker_recognition.h>
import "C"

// VoiceProfile is the class that defines voice profiles used in speaker recognition scenarios.
type VoiceProfile struct {
	handle     C.SPXHANDLE
}

// newVoiceProfileFromHandle creates a VoiceProfile instance from a valid handle. This is for internal use only.
func newVoiceProfileFromHandle(handle common.SPXHandle) (*VoiceProfile, error) {
	profile := new(VoiceProfile)
	profile.handle = uintptr2handle(handle)
	return profile, nil
}

// NewVoiceProfileFromIdAndType creates an instance of the voice profile with specified id and type.
func NewVoiceProfileFromIdAndType(id string, profileType common.VoiceProfileType) (*VoiceProfile, error) {
	var handle C.SPXHANDLE
	profileId := C.CString(id)
	defer C.free(unsafe.Pointer(profileId))
	ret := uintptr(C.create_voice_profile_from_id_and_type(&handle, profileId, (C.int)(profileType)))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	
	return newVoiceProfileFromHandle(handle2uintptr(handle))
}

// Return the id of the given voice profile 
func (profile *VoiceProfile) Id() (string, error) {
	var sz C.uint32_t
	ret := uintptr(C.voice_profile_get_id(profile.handle, nil, &sz))
	if ret != C.SPX_NOERROR {
		return "", common.NewCarbonError(ret)
	}
	buffer := C.malloc(C.sizeof_char * (C.size_t)(sz))
	defer C.free(unsafe.Pointer(buffer))
	ret = uintptr(C.voice_profile_get_id(profile.handle, (*C.char)(buffer), &sz))
	if ret != C.SPX_NOERROR {
		return "", common.NewCarbonError(ret)
	}
	id := C.GoString((*C.char)(buffer))
	return id, nil
}

// Return the type of the given voice profile 
func (profile *VoiceProfile) Type() (common.VoiceProfileType, error) {
	var profileType C.int
	ret := uintptr(C.voice_profile_get_type(profile.handle, &profileType))
	if ret != C.SPX_NOERROR {
		return common.VoiceProfileType(1), common.NewCarbonError(ret)
	}
	return common.VoiceProfileType(profileType), nil
}

// Close disposes the associated resources.
func (profile *VoiceProfile) Close() {
	C.voice_profile_release_handle(profile.handle)
}

func (profile *VoiceProfile) GetHandle() common.SPXHandle {
	return handle2uintptr(profile.handle)
}
