// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// SpeechSynthesisBoundaryType defines the boundary type of speech synthesis boundary event.
type SpeechSynthesisBoundaryType int

const (
	// WordBoundary indicates word boundary.
	WordBoundary SpeechSynthesisBoundaryType = 0

	// PunctuationBoundary indicates punctuation boundary.
	PunctuationBoundary SpeechSynthesisBoundaryType = 1

	// SentenceBoundary indicates sentence boundary.
	SentenceBoundary SpeechSynthesisBoundaryType = 2
)
