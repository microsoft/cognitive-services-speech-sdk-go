package common

// ProfanityOption defines the profanity option.
type ProfanityOption int

const (
	// Masked profanity option.
	Masked  ProfanityOption = 0

	// Removed profanity option
	Removed ProfanityOption = 1

	// Raw profanity option
	Raw     ProfanityOption = 2
)
