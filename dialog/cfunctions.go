package dialog
// This file defines the proxy functions required to use callbacks

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// extern void dialogFireEventSessionStarted(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_dialog_session_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//    dialogFireEventSessionStarted(handle, event);
// }
import "C"