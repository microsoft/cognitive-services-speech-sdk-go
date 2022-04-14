// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_synthesizer.h>
// #include <speechapi_c_property_bag.h>
import "C"

// SpeechSynthesisWordBoundaryEventArgs represents the speech synthesis word boundary event arguments.
type SpeechSynthesisWordBoundaryEventArgs struct {
	handle C.SPXHANDLE

	// AudioOffset is the audio offset of the word boundary event, in ticks (100 nanoseconds).
	AudioOffset uint64

	// Duration is the duration of the word boundary event.
	Duration time.Duration

	// TextOffset is the text offset.
	TextOffset uint

	// WordLength is the length of the word.
	WordLength uint

	// Text is the text.
	Text string

	// BoundaryType is the boundary type.
	BoundaryType common.SpeechSynthesisBoundaryType
}

// Close releases the underlying resources
func (event SpeechSynthesisWordBoundaryEventArgs) Close() {
	C.synthesizer_event_handle_release(event.handle)
}

// NewSpeechSynthesisWordBoundaryEventArgsFromHandle creates the object from the handle (for internal use)
func NewSpeechSynthesisWordBoundaryEventArgsFromHandle(handle common.SPXHandle) (*SpeechSynthesisWordBoundaryEventArgs, error) {
	event := new(SpeechSynthesisWordBoundaryEventArgs)
	event.handle = uintptr2handle(handle)
	var cAudioOffset, cDuration C.uint64_t
	var cTextOffset, cWordLength C.uint32_t
	var cBoundaryType C.SpeechSynthesis_BoundaryType
	ret := uintptr(C.synthesizer_word_boundary_event_get_values(event.handle, &cAudioOffset, &cDuration, &cTextOffset, &cWordLength, &cBoundaryType))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	event.AudioOffset = uint64(cAudioOffset)
	event.Duration = time.Duration(cDuration*100) * time.Nanosecond
	event.TextOffset = uint(cTextOffset)
	event.WordLength = uint(cWordLength)
	event.BoundaryType = (common.SpeechSynthesisBoundaryType)(cBoundaryType)
	/* Text */
	value := C.synthesizer_event_get_text(event.handle)
	event.Text = C.GoString(value)
	C.property_bag_free_string(value)
	return event, nil
}

// SpeechSynthesisWordBoundaryEventHandler is the type of the event handler that receives SpeechSynthesisWordBoundaryEventArgs
type SpeechSynthesisWordBoundaryEventHandler func(event SpeechSynthesisWordBoundaryEventArgs)
