// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_source_lang_config.h>
import "C"

// SourceLanguageConfig defines source language configuration.
type SourceLanguageConfig struct {
	handle     C.SPXHANDLE
	properties *common.PropertyCollection
}

func newSourceLanguageConfigFromHandle(handle C.SPXHANDLE) (*SourceLanguageConfig, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.source_lang_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		C.source_lang_config_release(handle)
		return nil, common.NewCarbonError(ret)
	}
	config := new(SourceLanguageConfig)
	config.handle = handle
	config.properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return config, nil
}

// NewSourceLanguageConfigFromLanguage creates an instance of the SourceLanguageConfig with source language
func NewSourceLanguageConfigFromLanguage(language string) (*SourceLanguageConfig, error) {
	var handle C.SPXHANDLE
	languageCStr := C.CString(language)
	defer C.free(unsafe.Pointer(languageCStr))
	ret := uintptr(C.source_lang_config_from_language(&handle, languageCStr))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSourceLanguageConfigFromHandle(handle)
}

// NewSourceLanguageConfigFromLanguageAndEndpointId creates an instance of the SourceLanguageConfig with source language and custom endpoint id. A custom endpoint id corresponds to custom models.
//nolint:revive
func NewSourceLanguageConfigFromLanguageAndEndpointId(language string, endpointID string) (*SourceLanguageConfig, error) {
	var handle C.SPXHANDLE
	languageCStr := C.CString(language)
	defer C.free(unsafe.Pointer(languageCStr))
	endpointCStr := C.CString(endpointID)
	defer C.free(unsafe.Pointer(endpointCStr))
	ret := uintptr(C.source_lang_config_from_language_and_endpointId(&handle, languageCStr, endpointCStr))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newSourceLanguageConfigFromHandle(handle)
}

func (config SourceLanguageConfig) getHandle() C.SPXHANDLE {
	return config.handle
}

// Close performs cleanup of resources.
func (config SourceLanguageConfig) Close() {
	config.properties.Close()
	C.source_lang_config_release(config.handle)
}
