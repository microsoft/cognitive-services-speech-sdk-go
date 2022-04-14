// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_synthesizer.h>
// #include <speechapi_c_property_bag.h>
import "C"

// SpeechSynthesisVisemeEventArgs represents the speech synthesis viseme event arguments.
type SpeechSynthesisVisemeEventArgs struct {
	handle C.SPXHANDLE

	// AudioOffset is the audio offset of the viseme event, in ticks (100 nanoseconds).
	AudioOffset uint64

	// VisemeID is the viseme ID.
	VisemeID uint

	// Animation is the animation.
	Animation string
}

// Close releases the underlying resources
func (event SpeechSynthesisVisemeEventArgs) Close() {
	C.synthesizer_event_handle_release(event.handle)
}

// NewSpeechSynthesisVisemeEventArgsFromHandle creates the object from the handle (for internal use)
func NewSpeechSynthesisVisemeEventArgsFromHandle(handle common.SPXHandle) (*SpeechSynthesisVisemeEventArgs, error) {
	event := new(SpeechSynthesisVisemeEventArgs)
	event.handle = uintptr2handle(handle)
	/* AudioOffset and VisemeID */
	var cAudioOffset C.uint64_t
	var cVisemeID C.uint32_t
	ret := uintptr(C.synthesizer_viseme_event_get_values(event.handle, &cAudioOffset, &cVisemeID))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	event.AudioOffset = uint64(cAudioOffset)
	event.VisemeID = uint(cVisemeID)
	/* Animation */
	value := C.synthesizer_viseme_event_get_animation(event.handle)
	event.Animation = C.GoString(value)
	C.property_bag_free_string(value)
	return event, nil
}

// SpeechSynthesisVisemeEventHandler is the type of the event handler that receives SpeechSynthesisVisemeEventArgs
type SpeechSynthesisVisemeEventHandler func(event SpeechSynthesisVisemeEventArgs)
