package dialog

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// #include <stdlib.h>
// #include <speechapi_c_common.h>
// #include <speechapi_c_factory.h>
// #include <speechapi_c_dialog_service_connector.h>
//
// /* Proxy functions forward declarations */
// void cgo_dialog_session_started(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
//
import "C"
import "unsafe"

// DialogServiceConnector connects to a speech enabled dialog backend.
type DialogServiceConnector struct {
	Properties common.PropertyCollection
	handle C.SPXHANDLE
}

func newDialogServiceConnectorFromHandle(handle C.SPXHANDLE) (*DialogServiceConnector, error) {
	var propBagHandle C.SPXHANDLE
	ret := uintptr(C.dialog_service_connector_get_property_bag(handle, &propBagHandle))
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	connector := new(DialogServiceConnector)
	connector.handle = handle
	connector.Properties = common.NewPropertyCollectionFromHandle(handle2uintptr(propBagHandle))
	return connector, nil
}

// NewDialogServiceConnectorFromConfig creates a dialog service connector from a dialog service config and an audio config.
// Users should use this function to create a dialog service connector.
func NewDialogServiceConnectorFromConfig(config *DialogServiceConfig, audioConfig *audio.AudioConfig) (*DialogServiceConnector, error) {
	var handle C.SPXHANDLE
	if config == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	configHandle := (*config).getHandle()
	var audioHandle C.SPXHANDLE
	if audioConfig == nil {
		audioHandle = nil
	} else {
		audioHandle = uintptr2handle((*audioConfig).GetHandle())
	}
	ret := uintptr(C.dialog_service_connector_create_dialog_service_connector_from_config(&handle, configHandle, audioHandle));
	if ret != C.SPX_NOERROR {
		return nil, common.NewCarbonError(ret)
	}
	return newDialogServiceConnectorFromHandle(handle)
}

// Close performs cleanup of resources.
func (connector DialogServiceConnector) Close() {
	connector.Properties.Close()
	C.dialog_service_connector_handle_release(connector.handle)
}

func (connector DialogServiceConnector) SessionStarted(handler speech.SessionEventHandler) {
	registerSessionStartedCallback(handler, connector.handle)
	if handler != nil {
		C.dialog_service_connector_session_started_set_callback(
			connector.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_dialog_session_started)),
			nil)
	} else {
		C.dialog_service_connector_session_started_set_callback(connector.handle, nil, nil)
	}
}