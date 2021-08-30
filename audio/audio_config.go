// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package audio

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_property_bag.h>
// #include <speechapi_c_audio_config.h>
import "C"

// AudioConfig represents specific audio configuration, such as microphone, file, or custom audio streams.
type AudioConfig struct {
	handle     C.SPXHANDLE
	properties *common.PropertyCollection
}

// GetHandle gets the handle to the resource (for internal use)
func (config AudioConfig) GetHandle() common.SPXHandle {
	return handle2uintptr(config.handle)
}

// Close releases the underlying resources
func (config AudioConfig) Close() {
	config.properties.Close()
	C.audio_config_release(config.handle)
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

// NewAudioConfigFromWavFileInput creates an AudioConfig object representing the specified file.
func NewAudioConfigFromWavFileInput(filename string) (*AudioConfig, error) {
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
func NewAudioConfigFromDefaultSpeakerOutput() (*AudioConfig, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_config_create_audio_output_from_default_speaker(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

// NewAudioConfigFromSpeakerOutput creates an AudioConfig object representing the specific audio output device
// (speaker) on the system.
func NewAudioConfigFromSpeakerOutput(deviceName string) (*AudioConfig, error) {
	var handle C.SPXHANDLE
	dn := C.CString(deviceName)
	defer C.free(unsafe.Pointer(dn))
	ret := uintptr(C.audio_config_create_audio_output_from_a_speaker(&handle, dn))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

// NewAudioConfigFromWavFileOutput creates an AudioConfig object representing the specified file for audio output.
func NewAudioConfigFromWavFileOutput(filename string) (*AudioConfig, error) {
	var handle C.SPXHANDLE
	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))
	ret := uintptr(C.audio_config_create_audio_output_from_wav_file_name(&handle, fn))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

// NewAudioConfigFromStreamOutput creates an AudioConfig object representing the specified output stream.
func NewAudioConfigFromStreamOutput(stream AudioOutputStream) (*AudioConfig, error) {
	var handle C.SPXHANDLE
	streamHandle := stream.getHandle()
	ret := uintptr(C.audio_config_create_audio_output_from_stream(&handle, streamHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAudioConfigFromHandle(handle)
}

// SetProperty sets a property value by ID.
func (config AudioConfig) SetProperty(id common.PropertyID, value string) error {
	return config.properties.SetProperty(id, value)
}

// GetProperty gets a property value by ID.
func (config AudioConfig) GetProperty(id common.PropertyID) string {
	return config.properties.GetProperty(id, "")
}

// SetPropertyByString sets a property value by name.
func (config AudioConfig) SetPropertyByString(name string, value string) error {
	return config.properties.SetPropertyByString(name, value)
}

// GetPropertyByString gets a property value by name.
func (config AudioConfig) GetPropertyByString(name string) string {
	return config.properties.GetPropertyByString(name, "")
}
