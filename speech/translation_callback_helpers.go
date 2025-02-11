// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"sync"
	"unsafe"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_recognizer.h>
//
import "C"

var (
	translationRecognizingCallbacks = make(map[uintptr]TranslationRecognitionEventHandler)
	translationRecognizedCallbacks  = make(map[uintptr]TranslationRecognitionEventHandler)
	translationCanceledCallbacks    = make(map[uintptr]TranslationRecognitionCanceledEventHandler)
	translationSynthesisCallbacks   = make(map[uintptr]TranslationSynthesisEventHandler)
	translationCallbacksLock        sync.Mutex
)

func registerTranslationRecognizingCallback(callback TranslationRecognitionEventHandler, handle C.SPXHANDLE) {
	translationCallbacksLock.Lock()
	defer translationCallbacksLock.Unlock()
	if callback == nil {
		delete(translationRecognizingCallbacks, uintptr(handle))
	} else {
		translationRecognizingCallbacks[uintptr(handle)] = callback
	}
}

func registerTranslationRecognizedCallback(callback TranslationRecognitionEventHandler, handle C.SPXHANDLE) {
	translationCallbacksLock.Lock()
	defer translationCallbacksLock.Unlock()
	if callback == nil {
		delete(translationRecognizedCallbacks, uintptr(handle))
	} else {
		translationRecognizedCallbacks[uintptr(handle)] = callback
	}
}

func registerTranslationCanceledCallback(callback TranslationRecognitionCanceledEventHandler, handle C.SPXHANDLE) {
	translationCallbacksLock.Lock()
	defer translationCallbacksLock.Unlock()
	if callback == nil {
		delete(translationCanceledCallbacks, uintptr(handle))
	} else {
		translationCanceledCallbacks[uintptr(handle)] = callback
	}
}

func registerTranslationSynthesisCallback(callback TranslationSynthesisEventHandler, handle C.SPXHANDLE) {
	translationCallbacksLock.Lock()
	defer translationCallbacksLock.Unlock()
	if callback == nil {
		delete(translationSynthesisCallbacks, uintptr(handle))
	} else {
		translationSynthesisCallbacks[uintptr(handle)] = callback
	}
}

//export cgoTranslationRecognizing
func cgoTranslationRecognizing(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE, context unsafe.Pointer) {
	translationCallbacksLock.Lock()
	callback := translationRecognizingCallbacks[uintptr(handle)]
	translationCallbacksLock.Unlock()
	if callback != nil {
		eventArgs, _ := NewTranslationRecognitionEventArgsFromHandle(uintptr(eventHandle))
		callback(*eventArgs)
	}
}

//export cgoTranslationRecognized
func cgoTranslationRecognized(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE, context unsafe.Pointer) {
	translationCallbacksLock.Lock()
	callback := translationRecognizedCallbacks[uintptr(handle)]
	translationCallbacksLock.Unlock()
	if callback != nil {
		eventArgs, _ := NewTranslationRecognitionEventArgsFromHandle(uintptr(eventHandle))
		callback(*eventArgs)
	}
}

//export cgoTranslationCanceled
func cgoTranslationCanceled(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE, context unsafe.Pointer) {
	translationCallbacksLock.Lock()
	callback := translationCanceledCallbacks[uintptr(handle)]
	translationCallbacksLock.Unlock()
	if callback != nil {
		eventArgs, _ := NewTranslationRecognitionCanceledEventArgsFromHandle(uintptr(eventHandle))
		callback(*eventArgs)
	}
}

//export cgoTranslationSynthesis
func cgoTranslationSynthesis(handle C.SPXRECOHANDLE, eventHandle C.SPXEVENTHANDLE, context unsafe.Pointer) {
	translationCallbacksLock.Lock()
	callback := translationSynthesisCallbacks[uintptr(handle)]
	translationCallbacksLock.Unlock()
	if callback != nil {
		eventArgs, _ := NewTranslationSynthesisEventArgsFromHandle(uintptr(eventHandle))
		callback(*eventArgs)
	}
}
