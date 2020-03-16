//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package common

// CancellationErrorCode defines error code in case that CancellationReason is Error.
type CancellationErrorCode int

const (

	// No error.
	// If CancellationReason is EndOfStream, CancellationErrorCode
	// is set to NoError.
	NoError CancellationErrorCode = 0

	// Indicates an authentication error.
	// An authentication error occurs if subscription key or authorization token is invalid, expired,
	// or does not match the region being used.
	AuthenticationFailure CancellationErrorCode = 1

	// Indicates that one or more recognition parameters are invalid or the audio format is not supported.

	BadRequest CancellationErrorCode = 2

	// Indicates that the number of parallel requests exceeded the number of allowed concurrent transcriptions for the subscription.
	TooManyRequests CancellationErrorCode = 3

	// Indicates that the free subscription used by the request ran out of quota.
	Forbidden CancellationErrorCode = 4

	// Indicates a connection error.
	ConnectionFailure CancellationErrorCode = 5

	// Indicates a time-out error when waiting for response from service.
	ServiceTimeout CancellationErrorCode = 6

	// Indicates that an error is returned by the service.
	ServiceError CancellationErrorCode = 7

	// Indicates that the service is currently unavailable.
	ServiceUnavailable CancellationErrorCode = 8

	// Indicates an unexpected runtime error.
	RuntimeError CancellationErrorCode = 9
)
