// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

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
// extern void synthesizerFireEventSynthesisStarted(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_synthesizer_synthesis_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     synthesizerFireEventSynthesisStarted(handle, event);
// }
//
// extern void synthesizerFireEventSynthesizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_synthesizer_synthesizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     synthesizerFireEventSynthesizing(handle, event);
// }
//
// extern void synthesizerFireEventSynthesisCompleted(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_synthesizer_synthesis_completed(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     synthesizerFireEventSynthesisCompleted(handle, event);
// }
//
// extern void synthesizerFireEventSynthesisCanceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_synthesizer_synthesis_canceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     synthesizerFireEventSynthesisCanceled(handle, event);
// }
//
// extern void synthesizerFireEventWordBoundary(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_synthesizer_word_boundary(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     synthesizerFireEventWordBoundary(handle, event);
// }
//
// extern void synthesizerFireEventVisemeReceived(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_synthesizer_viseme_received(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     synthesizerFireEventVisemeReceived(handle, event);
// }
//
// extern void synthesizerFireEventBookmarkReached(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_synthesizer_bookmark_reached(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     synthesizerFireEventBookmarkReached(handle, event);
// }
//
// extern void translatorFireEventSynthesizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event);
//
// void cgo_translator_synthesizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context)
// {
//     translatorFireEventSynthesizing(handle, event);
// }
import "C"
