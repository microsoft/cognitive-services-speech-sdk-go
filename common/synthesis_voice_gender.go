// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// SynthesisVoiceGender defines the gender of a synthesis voice.
type SynthesisVoiceGender int

const (
	// GenderUnknown means the gender is unknown.
	GenderUnknown SynthesisVoiceGender = 0

	// Female indicates female.
	Female SynthesisVoiceGender = 1

	// Male indicates male.
	Male SynthesisVoiceGender = 2
)
