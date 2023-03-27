// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"io"
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <string.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_audio_stream.h>
//
import "C"

// AudioDataStream represents audio data stream used for operating audio data as a stream.
// Added in version 1.17.0
type AudioDataStream struct {
	handle C.SPXHANDLE

	// Properties represents the collection of additional properties.
	Properties *common.PropertyCollection
}

// Close disposes the associated resources.
func (stream AudioDataStream) Close() {
	stream.Properties.Close()
	C.audio_data_stream_release(stream.handle)
}

// NewAudioDataStreamFromHandle creates a new AudioDataStream from a handle (for internal use)
func NewAudioDataStreamFromHandle(handle common.SPXHandle) (*AudioDataStream, error) {
	stream := new(AudioDataStream)
	stream.handle = uintptr2handle(handle)
	/* Properties */
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.audio_data_stream_get_property_bag(uintptr2handle(handle), &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	stream.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return stream, nil
}

// NewAudioDataStreamFromWavFileInput creates a memory backed AudioDataStream for the specified audio input file.
func NewAudioDataStreamFromWavFileInput(filename string) (*AudioDataStream, error) {
	var handle C.SPXHANDLE
	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))
	ret := uintptr(C.audio_data_stream_create_from_file(&handle, fn))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return NewAudioDataStreamFromHandle(handle2uintptr(handle))
}

// NewAudioDataStreamFromSpeechSynthesisResult creates a memory backed AudioDataStream from given speech synthesis result.
func NewAudioDataStreamFromSpeechSynthesisResult(result *SpeechSynthesisResult) (*AudioDataStream, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_data_stream_create_from_result(&handle, result.handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return NewAudioDataStreamFromHandle(handle2uintptr(handle))
}

// GetStatus gets the current status of the audio data stream.
func (stream AudioDataStream) GetStatus() (common.StreamStatus, error) {
	var cStatus C.Stream_Status
	ret := uintptr(C.audio_data_stream_get_status(stream.handle, &cStatus))
	if ret != C.SPX_NOERROR {
		return common.StreamStatusUnknown, common.NewCarbonError(ret)
	}
	return (common.StreamStatus)(cStatus), nil
}

// CanReadData checks whether the stream has enough data to be read.
func (stream AudioDataStream) CanReadData(bytesRequested uint) bool {
	return (bool)(C.audio_data_stream_can_read_data(stream.handle, (C.uint32_t)(bytesRequested)))
}

// CanReadDataAt checks whether the stream has enough data to be read, at the specified offset.
func (stream AudioDataStream) CanReadDataAt(bytesRequested uint, off int64) bool {
	return (bool)(C.audio_data_stream_can_read_data_from_position(stream.handle, (C.uint32_t)(bytesRequested), (C.uint32_t)(off)))
}

// Read reads a chunk of the audio data stream and fill it to given buffer.
// It returns size of data filled to the buffer and any write error encountered.
func (stream AudioDataStream) Read(buffer []byte) (int, error) {
	if len(buffer) == 0 {
		return 0, common.NewCarbonError(0x005) // SPXERR_INVALID_ARG
	}
	var outSize C.uint32_t
	ret := uintptr(C.audio_data_stream_read(stream.handle, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), (C.uint32_t)(len(buffer)), &outSize))
	if ret != C.SPX_NOERROR {
		return 0, common.NewCarbonError(ret)
	}
	if outSize == 0 {
		return 0, io.EOF
	}
	return (int)(outSize), nil
}

// ReadAt reads a chunk of the audio data stream and fill it to given buffer, at specified offset.
// It returns size of data filled to the buffer and any write error encountered.
func (stream AudioDataStream) ReadAt(buffer []byte, off int64) (int, error) {
	if len(buffer) == 0 {
		return 0, common.NewCarbonError(0x005) // SPXERR_INVALID_ARG
	}
	var outSize C.uint32_t
	ret := uintptr(C.audio_data_stream_read_from_position(stream.handle, (*C.uint8_t)(unsafe.Pointer(&buffer[0])), (C.uint32_t)(len(buffer)), (C.uint32_t)(off), &outSize))
	if ret != C.SPX_NOERROR {
		return 0, common.NewCarbonError(ret)
	}
	if outSize == 0 {
		return 0, io.EOF
	}
	return (int)(outSize), nil
}

// SaveToWavFileAsync saves the audio data to a file, asynchronously.
func (stream AudioDataStream) SaveToWavFileAsync(filename string) chan error {
	outcome := make(chan error)
	go func() {
		fn := C.CString(filename)
		defer C.free(unsafe.Pointer(fn))
		ret := uintptr(C.audio_data_stream_save_to_wave_file(stream.handle, fn))
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
		} else {
			outcome <- nil
		}
	}()
	return outcome
}

// GetOffset gets current offset of the audio data stream.
func (stream AudioDataStream) GetOffset() (int, error) {
	var position C.uint32_t
	ret := uintptr(C.audio_data_stream_get_position(stream.handle, &position))
	if ret != C.SPX_NOERROR {
		return 0, common.NewCarbonError(ret)
	}
	return (int)(position), nil
}

// SetOffset sets current offset of the audio data stream.
func (stream AudioDataStream) SetOffset(offset int) error {
	ret := uintptr(C.audio_data_stream_set_position(stream.handle, (C.uint32_t)(offset)))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}
