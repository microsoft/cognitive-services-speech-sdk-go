package dialog
// This file defines the proxy functions required to use callbacks

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// extern void fireEventSessionStarted(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_session_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//    fireEventSessionStarted(handle, event);
// }
import "C"