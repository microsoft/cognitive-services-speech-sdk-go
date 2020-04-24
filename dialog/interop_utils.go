// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package dialog

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <speechapi_c_common.h>
import "C"
import "unsafe"

func uintptr2handle(h common.SPXHandle) C.SPXHANDLE {
	return (C.SPXHANDLE)(unsafe.Pointer(h))
}

func handle2uintptr(h C.SPXHANDLE) common.SPXHandle {
	return (common.SPXHandle)(unsafe.Pointer(h))
}
