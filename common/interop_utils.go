// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// #include <speechapi_c_common.h>
import "C"
import "unsafe"

// SPXHandle is the internal handle type
type SPXHandle uintptr

func uintptr2handle(h SPXHandle) C.SPXHANDLE {
	return (C.SPXHANDLE)(unsafe.Pointer(h)) //nolint:govet
}
