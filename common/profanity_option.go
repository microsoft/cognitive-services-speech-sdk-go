// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// ProfanityOption defines the profanity option.
type ProfanityOption int

const (
	// Masked profanity option.
	Masked ProfanityOption = 0

	// Removed profanity option
	Removed ProfanityOption = 1

	// Raw profanity option
	Raw ProfanityOption = 2
)
