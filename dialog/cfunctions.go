package dialog
// This file defines the proxy functions required to use callbacks

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// extern void dialogFireEventSessionStarted(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_dialog_session_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     dialogFireEventSessionStarted(handle, event);
// }
//
// extern void dialogFireEventSessionStopped(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_dialog_session_stopped(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     dialogFireEventSessionStopped(handle, event);
// }
//
// extern void dialogFireEventRecognized(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_dialog_recognized(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     dialogFireEventRecognized(handle, event);
// }
//
// extern void dialogFireEventRecognizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_dialog_recognizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     dialogFireEventRecognizing(handle, event);
// }
//
// extern void dialogFireEventCanceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_dialog_canceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     dialogFireEventCanceled(handle, event);
// }
//
// extern void dialogFireEventActivityReceived(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_dialog_activity_received(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     dialogFireEventActivityReceived(handle, event);
// }
//
import "C"