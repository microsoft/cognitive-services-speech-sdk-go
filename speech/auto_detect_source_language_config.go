// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"strings"
	"unsafe"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <speechapi_c_auto_detect_source_lang_config.h>
import "C"

// AutoDetectSourceLanguageConfig defines auto detection source configuration
type AutoDetectSourceLanguageConfig struct {
	handle     C.SPXHANDLE
	properties *common.PropertyCollection
}

func newAutoDetectSourceLanguageConfigFromHandle(handle C.SPXHANDLE) (*AutoDetectSourceLanguageConfig, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.auto_detect_source_lang_config_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		C.auto_detect_source_lang_config_release(handle)
		return nil, common.NewCarbonError(ret)
	}
	config := new(AutoDetectSourceLanguageConfig)
	config.handle = handle
	config.properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return config, nil
}

// NewAutoDetectSourceLanguageConfigFromOpenRange creates an instance of the AutoDetectSourceLanguageConfig with open range as source languages
func NewAutoDetectSourceLanguageConfigFromOpenRange() (*AutoDetectSourceLanguageConfig, error) {
	var handle C.SPXHANDLE
	ret := uintptr(C.create_auto_detect_source_lang_config_from_open_range(&handle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAutoDetectSourceLanguageConfigFromHandle(handle)
}

// NewAutoDetectSourceLanguageConfigFromLanguages creates an instance of the AutoDetectSourceLanguageConfig with source languages
func NewAutoDetectSourceLanguageConfigFromLanguages(languages []string) (*AutoDetectSourceLanguageConfig, error) {
	var handle C.SPXHANDLE
	languageStr := strings.Join(languages, ",")
	languageCStr := C.CString(languageStr)
	defer C.free(unsafe.Pointer(languageCStr))
	ret := uintptr(C.create_auto_detect_source_lang_config_from_languages(&handle, languageCStr))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newAutoDetectSourceLanguageConfigFromHandle(handle)
}

// NewAutoDetectSourceLanguageConfigFromLanguageConfigs creates an instance of the AutoDetectSourceLanguageConfig with a list of source language config
func NewAutoDetectSourceLanguageConfigFromLanguageConfigs(configs []*SourceLanguageConfig) (*AutoDetectSourceLanguageConfig, error) {
	if len(configs) == 0 {
		return nil, common.NewCarbonError(C.SPXERR_INVALID_ARG)
	}
	var handle C.SPXHANDLE
	var ret uintptr
	first := true
	for i := 0; i < len(configs); i++ {
		c := configs[i]
		if c == nil {
			if !first {
				C.auto_detect_source_lang_config_release(handle)
			}
			return nil, common.NewCarbonError(C.SPXERR_INVALID_ARG)
		}
		if first {
			ret = uintptr(C.create_auto_detect_source_lang_config_from_source_lang_config(&handle, c.getHandle()))
			if ret != C.SPX_NOERROR {
				return nil, common.NewCarbonError(ret)
			}
		} else {
			ret = uintptr(C.add_source_lang_config_to_auto_detect_source_lang_config(handle, c.getHandle()))
			if ret != C.SPX_NOERROR {
				return nil, common.NewCarbonError(ret)
			}
		}
	}
	return newAutoDetectSourceLanguageConfigFromHandle(handle)
}

func (config AutoDetectSourceLanguageConfig) getHandle() C.SPXHANDLE {
	return config.handle
}

// Close performs cleanup of resources.
func (config AutoDetectSourceLanguageConfig) Close() {
	config.properties.Close()
	C.auto_detect_source_lang_config_release(config.handle)
}
