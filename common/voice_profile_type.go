// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// VoiceProfileType defines the type of scenario a voice profile has created for.
type VoiceProfileType int

const (
	// Text independent speaker identification
	TextIndependentIdentification VoiceProfileType = 1

	// Text dependent speaker verification
	TextDependentVerification VoiceProfileType = 2

	// Text independent speaker verification
	TextIndependentVerification VoiceProfileType = 3
)
