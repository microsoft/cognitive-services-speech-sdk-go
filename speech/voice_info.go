// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"strings"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_result.h>
// #include <speechapi_c_property_bag.h>
//
import "C"

// VoiceInfo contains information about result from voices list of speech synthesizers.
type VoiceInfo struct {
	handle C.SPXHANDLE

	// Name specifies the voice name.
	Name string

	// Locale specifies the locale of the voice
	Locale string

	// ShortName specifies the voice name in short format
	ShortName string

	// LocalName specifies the local name of the voice
	LocalName string

	// Gender specifies the gender of the voice.
	Gender common.SynthesisVoiceGender

	// VoiceType specifies the voice type.
	VoiceType common.SynthesisVoiceType

	// StyleList specifies the styles the voice supports.
	StyleList []string

	// VoicePath specifies the voice path
	VoicePath string

	// Collection of additional properties.
	Properties *common.PropertyCollection
}

// Close releases the underlying resources
func (result VoiceInfo) Close() {
	result.Properties.Close()
	C.voice_info_handle_release(result.handle)
}

// NewVoiceInfoFromHandle creates a VoiceInfo from a handle (for internal use)
func NewVoiceInfoFromHandle(handle common.SPXHandle) (*VoiceInfo, error) {
	voiceInfo := new(VoiceInfo)
	voiceInfo.handle = uintptr2handle(handle)
	/* Name */
	value := C.voice_info_get_name(voiceInfo.handle)
	voiceInfo.Name = C.GoString(value)
	C.property_bag_free_string(value)
	/* Locale */
	value = C.voice_info_get_locale(voiceInfo.handle)
	voiceInfo.Locale = C.GoString(value)
	C.property_bag_free_string(value)
	/* ShortName */
	value = C.voice_info_get_short_name(voiceInfo.handle)
	voiceInfo.ShortName = C.GoString(value)
	C.property_bag_free_string(value)
	/* LocalName */
	value = C.voice_info_get_local_name(voiceInfo.handle)
	voiceInfo.LocalName = C.GoString(value)
	C.property_bag_free_string(value)
	/* StyleList */
	value = C.voice_info_get_style_list(voiceInfo.handle)
	voiceInfo.StyleList = strings.Split(C.GoString(value), "|")
	C.property_bag_free_string(value)
	/* VoiceType */
	var cVoiceType C.Synthesis_VoiceType
	ret := uintptr(C.voice_info_get_voice_type(voiceInfo.handle, &cVoiceType))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	voiceInfo.VoiceType = (common.SynthesisVoiceType)(cVoiceType)
	/* VoicePath */
	value = C.voice_info_get_voice_path(voiceInfo.handle)
	voiceInfo.VoicePath = C.GoString(value)
	C.property_bag_free_string(value)
	/* Properties */
	var propBagHandle C.SPXHANDLE
	ret = uintptr(C.voice_info_get_property_bag(uintptr2handle(handle), &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	voiceInfo.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	gender := voiceInfo.Properties.GetPropertyByString("Gender", "")
	if gender == "Female" {
		voiceInfo.Gender = common.Female
	} else if gender == "Male" {
		voiceInfo.Gender = common.Male
	} else {
		voiceInfo.Gender = common.GenderUnknown
	}
	return voiceInfo, nil
}
