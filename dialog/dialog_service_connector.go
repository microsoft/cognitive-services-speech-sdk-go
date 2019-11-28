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
// void cgo_dialog_session_stopped(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_dialog_recognized(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_dialog_recognizing(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
// void cgo_dialog_canceled(SPXRECOHANDLE handle, SPXEVENTHANDLE event, void* context);
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
func NewDialogServiceConnectorFromConfig(config DialogServiceConfig, audioConfig *audio.AudioConfig) (*DialogServiceConnector, error) {
	var handle C.SPXHANDLE
	if config == nil {
		return nil, common.NewCarbonError(uintptr(C.SPXERR_INVALID_ARG))
	}
	configHandle := config.getHandle()
	var audioHandle C.SPXHANDLE
	if audioConfig == nil {
		audioHandle = nil
	} else {
		audioHandle = uintptr2handle(audioConfig.GetHandle())
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

// ConnectAsync connects with the back end.
func (connector DialogServiceConnector) ConnectAsync() chan error {
	outcome := make(chan error)
	go func() {
		ret := uintptr(C.dialog_service_connector_connect(connector.handle))
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
		} else {
			outcome <- nil
		}
	}()
	return outcome
}

// DisconnectAsync disconnects from the back end.
func (connector DialogServiceConnector) DisconnectAsync() chan error {
	outcome := make(chan error)
	go func() {
		ret := uintptr(C.dialog_service_connector_disconnect(connector.handle))
		if ret != C.SPX_NOERROR {
			outcome <- common.NewCarbonError(ret)
		} else {
			outcome <- nil
		}
	}()
	return outcome
}

type SendActivityOutcome struct {
	common.OperationOutcome

	// InteractionID is the identifier associated with the interaction
	InteractionID string
}

// SendActivityAsync sends an activity to the backing dialog.
func (connector DialogServiceConnector) SendActivityAsync(message string) chan SendActivityOutcome {
	outcome := make(chan SendActivityOutcome)
	go func() {
		msg := C.CString(message)
		defer C.free(unsafe.Pointer(msg))
		buffer := C.malloc(C.sizeof_char * 37)
		defer C.free(unsafe.Pointer(buffer))
		ret := uintptr(C.dialog_service_connector_send_activity(connector.handle, msg, (*C.char)(buffer)));
		if ret != C.SPX_NOERROR {
			outcome <- SendActivityOutcome{ InteractionID: "", OperationOutcome: common.OperationOutcome{ common.NewCarbonError(ret) } }
		} else {
			interactionID := C.GoString((*C.char)(buffer))
			outcome <- SendActivityOutcome{ InteractionID: interactionID, OperationOutcome: common.OperationOutcome{ nil } }
		}
	}()
	return outcome
}

// ListenOnceAsync starts a listening session that will terminate after the first utterance.
func (connector DialogServiceConnector) ListenOnceAsync() <-chan speech.SpeechRecognitionOutcome {
	outcome := make(chan speech.SpeechRecognitionOutcome)
	go func() {
		var handle C.SPXRESULTHANDLE
		ret := uintptr(C.dialog_service_connector_listen_once(connector.handle, &handle))
		if ret != C.SPX_NOERROR {
			outcome <- speech.SpeechRecognitionOutcome{ Result: nil, OperationOutcome: common.OperationOutcome{ common.NewCarbonError(ret) } }
		} else {
			result, err := speech.NewSpeechRecognitionResultFromHandle(handle2uintptr(handle))
			outcome <- speech.SpeechRecognitionOutcome{ Result: result, OperationOutcome: common.OperationOutcome{ err } }
		}
	}()
	return outcome
}

// SetAuthorizationToken sets the authorization token that will be used for connecting to the service.
// Note: The caller needs to ensure that the authorization token is valid. Before the authorization token
// expires, the caller needs to refresh it by calling this setter with a new valid token.
// Otherwise, the connector will encounter errors during its operation.
func (connector DialogServiceConnector) SetAuthorizationToken(token string) error {
	return connector.Properties.SetProperty(common.SpeechServiceAuthorizationToken, token)
}

// AuthorizationToken is the authorization token.
func (connector DialogServiceConnector) AuthorizationToken() string {
	return connector.Properties.GetProperty(common.SpeechServiceAuthorizationToken, "")
}

// Recognized signals events containing speech recognition results.
func (connector DialogServiceConnector) Recognized(handler speech.SpeechRecognitionEventHandler) {
	registerRecognizedCallback(handler, connector.handle)
	if handler != nil {
		C.dialog_service_connector_recognized_set_callback(
			connector.handle,
			(C.PRECOGNITION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_dialog_recognized)),
			nil)
	} else {
		C.dialog_service_connector_recognized_set_callback(connector.handle, nil, nil)
	}
}

// Recognizing signals events containing intermediate recognition results.
func (connector DialogServiceConnector) Recognizing(handler speech.SpeechRecognitionEventHandler) {
	registerRecognizingCallback(handler, connector.handle)
	if handler != nil {
		C.dialog_service_connector_recognizing_set_callback(
			connector.handle,
			(C.PRECOGNITION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_dialog_recognizing)),
			nil)
	} else {
		C.dialog_service_connector_recognizing_set_callback(connector.handle, nil, nil)
	}
}

// SessionStarted signals the start of a listening session.
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

// SessionStopped signals the end of a listening session.
func (connector DialogServiceConnector) SessionStopped(handler speech.SessionEventHandler) {
	registerSessionStoppedCallback(handler, connector.handle)
	if handler != nil {
		C.dialog_service_connector_session_stopped_set_callback(
			connector.handle,
			(C.PSESSION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_dialog_session_stopped)),
			nil)
	} else {
		C.dialog_service_connector_session_stopped_set_callback(connector.handle, nil, nil)
	}
}

// Canceled signals events relating to the cancellation of an interaction. The event indicates if the reason is a direct cancellation or an error.
func (connector DialogServiceConnector) Canceled(handler speech.SpeechRecognitionCanceledEventHandler) {
	registerCanceledCallback(handler, connector.handle)
	if handler != nil {
		C.dialog_service_connector_canceled_set_callback(
			connector.handle,
			(C.PRECOGNITION_CALLBACK_FUNC)(unsafe.Pointer(C.cgo_dialog_canceled)),
			nil)
	} else {
		C.dialog_service_connector_canceled_set_callback(connector.handle, nil, nil)
	}
}