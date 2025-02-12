// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"sync"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_recognizer.h>
//
import "C"

var (
	translationRecognizingCallbacks = make(map[C.SPXHANDLE]TranslationRecognitionEventHandler)
	translationRecognizedCallbacks  = make(map[C.SPXHANDLE]TranslationRecognitionEventHandler)
	translationCanceledCallbacks    = make(map[C.SPXHANDLE]TranslationRecognitionCanceledEventHandler)
	translationSynthesisCallbacks   = make(map[C.SPXHANDLE]TranslationSynthesisEventHandler)
	translationCallbacksLock        sync.Mutex
)

func registerTranslationRecognizingCallback(callback TranslationRecognitionEventHandler, handle C.SPXHANDLE) {
	translationCallbacksLock.Lock()
	defer translationCallbacksLock.Unlock()
	translationRecognizingCallbacks[handle] = callback
}

func registerTranslationRecognizedCallback(callback TranslationRecognitionEventHandler, handle C.SPXHANDLE) {
	translationCallbacksLock.Lock()
	defer translationCallbacksLock.Unlock()
	translationRecognizedCallbacks[handle] = callback
}

func registerTranslationCanceledCallback(callback TranslationRecognitionCanceledEventHandler, handle C.SPXHANDLE) {
	translationCallbacksLock.Lock()
	defer translationCallbacksLock.Unlock()
	translationCanceledCallbacks[handle] = callback
}

func registerTranslationSynthesisCallback(callback TranslationSynthesisEventHandler, handle C.SPXHANDLE) {
	translationCallbacksLock.Lock()
	defer translationCallbacksLock.Unlock()
	translationSynthesisCallbacks[handle] = callback
}

//export cgoTranslationRecognizing
func cgoTranslationRecognizing(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	translationCallbacksLock.Lock()
	callback := translationRecognizingCallbacks[handle]
	translationCallbacksLock.Unlock()
	if callback != nil {
		eventArgs, _ := NewTranslationRecognitionEventArgsFromHandle(handle2uintptr(eventHandle))
		callback(*eventArgs)
	}
}

//export cgoTranslationRecognized
func cgoTranslationRecognized(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	translationCallbacksLock.Lock()
	callback := translationRecognizedCallbacks[handle]
	translationCallbacksLock.Unlock()
	if callback != nil {
		eventArgs, _ := NewTranslationRecognitionEventArgsFromHandle(handle2uintptr(eventHandle))
		callback(*eventArgs)
	}
}

//export cgoTranslationCanceled
func cgoTranslationCanceled(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	translationCallbacksLock.Lock()
	callback := translationCanceledCallbacks[handle]
	translationCallbacksLock.Unlock()
	if callback != nil {
		eventArgs, _ := NewTranslationRecognitionCanceledEventArgsFromHandle(handle2uintptr(eventHandle))
		callback(*eventArgs)
	}
}

//export cgoTranslationSynthesis
func cgoTranslationSynthesis(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE) {
	translationCallbacksLock.Lock()
	callback := translationSynthesisCallbacks[handle]
	translationCallbacksLock.Unlock()
	if callback != nil {
		eventArgs, _ := NewTranslationSynthesisEventArgsFromHandle(handle2uintptr(eventHandle))
		callback(*eventArgs)
	}
}
