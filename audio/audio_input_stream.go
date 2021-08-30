// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package audio

import (
	"sync"
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <string.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_audio_stream.h>
//
// /* Proxy functions forward declarations */
// int cgo_audio_read_callback_wrapper(void *context, uint8_t *buffer, uint32_t size);
// void cgo_audio_get_property_callback_wrapper(void* context, int id, uint8_t* value, uint32_t size);
// void cgo_audio_close_callback_wrapper(void *context);
import "C"

// AudioInputStream represents audio input stream used for custom audio input configurations
type AudioInputStream interface {
	Close()
	getHandle() C.SPXHANDLE
}

type audioInputStreamBase struct {
	handle C.SPXHANDLE
}

func (stream audioInputStreamBase) getHandle() C.SPXHANDLE {
	return stream.handle
}

func (stream audioInputStreamBase) Close() {
	C.audio_stream_release(stream.handle)
}

// PushAudioInputStream represents memory backed push audio input stream used for custom audio input configurations.
type PushAudioInputStream struct {
	audioInputStreamBase
}

// CreatePushAudioInputStreamFromFormat creates a memory backed PushAudioInputStream with the specified audio format.
// Currently, only WAV / PCM with 16-bit samples, 16 kHz sample rate, and a single channel (Mono) is supported. When used
// with Conversation Transcription, eight channels are supported.
func CreatePushAudioInputStreamFromFormat(format *AudioStreamFormat) (*PushAudioInputStream, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_create_push_audio_input_stream(&handle, format.handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	stream := new(PushAudioInputStream)
	stream.handle = handle
	return stream, nil
}

// CreatePushAudioInputStream creates a memory backed PushAudioInputStream using the default format (16 kHz, 16 bit, mono PCM).
func CreatePushAudioInputStream() (*PushAudioInputStream, error) {
	format, err := GetDefaultInputFormat()
	if err != nil {
		return nil, err
	}
	return CreatePushAudioInputStreamFromFormat(format)
}

// Write writes the audio data specified by making an internal copy of the data.
// Note: The dataBuffer should not contain any audio header.
func (stream PushAudioInputStream) Write(buffer []byte) error {
	size := uint(len(buffer))
	cBuffer := C.CBytes(buffer)
	defer C.free(cBuffer)
	ret := uintptr(C.push_audio_input_stream_write(stream.handle, (*C.uint8_t)(cBuffer), (C.uint32_t)(size)))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// SetProperty sets value of a property. The properties of the audio data should be set before writing the audio data.
func (stream PushAudioInputStream) SetProperty(id common.PropertyID, value string) error {
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))
	ret := uintptr(C.push_audio_input_stream_set_property_by_id(stream.handle, (C.int)(id), v))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// SetPropertyByName sets value of a property. The properties of the audio data should be set before writing the audio data.
func (stream PushAudioInputStream) SetPropertyByName(name string, value string) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))
	ret := uintptr(C.push_audio_input_stream_set_property_by_name(stream.handle, n, v))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// CloseStream closes the stream.
func (stream PushAudioInputStream) CloseStream() {
	C.push_audio_input_stream_close(stream.handle)
}

// PullAudioInputStream represents audio input stream used for custom audio input configurations.
type PullAudioInputStream struct {
	audioInputStreamBase
}

// PullAudioInputStreamCallback interface that defines callback methods (Read(), GetProperty() and CloseStream()) for custom
// audio input streams).
type PullAudioInputStreamCallback interface {
	Read(maxSize uint32) ([]byte, int)
	GetProperty(id common.PropertyID) string
	CloseStream()
}

var mu sync.Mutex
var pullStreamCallbacks = make(map[C.SPXHANDLE]PullAudioInputStreamCallback)

func registerCallback(handle C.SPXHANDLE, callback PullAudioInputStreamCallback) {
	mu.Lock()
	defer mu.Unlock()
	pullStreamCallbacks[handle] = callback
}

func getCallback(handle C.SPXHANDLE) *PullAudioInputStreamCallback {
	mu.Lock()
	defer mu.Unlock()
	cb, ok := pullStreamCallbacks[handle]
	if ok {
		return &cb
	}
	return nil
}

//nolint:deadcode
func deregisterCallback(handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	pullStreamCallbacks[handle] = nil
}

//export cgoAudioCallReadCallback
func cgoAudioCallReadCallback(handle C.SPXRECOHANDLE, dataBuffer *C.uint8_t, size C.uint32_t) int {
	callback := getCallback(handle)
	if callback != nil {
		goBuffer, readSize := (*callback).Read(uint32(size))
		buffer := C.CBytes(goBuffer)
		defer C.free(buffer)
		C.memcpy(unsafe.Pointer(dataBuffer), buffer, (C.size_t)(readSize))
		return readSize
	}
	return 0
}

//export cgoAudioCallGetPropertyCallback
func cgoAudioCallGetPropertyCallback(handle C.SPXHANDLE, id int, value *C.uint8_t, size C.uint32_t) {
	callback := getCallback(handle)
	if callback != nil {
		propValue := (*callback).GetProperty((common.PropertyID)(id))
		buffer := C.CString(propValue)
		defer C.free(unsafe.Pointer(buffer))
		s := size
		if uintptr(len(propValue)) < uintptr(size) {
			s = (C.uint32_t)(len(propValue))
		}
		C.memcpy(unsafe.Pointer(value), unsafe.Pointer(buffer), (C.size_t)(s))
	}
}

//export cgoAudioCallCloseCallback
func cgoAudioCallCloseCallback(handle C.SPXHANDLE) {
	callback := getCallback(handle)
	if callback != nil {
		(*callback).CloseStream()
	}
}

// CreatePullStreamFromFormat creates a PullAudioInputStream that delegates to the specified callback interface for Read()
// and CloseStream() methods and the specified format.
// Currently, only WAV / PCM with 16-bit samples, 16 kHz sample rate, and a single channel (Mono) is supported. When used with Conversation Transcription, eight channels are supported.
func CreatePullStreamFromFormat(callback PullAudioInputStreamCallback, format *AudioStreamFormat) (*PullAudioInputStream, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_create_pull_audio_input_stream(&handle, format.handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.pull_audio_input_stream_set_callbacks(
		handle,
		unsafe.Pointer(handle),
		(C.CUSTOM_AUDIO_PULL_STREAM_READ_CALLBACK)(unsafe.Pointer(C.cgo_audio_read_callback_wrapper)),
		(C.CUSTOM_AUDIO_PULL_STREAM_CLOSE_CALLBACK)(unsafe.Pointer(C.cgo_audio_close_callback_wrapper))))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.pull_audio_input_stream_set_getproperty_callback(
		handle,
		unsafe.Pointer(handle),
		(C.CUSTOM_AUDIO_PULL_STREAM_GET_PROPERTY_CALLBACK)(unsafe.Pointer(C.cgo_audio_get_property_callback_wrapper))))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	registerCallback(handle, callback)
	stream := new(PullAudioInputStream)
	stream.handle = handle
	return stream, nil
}

// CreatePullStream creates a PullAudioInputStream that delegates to the specified callback interface for Read() and CloseStream()
// methods using the default format (16 kHz, 16 bit, mono PCM).
func CreatePullStream(callback PullAudioInputStreamCallback) (*PullAudioInputStream, error) {
	format, err := GetDefaultInputFormat()
	if err != nil {
		return nil, err
	}
	return CreatePullStreamFromFormat(callback, format)
}
