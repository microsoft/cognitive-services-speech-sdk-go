//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package speech

// This file defines the proxy functions required to use callbacks

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// extern void recognizerFireEventSessionStarted(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_recognizer_session_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     recognizerFireEventSessionStarted(handle, event);
// }
//
// extern void recognizerFireEventSessionStopped(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_recognizer_session_stopped(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     recognizerFireEventSessionStopped(handle, event);
// }
//
// extern void recognizerFireEventSpeechStartDetected(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_recognizer_speech_start_detected(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     recognizerFireEventSpeechStartDetected(handle, event);
// }
//
// extern void recognizerFireEventSpeechEndDetected(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_recognizer_speech_end_detected(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     recognizerFireEventSpeechEndDetected(handle, event);
// }
//
// extern void recognizerFireEventRecognized(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_recognizer_recognized(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     recognizerFireEventRecognized(handle, event);
// }
//
// extern void recognizerFireEventRecognizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_recognizer_recognizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     recognizerFireEventRecognizing(handle, event);
// }
//
// extern void recognizerFireEventCanceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_recognizer_canceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     recognizerFireEventCanceled(handle, event);
// }
//
import "C"
