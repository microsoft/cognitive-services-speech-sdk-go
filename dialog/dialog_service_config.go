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

// DialogServiceConfig defines base configurations for the dialog service connector object.
type DialogServiceConfig interface {
	SetProperty(id common.PropertyID, value string) error
	GetProperty(id common.PropertyID) string
	SetPropertyByString(name string, value string) error
	GetPropertyByString(name string) string
	SetServiceProperty(name string, value string, channel common.ServicePropertyChannel) error
	SetProxy(hostname string, port uint64) error
	SetProxyWithUsernameAndPassword(hostname string, port uint64, username string, password string) error
	SetLanguage(lang string) error
	GetLanguage() string
	Close()
	getHandle() C.SPXHANDLE
}

type dialogServiceConfigBase struct {
	config speech.SpeechConfig
	handle C.SPXHANDLE
}

// SetProperty sets a property value by ID.
func (config *dialogServiceConfigBase) SetProperty(id common.PropertyID, value string) error {
	return config.config.SetProperty(id, value)
}

// GetProperty gets a property value by ID.
func (config *dialogServiceConfigBase) GetProperty(id common.PropertyID) string {
	return config.config.GetProperty(id)
}

// SetPropertyByString sets a property value by name.
func (config *dialogServiceConfigBase) SetPropertyByString(name string, value string) error {
	return config.config.SetPropertyByString(name, value)
}

// GetPropertyByString gets a property value by name.
func (config *dialogServiceConfigBase) GetPropertyByString(name string) string {
	return config.config.GetPropertyByString(name);
}

// SetServiceProperty sets a property value that will be passed to service using the specified channel.
func (config *dialogServiceConfigBase) SetServiceProperty(name string, value string, channel common.ServicePropertyChannel) error {
	return config.config.SetServiceProperty(name, value, channel)
}

// SetProxy sets proxy configuration
//
// Note: Proxy functionality is not available on macOS. This function will have no effect on this platform.
func (config *dialogServiceConfigBase) SetProxy(hostname string, port uint64) error {
	return config.config.SetProxy(hostname, port);
}

// SetProxyWithUsernameAndPassword sets proxy configuration with username and password
//
// Note: Proxy functionality is not available on macOS. This function will have no effect on this platform.
func (config *dialogServiceConfigBase) SetProxyWithUsernameAndPassword(hostname string, port uint64, username string, password string) error {
	return config.config.SetProxyWithUsernameAndPassword(hostname, port, username, password)
}

// SetLanguage sets the input language to the connector.
// The language is specified in BCP-47 format.
func (config *dialogServiceConfigBase) SetLanguage(lang string) error {
	return config.SetProperty(common.SpeechServiceConnectionRecoLanguage, lang);
}

// GetLanguage gets the input language to the connector.
// The language is specified in BCP-47 format.
func (config *dialogServiceConfigBase) GetLanguage() string {
	return config.GetProperty(common.SpeechServiceConnectionRecoLanguage)
}

// Close disposes the associated resources.
func (config *dialogServiceConfigBase) Close() {
	config.config.Close()
}

func (config *dialogServiceConfigBase) getHandle() C.SPXHANDLE {
	return config.handle
}

// BotFrameworkConfig defines configurations for the dialog service connector object for using a Bot Framework backend.
type BotFrameworkConfig struct {
	dialogServiceConfigBase
}

// NewBotFrameworkConfigFromSubscription creates a bot framework service config instance with the specified subscription
// key and region.
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
	config.config = *speechConfig
	config.handle = handle
	return config, nil
}

// NewBotFrameworkConfigFromAuthorizationToken creates a bot framework service config instance with the specified authorization
// token and region.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token
// expires, the caller needs to refresh it by calling this setter with a new valid token.
// As configuration values are copied when creating a new connector, the new token value will not apply to connectors that have
// already been created.
// For connectors that have been created before, you need to set authorization token of the corresponding connector
// to refresh the token. Otherwise, the connectors will encounter errors during operation.
func NewBotFrameworkConfigFromAuthorizationToken(authorizationToken string, region string) (*BotFrameworkConfig, error) {
	var handle C.SPXHANDLE
	at := C.CString(authorizationToken)
	defer C.free(unsafe.Pointer(at))
	r := C.CString(region)
	defer C.free(unsafe.Pointer(r))
	ret := uintptr(C.bot_framework_config_from_authorization_token(&handle, at, r))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	speechConfig, err := speech.NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}
	config := new(BotFrameworkConfig)
	config.config = *speechConfig
	config.handle = handle
	return config, nil
}

// CustomCommandsConfig defines configurations for the dialog service connector object for using a CustomCommands backend.
type CustomCommandsConfig struct {
	dialogServiceConfigBase
}

// NewCustomCommandsConfigFromSubscription creates a Custom Commands config instance with the specified application id,
// subscription key and region.
func NewCustomCommandsConfigFromSubscription(applicationID string, subscriptionKey string, region string) (*CustomCommandsConfig, error) {
	var handle C.SPXHANDLE
	appID := C.CString(applicationID)
	defer C.free(unsafe.Pointer(appID))
	sk := C.CString(subscriptionKey)
	defer C.free(unsafe.Pointer(sk))
	r := C.CString(region)
	defer C.free(unsafe.Pointer(r))
	ret := uintptr(C.custom_commands_config_from_subscription(&handle, appID, sk, r))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	speechConfig, err := speech.NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}
	config := new(CustomCommandsConfig)
	config.config = *speechConfig
	config.handle = handle
	return config, nil
}

// NewCustomCommandsConfigFromAuthorizationToken creates a Custom Commands config instance with the specified application id
// authorization token and region.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token
// expires, the caller needs to refresh it by calling this setter with a new valid token.
// As configuration values are copied when creating a new connector, the new token value will not apply to connectors that have
// already been created.
// For connectors that have been created before, you need to set authorization token of the corresponding connector
// to refresh the token. Otherwise, the connectors will encounter errors during operation.
func NewCustomCommandsConfigFromAuthorizationToken(applicationID string, authorizationToken string, region string) (*CustomCommandsConfig, error) {
	var handle C.SPXHANDLE
	appID := C.CString(applicationID)
	defer C.free(unsafe.Pointer(appID))
	at := C.CString(authorizationToken)
	defer C.free(unsafe.Pointer(at))
	r := C.CString(region)
	defer C.free(unsafe.Pointer(r))
	ret := uintptr(C.custom_commands_config_from_authorization_token(&handle, appID, at, r))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	speechConfig, err := speech.NewSpeechConfigFromHandle(handle2uintptr(handle))
	if err != nil {
		return nil, err
	}
	config := new(CustomCommandsConfig)
	config.config = *speechConfig
	config.handle = handle
	return config, nil
}

// ApplicationID is the corresponding backend application identifier.
func (config *CustomCommandsConfig) ApplicationID() string {
	return config.GetProperty(common.ConversationApplicationID)
}

// SetApplicationID sets the corresponding backend application identifier.
func (config *CustomCommandsConfig) SetApplicationID(appID string) error {
	return config.SetProperty(common.ConversationApplicationID, appID)
}

