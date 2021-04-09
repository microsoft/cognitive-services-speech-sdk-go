// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// SynthesisVoiceType defines the type of a synthesis voice.
type SynthesisVoiceType int

const (
	// OnlineNeural indicates online neural voice.
	OnlineNeural SynthesisVoiceType = 1

	// OnlineStandard indicates online standard voice.
	OnlineStandard SynthesisVoiceType = 2

	// OfflineNeural indicates offline neural voice.
	OfflineNeural SynthesisVoiceType = 3

	// OfflineStandard indicates offline started voice.
	OfflineStandard SynthesisVoiceType = 4
)
