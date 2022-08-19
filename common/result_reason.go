// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// ResultReason specifies the possible reasons a recognition result might be generated.
type ResultReason int

const (
	// NoMatch indicates speech could not be recognized. More details can be found in the NoMatchDetails object.
	NoMatch ResultReason = 0

	// Canceled indicates that the recognition was canceled. More details can be found using the CancellationDetails object.
	Canceled ResultReason = 1

	// RecognizingSpeech indicates the speech result contains hypothesis text.
	RecognizingSpeech ResultReason = 2

	// RecognizedSpeech indicates the speech result contains final text that has been recognized.
	// Speech Recognition is now complete for this phrase.
	RecognizedSpeech ResultReason = 3

	// RecognizingIntent indicates the intent result contains hypothesis text and intent.
	RecognizingIntent ResultReason = 4

	// RecognizedIntent indicates the intent result contains final text and intent.
	// Speech Recognition and Intent determination are now complete for this phrase.
	RecognizedIntent ResultReason = 5

	// TranslatingSpeech indicates the translation result contains hypothesis text and its translation(s).
	TranslatingSpeech ResultReason = 6

	// TranslatedSpeech indicates the translation result contains final text and corresponding translation(s).
	// Speech Recognition and Translation are now complete for this phrase.
	TranslatedSpeech ResultReason = 7

	// SynthesizingAudio indicates the synthesized audio result contains a non-zero amount of audio data
	SynthesizingAudio ResultReason = 8

	// SynthesizingAudioCompleted indicates the synthesized audio is now complete for this phrase.
	SynthesizingAudioCompleted ResultReason = 9

	// RecognizingKeyword indicates the speech result contains (unverified) keyword text.
	RecognizingKeyword ResultReason = 10

	// RecognizedKeyword indicates that keyword recognition completed recognizing the given keyword.
	RecognizedKeyword ResultReason = 11

	// SynthesizingAudioStarted indicates the speech synthesis is now started
	SynthesizingAudioStarted ResultReason = 12

	// EnrollingVoiceProfile indicates the voice profile is being enrolling and customers need to send more audio to create a voice profile.
	EnrollingVoiceProfile ResultReason = 17

	// EnrolledVoiceProfile indicates the voice profile has been enrolled.
	EnrolledVoiceProfile ResultReason = 18

	// RecognizedSpeakers indicates some speakers have been successfully identified.
	RecognizedSpeakers ResultReason = 19

	// RecognizedSpeaker indicates one speaker has been successfully verified.
	RecognizedSpeaker ResultReason = 20

	// ResetVoiceProfile indicates the voice profile has been reset successfully.
	ResetVoiceProfile ResultReason = 21

	// DeletedVoiceProfile indicates the voice profile has been deleted successfully.
	DeletedVoiceProfile ResultReason = 22

	// VoicesListRetrieved indicates the voices list has been retrieved successfully.
	VoicesListRetrieved ResultReason = 23
)

//go:generate stringer -type=ResultReason -output=result_reason_string.go
