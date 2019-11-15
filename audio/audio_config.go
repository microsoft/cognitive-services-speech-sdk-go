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

// NewAudioConfigFromFileInput creates an AudioConfig object representing the specified file.
func NewAudioConfigFromFileInput(filename string) (*AudioConfig, error) {
	var handle C.SPXHANDLE
	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))
	ret := uintptr(C.audio_config_create_audio_input_from_wav_file_name(&handle, fn))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

// NewAudioConfigFromStreamInput creates an AudioConfig object representing the specified stream.
func NewAudioConfigFromStreamInput(stream AudioInputStream) (*AudioConfig, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_config_create_audio_input_from_stream(&handle, stream.getHandle()))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

// NewAudioConfigFromDefaultSpeakerOutput creates an AudioConfig object representing the default audio output device
// (speaker) on the system.
// Added in version 1.4.0
func NewAudioConfigFromDefaultSpeakerOutput() (*AudioConfig, error) {
	var handle C.SPXHANDLE
	ret := C.audio_config_create_audio_output_from_default_speaker(&handle)
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

// NewAudioConfigFromWavFileOutput creates an AudioConfig object representing the specified file for audio output.
// Added in version 1.4.0
func NewAudioConfigFromWavFileOutput(filename string) (*AudioConfig, error) {
	var handle C.SPXHANDLE
	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))
	ret := C.audio_config_create_audio_output_from_default_speaker(&handle, filename)
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

