package common

// OperationOutcome is the base type of operation outcomes.
type OperationOutcome struct {
	// Error is present (not nil) if the operation failed
	Error error
}

// Failed checks if the operation failed
func (outcome OperationOutcome) Failed() bool {
	return outcome.Error != nil
}