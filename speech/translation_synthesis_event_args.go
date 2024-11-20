package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_recognizer.h>
import "C"

type TranslationSynthesisEventArgs struct {
	handle C.SPXHANDLE
	Result TranslationResult
}

// Close releases the underlying resources
func (event TranslationSynthesisEventArgs) Close() {
	event.Result.Close()
	C.recognizer_event_handle_release(event.handle)
}

func NewTranslationSynthesisEventArgsFromHandle(handle common.SPXHandle) (*TranslationSynthesisEventArgs, error) {
	event := new(TranslationSynthesisEventArgs)
	event.handle = uintptr2handle(handle)
	var resultHandle C.SPXHANDLE
	ret := uintptr(C.recognizer_recognition_event_get_result(event.handle, &resultHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	result, err := NewTranslationResultFromHandle(handle2uintptr(resultHandle))
	if err != nil {
		return nil, err
	}
	event.Result = *result
	return event, nil
}

type TranslationSynthesisEventHandler func(event TranslationSynthesisEventArgs)
