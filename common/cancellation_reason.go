//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package common

// CancellationReason defines the possible reasons a recognition result might be canceled.
type CancellationReason int

const (
	// Indicates that an error occurred during speech recognition.
	Error CancellationReason = 1

	// Indicates that the end of the audio stream was reached.
	EndOfStream CancellationReason = 2
)
