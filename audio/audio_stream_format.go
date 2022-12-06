// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package audio

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_audio_stream_format.h>
import "C"

// AudioStreamFormat represents the audio stream format used for custom audio input configurations.
// Updated in version 1.5.0.
type AudioStreamFormat struct {
	handle C.SPXHANDLE
}

// GetDefaultInputFormat creates an audio stream format object representing the default audio stream format
// (16 kHz, 16 bit, mono PCM).
func GetDefaultInputFormat() (*AudioStreamFormat, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_format_create_from_default_input(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	format := new(AudioStreamFormat)
	format.handle = handle
	return format, nil
}

// GetWaveFormat creates an audio stream format object with the specified waveformat characteristics.
func GetWaveFormat(samplesPerSecond uint32, bitsPerSample uint8, channels uint8, waveFormat AudioStreamWaveFormat) (*AudioStreamFormat, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_format_create_from_waveformat(
		&handle,
		(C.uint32_t)(samplesPerSecond),
		(C.uint8_t)(bitsPerSample),
		(C.uint8_t)(channels),
		(C.Audio_Stream_Wave_Format)(waveFormat)))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	format := new(AudioStreamFormat)
	format.handle = handle
	return format, nil
}

// GetWaveFormatPCM creates an audio stream format object with the specified PCM waveformat characteristics.
// Note: Currently, only WAV / PCM with 16-bit samples, 16 kHz sample rate, and a single channel (Mono) is supported. When
// used with Conversation Transcription, eight channels are supported.
func GetWaveFormatPCM(samplesPerSecond uint32, bitsPerSample uint8, channels uint8) (*AudioStreamFormat, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_format_create_from_waveformat_pcm(
		&handle,
		(C.uint32_t)(samplesPerSecond),
		(C.uint8_t)(bitsPerSample),
		(C.uint8_t)(channels)))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	format := new(AudioStreamFormat)
	format.handle = handle
	return format, nil
}

// GetDefaultOutputFormat creates an audio stream format object representing the default audio stream format
// (16 kHz, 16 bit, mono PCM).
func GetDefaultOutputFormat() (*AudioStreamFormat, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_format_create_from_default_output(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	format := new(AudioStreamFormat)
	format.handle = handle
	return format, nil
}

// GetCompressedFormat creates an audio stream format object with the specified compressed audio container format, to be
// used as input format.
func GetCompressedFormat(compressedFormat AudioStreamContainerFormat) (*AudioStreamFormat, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_format_create_from_compressed_format(&handle, (C.Audio_Stream_Container_Format)(compressedFormat)))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	format := new(AudioStreamFormat)
	format.handle = handle
	return format, nil
}

// Close disposes the associated resources.
func (format *AudioStreamFormat) Close() {
	C.audio_stream_format_release(format.handle)
}
