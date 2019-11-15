package audio

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_audio_stream.h>
import "C"
import "unsafe"

// AudioInputStream represents audio input stream used for custom audio input configurations
type AudioInputStream interface {
	Close()
}

type audioInputStreamBase struct {
	handle C.SPXHANDLE
}

func (stream audioInputStreamBase) Close() {
	C.audio_stream_release(stream.handle)
}

// PushAudioInputStream represents memory backed push audio input stream used for custom audio input configurations.
type PushAudioInputStream struct {
	audioInputStreamBase
}

// CreatePushStreamFromFormat creates a memory backed PushAudioInputStream with the specified audio format.
// Currently, only WAV / PCM with 16-bit samples, 16 kHz sample rate, and a single channel (Mono) is supported. When used
// with Conversation Transcription, eight channels are supported.
func CreatePushStreamFromFormat(format AudioStreamFormat) (*PushAudioInputStream, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.audio_stream_create_push_audio_input_stream(&handle, format.handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	stream := new(PushAudioInputStream)
	stream.handle = handle
	return stream, nil
}

// CreatePushStream creates a memory backed PushAudioInputStream using the default format (16 kHz, 16 bit, mono PCM).
func CreatePushStream() (*PushAudioInputStream, error) {
	format, err := GetDefaultInputFormat()
	if err != nil {
		return nil, err
	}
	return CreatePushStreamFromFormat(*format)
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
// Added in version 1.5.0.
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

// Close closes the stream.
func (stream PushAudioInputStream) Close() {
	C.push_audio_input_stream_close(stream.handle)
	stream.audioInputStreamBase.Close()
}

// PullAudioInputStream represents audio input stream used for custom audio input configurations.
type PullAudioInputStream struct {
	audioInputStreamBase
}
