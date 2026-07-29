package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdbool.h>
// #include <speechapi_c_synthesis_request.h>
//
// SPXHR create_text_stream_request(SPXREQUESTHANDLE* hrequest) {
//     return speech_synthesis_request_create(true, false, NULL, 0, hrequest);
// }
import "C"

// SpeechSynthesisRequestInputType defines the type of input for speech synthesis request.
type SpeechSynthesisRequestInputType int

const (
	// SpeechSynthesisRequestInputType_TextStream indicates that the input is a text stream.
	SpeechSynthesisRequestInputType_TextStream SpeechSynthesisRequestInputType = 0
)

// SpeechSynthesisRequest represents a speech synthesis request.
type SpeechSynthesisRequest struct {
	handle C.SPXHANDLE
}

// NewSpeechSynthesisRequest creates a new speech synthesis request.
func NewSpeechSynthesisRequest(inputType SpeechSynthesisRequestInputType) (*SpeechSynthesisRequest, error) {
	var handle C.SPXHANDLE
	// Currently only TextStream is supported via this API in our extension
    // Use C helper function to ensure bools are passed correctly
    if inputType == SpeechSynthesisRequestInputType_TextStream {
	    ret := uintptr(C.create_text_stream_request(&handle))
	    if ret != C.SPX_NOERROR {
		    return nil, common.NewCarbonError(ret)
	    }
        return &SpeechSynthesisRequest{handle: handle}, nil
    }

    return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
}

// Close releases the resources associated with the request.
func (req *SpeechSynthesisRequest) Close() {
	if req.handle != C.SPXHANDLE_INVALID {
		C.speech_synthesis_request_release(req.handle)
		req.handle = C.SPXHANDLE_INVALID
	}
}

// InputStream returns the input stream for the request.
func (req *SpeechSynthesisRequest) InputStream() *SpeechSynthesisInputStream {
	return &SpeechSynthesisInputStream{request: req}
}

// SpeechSynthesisInputStream represents the input stream for speech synthesis.
type SpeechSynthesisInputStream struct {
	request *SpeechSynthesisRequest
}

// Write writes text to the input stream.
func (stream *SpeechSynthesisInputStream) Write(text string) error {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	ret := uintptr(C.speech_synthesis_request_send_text_piece(stream.request.handle, cText, C.uint32_t(len(text))))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// Close closes the input stream.
func (stream *SpeechSynthesisInputStream) Close() error {
	ret := uintptr(C.speech_synthesis_request_finish(stream.request.handle))
	if ret != C.SPX_NOERROR {
		return common.NewCarbonError(ret)
	}
	return nil
}

// SetVoice sets the voice for the request.
func (req *SpeechSynthesisRequest) SetVoice(voice string) error {
    cVoice := C.CString(voice)
    defer C.free(unsafe.Pointer(cVoice))
    ret := uintptr(C.speech_synthesis_request_set_voice(req.handle, cVoice, nil, nil))
    if ret != C.SPX_NOERROR {
        return common.NewCarbonError(ret)
    }
    return nil
}

func (req *SpeechSynthesisRequest) getHandle() C.SPXHANDLE {
    return req.handle
}
