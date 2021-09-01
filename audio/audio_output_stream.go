// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package audio

import (
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
// int cgo_audio_push_stream_write_callback_wrapper(void *context, uint8_t* buffer, uint32_t size);
// void cgo_audio_push_stream_close_callback_wrapper(void *context);
import "C"

// AudioOutputStream represents audio output stream used for custom audio output configurations.
// Updated in version 1.7.0
type AudioOutputStream interface {
	Close()
	getHandle() C.SPXHANDLE
}

type audioOutputStreamBase struct {
	handle C.SPXHANDLE
}

func (stream *audioOutputStreamBase) getHandle() C.SPXHANDLE {
	return stream.handle
}

func (stream *audioOutputStreamBase) Close() {
	C.audio_stream_release(stream.handle)
}

// PullAudioOutputStream represents memory backed pull audio output stream used for custom audio output configurations.
type PullAudioOutputStream struct {
	audioOutputStreamBase
}

// NewPullAudioOutputStreamFromHandle creates a new PullAudioOutputStream from a handle (for internal use)
func NewPullAudioOutputStreamFromHandle(handle common.SPXHandle) *PullAudioOutputStream {
	stream := new(PullAudioOutputStream)
	stream.handle = uintptr2handle(handle)
	return stream
}

// CreatePullAudioOutputStream creates a memory backed PullAudioOutputStream.
func CreatePullAudioOutputStream() (*PullAudioOutputStream, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_create_pull_audio_output_stream(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return NewPullAudioOutputStreamFromHandle(handle2uintptr(handle)), nil
}

// Read reads audio from the stream.
// The maximal number of bytes to be read is determined from the size parameter.
// If there is no data immediately available, read() blocks until the next data becomes available.
func (stream PullAudioOutputStream) Read(size uint) ([]byte, error) {
	cBuffer := C.malloc(C.sizeof_char * (C.size_t)(size))
	defer C.free(unsafe.Pointer(cBuffer))
	var outSize C.uint32_t
	ret := uintptr(C.pull_audio_output_stream_read(stream.handle, (*C.uint8_t)(cBuffer), (C.uint32_t)(size), &outSize))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	buffer := C.GoBytes(cBuffer, (C.int)(outSize))
	return buffer, nil
}

// PushAudioOutputStream represents audio output stream used for custom audio output configurations.
type PushAudioOutputStream struct {
	audioOutputStreamBase
}

// PushAudioOutputStreamCallback an interface that defines callback methods (Write() and CloseStream()) for custom audio output
// streams).
type PushAudioOutputStreamCallback interface {
	Write(buffer []byte) int
	CloseStream()
}

var pushStreamCallbacks = make(map[C.SPXHANDLE]PushAudioOutputStreamCallback)

func registerPushStreamCallback(handle C.SPXHANDLE, callback PushAudioOutputStreamCallback) {
	mu.Lock()
	defer mu.Unlock()
	pushStreamCallbacks[handle] = callback
}

func getPushStreamCallback(handle C.SPXHANDLE) *PushAudioOutputStreamCallback {
	mu.Lock()
	defer mu.Unlock()
	cb, ok := pushStreamCallbacks[handle]
	if ok {
		return &cb
	}
	return nil
}

//nolint:deadcode
func deregisterPushStreamCallback(handle C.SPXHANDLE) {
	mu.Lock()
	defer mu.Unlock()
	pushStreamCallbacks[handle] = nil
}

//export cgoAudioOutputCallWriteCallback
func cgoAudioOutputCallWriteCallback(handle C.SPXHANDLE, buffer *C.uint8_t, size C.uint32_t) int {
	callback := getPushStreamCallback(handle)
	if callback != nil {
		goBuffer := C.GoBytes(unsafe.Pointer(buffer), (C.int)(size))
		return (*callback).Write(goBuffer)
	}
	return 0
}

//export cgoAudioOutputCallCloseCallback
func cgoAudioOutputCallCloseCallback(handle C.SPXHANDLE) {
	callback := getPushStreamCallback(handle)
	if callback != nil {
		(*callback).CloseStream()
	}
}

// CreatePushAudioOutputStream creates a PushAudioOutputStream that delegates to the specified callback interface for Write()
// and CloseStream() methods.
func CreatePushAudioOutputStream(callback PushAudioOutputStreamCallback) (*PushAudioOutputStream, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_create_push_audio_output_stream(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	ret = uintptr(C.push_audio_output_stream_set_callbacks(
		handle,
		unsafe.Pointer(handle),
		(C.CUSTOM_AUDIO_PUSH_STREAM_WRITE_CALLBACK)(unsafe.Pointer(C.cgo_audio_push_stream_write_callback_wrapper)),
		(C.CUSTOM_AUDIO_PUSH_STREAM_CLOSE_CALLBACK)(unsafe.Pointer(C.cgo_audio_push_stream_close_callback_wrapper))))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	registerPushStreamCallback(handle, callback)
	stream := new(PushAudioOutputStream)
	stream.handle = handle
	return stream, nil
}
