package audio

import (
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"testing"
)

func TestHandleConversion(t *testing.T) {
	orig := common.SPXHandle(3)
	handle := uintptr2handle(orig)
	dest := handle2uintptr(handle)
	if orig != dest {
		t.Error("Values are not equal")
	}
	if uintptr2handle(dest) != handle {
		t.Error("Values are not equal")
	}
}