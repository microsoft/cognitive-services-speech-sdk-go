package common

type CarbonError struct {
	Code uintptr
}

func NewCarbonError(code uintptr) CarbonError {
	var error CarbonError
	error.Code = code
	return error
}

func (e CarbonError) Error() string {
	return "";
}