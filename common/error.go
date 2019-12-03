//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
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