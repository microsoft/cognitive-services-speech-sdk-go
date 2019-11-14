package audio

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_property_bag.h>
// #include <speechapi_c_audio_config.h>
import "C"
import "unsafe"

// AudioConfig represents specific audio configuration, such as microphone, file, or custom audio streams.
type AudioConfig struct {
	handle C.SPXHANDLE
	properties common.PropertyCollection
}

func newAudioConfigFromHandle(handle C.SPXHANDLE) (*AudioConfig, error) {
	var propBagHandle C.SPXPROPERTYBAGHANDLE
	ret := uintptr(C.audio_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	config := new(AudioConfig)
	config.handle = handle
	config.properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return config, nil
}

// NewAudioConfigFromDefaultMicrophoneInput creates an AudioConfig object representing the default microphone on the system.
func NewAudioConfigFromDefaultMicrophoneInput() (*AudioConfig, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_config_create_audio_input_from_default_microphone(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

// NewAudioConfigFromMicrophoneInput creates an AudioConfig object representing a specific microphone on the system.
// Added in version 1.3.0.
func NewAudioConfigFromMicrophoneInput(deviceName string) (*AudioConfig, error) {
	var handle C.SPXHANDLE
	dn := C.CString(deviceName)
	defer C.free(unsafe.Pointer(dn))
	ret := uintptr(C.audio_config_create_audio_input_from_a_microphone(&handle, dn))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}