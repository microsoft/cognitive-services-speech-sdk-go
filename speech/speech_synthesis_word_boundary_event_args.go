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

// SpeechSynthesisWordBoundaryEventArgs represents the speech synthesis word boundary event arguments.
type SpeechSynthesisWordBoundaryEventArgs struct {
	handle      C.SPXHANDLE
	AudioOffset uint64
	TextOffset  uint
	WordLength  uint
}

// Close releases the underlying resources
func (event SpeechSynthesisWordBoundaryEventArgs) Close() {
	C.synthesizer_event_handle_release(event.handle)
}

// NewSpeechSynthesisWordBoundaryEventArgsFromHandle creates the object from the handle (for internal use)
func NewSpeechSynthesisWordBoundaryEventArgsFromHandle(handle common.SPXHandle) (*SpeechSynthesisWordBoundaryEventArgs, error) {
	event := new(SpeechSynthesisWordBoundaryEventArgs)
	event.handle = uintptr2handle(handle)
	var cAudioOffset C.uint64_t
	var cTextOffset, cWordLength C.uint32_t
	ret := uintptr(C.synthesizer_word_boundary_event_get_values(event.handle, &cAudioOffset, &cTextOffset, &cWordLength))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	event.AudioOffset = uint64(cAudioOffset)
	event.TextOffset = uint(cTextOffset)
	event.WordLength = uint(cWordLength)
	return event, nil
}

// SpeechSynthesisWordBoundaryEventHandler is the type of the event handler that receives SpeechSynthesisWordBoundaryEventArgs
type SpeechSynthesisWordBoundaryEventHandler func(event SpeechSynthesisWordBoundaryEventArgs)
