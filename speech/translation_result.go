package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_recognizer.h>
// #include <speechapi_c_translation_result.h>
//
import "C"

type TranslationResult struct {
	handle    C.SPXHANDLE
	AudioData []byte
}

func (result *TranslationResult) Close() {
	if result.handle != C.SPXHANDLE_INVALID {
		C.recognizer_result_handle_release(result.handle)
		result.handle = C.SPXHANDLE_INVALID
	}
}

func NewTranslationResultFromHandle(handle common.SPXHandle) (*TranslationResult, error) {
	result := new(TranslationResult)
	result.handle = uintptr2handle(handle)
	var size C.size_t
	ret := uintptr(C.translation_synthesis_result_get_audio_data(result.handle, (*C.uint8_t)(nil), &size))
	if ret != C.SPX_NOERROR {
		if ret == C.SPXERR_BUFFER_TOO_SMALL {
			buffer := C.malloc(C.sizeof_char * size)
			defer C.free(unsafe.Pointer(buffer))
			ret = uintptr(C.translation_synthesis_result_get_audio_data(result.handle, (*C.uint8_t)(buffer), &size))
			if ret != C.SPX_NOERROR {
				return nil, common.NewCarbonError(ret)
			}
			result.AudioData = C.GoBytes(buffer, C.int(size))
			return result, nil
		}
		return nil, common.NewCarbonError(ret)
	}
	return result, nil
}
