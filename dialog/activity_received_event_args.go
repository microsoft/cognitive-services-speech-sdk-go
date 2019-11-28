package dialog


import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// #include <stdlib.h>
// #include <stdint.h>
// #include <speechapi_c_dialog_service_connector.h>
import "C"
import "unsafe"

type ActivityReceivedEventArgs struct {
	handle C.SPXHANDLE
	Activity string
}

// Close releases the underlying resources
func (event ActivityReceivedEventArgs) Close() {
	C.dialog_service_connector_activity_received_event_release(event.handle)
}

// NewSpeechRecognitionCanceledEventArgsFromHandle creates the object from the handle (for internal use)
func NewActivityReceivedEventArgsFromHandle(handle common.SPXHandle) (*ActivityReceivedEventArgs, error) {
	event := new(ActivityReceivedEventArgs)
	event.handle = uintptr2handle(handle)
	var size C.size_t
	ret := uintptr(C.dialog_service_connector_activity_received_event_get_activity_size(event.handle, &size))
	if ret != C.SPX_NOERROR {
		event.Close()
		return nil, common.NewCarbonError(ret)
	}
	actBuffer := C.malloc(C.sizeof_char * (size + 1))
	defer C.free(unsafe.Pointer(actBuffer))
	ret = uintptr(C.dialog_service_connector_activity_received_event_get_activity(event.handle, (*C.char)(actBuffer), size + 1))
	if ret != C.SPX_NOERROR {
		event.Close()
		return nil, common.NewCarbonError(ret)
	}
	event.Activity = C.GoString((*C.char)(actBuffer))
	return event, nil
}

type ActivityReceivedEventHandler func (event ActivityReceivedEventArgs)