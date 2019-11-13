package speech

// #include <speechapi_c_common.h>
import "C"
import "unsafe"

// SPXHandle is the internal handle type
type SPXHandle uintptr

func uintptr2handle(h SPXHandle) C.SPXHANDLE {
	return (C.SPXHANDLE)(unsafe.Pointer(h))
}

func handle2uintptr(h C.SPXHANDLE) SPXHandle {
	return (SPXHandle)(unsafe.Pointer(h))
}