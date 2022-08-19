// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// CancellationReason defines the possible reasons a recognition result might be canceled.
type CancellationReason int

const (
	// Error indicates that an error occurred during speech recognition.
	Error CancellationReason = 1

	// EndOfStream indicates that the end of the audio stream was reached.
	EndOfStream CancellationReason = 2

	// CancelledByUser indicates that request was cancelled by the user.
	// Added in version 1.17.0
	CancelledByUser CancellationReason = 3
)

//go:generate stringer -type=CancellationReason -output=cancellation_reason_string.go
