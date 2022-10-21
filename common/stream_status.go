// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// StreamStatus defines the possible status of audio data stream.
type StreamStatus int

const (
	// StreamStatusUnknown indicates the audio data stream status is unknown.
	StreamStatusUnknown StreamStatus = 0

	// StreamStatusNoData indicates that the audio data stream contains no data.
	StreamStatusNoData StreamStatus = 1

	// StreamStatusPartialData indicates the audio data stream contains partial data of a speak request.
	StreamStatusPartialData StreamStatus = 2

	// StreamStatusAllData indicates the audio data stream contains all data of a speak request.
	StreamStatusAllData StreamStatus = 3

	// StreamStatusCanceled indicates the audio data stream was canceled.
	StreamStatusCanceled StreamStatus = 4
)

//go:generate stringer -type=StreamStatus -trimprefix=StreamStatus -output=stream_status_string.go
