package dialog

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_speech_config.h>
// #include <speechapi_c_property_bag.h>
// #include <speechapi_c_dialog_service_config.h>
import "C"
import "unsafe"

func handle2uintptr(h C.SPXHANDLE) speech.SPXHandle {
	return (speech.SPXHandle)(unsafe.Pointer(h))
}

// DialogServiceConfig defines base configurations for the dialog service connector object.
type DialogServiceConfig struct {
	config *speech.SpeechConfig
}

// SetProperty sets a property value by ID.
func (config *DialogServiceConfig) SetProperty(id common.PropertyID, value string) error {
	return config.config.SetProperty(id, value)
}

// GetProperty gets a property value by ID.
func (config *DialogServiceConfig) GetProperty(id common.PropertyID) string {
	return config.config.GetProperty(id)
}

// SetPropertyByString sets a property value by name.
func (config *DialogServiceConfig) SetPropertyByString(name string, value string) error {
	return config.config.SetPropertyByString(name, value)
}

// GetPropertyByString gets a property value by name.
func (config *DialogServiceConfig) GetPropertyByString(name string, value string) string {
	return config.config.GetPropertyByString(name);
}

// SetServiceProperty sets a property value that will be passed to service using the specified channel.
func (config *DialogServiceConfig) SetServiceProperty(name string, value string, channel common.ServicePropertyChannel) error {
	return config.config.SetServiceProperty(name, value, channel)
}

// SetProxy sets proxy configuration
//
// Note: Proxy functionality is not available on macOS. This function will have no effect on this platform.
func (config *DialogServiceConfig) SetProxy(hostname string, port uint64) error {
	return config.config.SetProxy(hostname, port);
}

// SetProxyWithUsernameAndPassword sets proxy configuration with username and password
//
// Note: Proxy functionality is not available on macOS. This function will have no effect on this platform.
func (config *DialogServiceConfig) SetProxyWithUsernameAndPassword(hostname string, port uint64, username string, password string) error {
	return config.config.SetProxyWithUsernameAndPassword(hostname, port, username, password)
}

// SetLanguage sets the input language to the connector.
// The language is specified in BCP-47 format.
func (config *DialogServiceConfig) SetLanguage(lang string) error {
	return config.SetProperty(common.SpeechServiceConnectionRecoLanguage, lang);
}

// GetLanguage gets the input language to the connector.
// The language is specified in BCP-47 format.
func (config *DialogServiceConfig) GetLanguage() string {
	return config.GetProperty(common.SpeechServiceConnectionRecoLanguage)
}

type BotFrameworkConfig struct {
	DialogServiceConfig
}

func NewBotFrameworkConfigFromSubscription(subscriptionKey string, region string) (*BotFrameworkConfig, error) {
	var handle C.SPXHANDLE
	sk := C.CString(subscriptionKey)
	defer C.free(unsafe.Pointer(sk))
	r := C.CString(region)
	defer C.free(unsafe.Pointer(r))

	ret := uintptr(C.bot_framework_config_from_subscription(&handle, sk, r))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	speechConfig, err := speech.NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}
	config := new(BotFrameworkConfig)
	config.config = speechConfig
	return config, nil
}