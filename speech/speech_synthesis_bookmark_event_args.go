// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"runtime"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_synthesizer.h>
// #include <speechapi_c_property_bag.h>
import "C"

// SpeechSynthesisBookmarkEventArgs represents the speech synthesis bookmark event arguments.
type SpeechSynthesisBookmarkEventArgs struct {
	handle      C.SPXHANDLE
	AudioOffset uint64
	Text        string
}

// Close releases the underlying resources
func (event *SpeechSynthesisBookmarkEventArgs) Close() {
	runtime.SetFinalizer(event, nil)
	C.synthesizer_event_handle_release(event.handle)
}

// NewSpeechSynthesisBookmarkEventArgsFromHandle creates the object from the handle (for internal use)
func NewSpeechSynthesisBookmarkEventArgsFromHandle(handle common.SPXHandle) (*SpeechSynthesisBookmarkEventArgs, error) {
	event := new(SpeechSynthesisBookmarkEventArgs)
	event.handle = uintptr2handle(handle)
	runtime.SetFinalizer(event, (*SpeechSynthesisBookmarkEventArgs).Close)
	/* AudioOffset */
	var cAudioOffset C.uint64_t
	ret := uintptr(C.synthesizer_bookmark_event_get_values(event.handle, &cAudioOffset))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	event.AudioOffset = uint64(cAudioOffset)
	/* Text */
	value := C.synthesizer_bookmark_event_get_text(event.handle)
	event.Text = C.GoString(value)
	C.property_bag_free_string(value)
	return event, nil
}

// SpeechSynthesisBookmarkEventHandler is the type of the event handler that receives SpeechSynthesisBookmarkEventArgs
type SpeechSynthesisBookmarkEventHandler func(event SpeechSynthesisBookmarkEventArgs)
