// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speaker

import (
	// "fmt"
	"testing"
	//"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

func TestNewVoiceProfile(t *testing.T) {
	id := "12345678-abcd-abcd-abcd-12345678abcd"
	profileType := common.VoiceProfileType(3)
	profile, err := NewVoiceProfileFromIdAndType(id, profileType)
	if err != nil {
		t.Error("Unexpected error")
	}
	defer profile.Close()
	profileId, err := profile.Id()
	if err != nil {
		t.Error("id not properly set")
	} else if profileId != id {
		t.Error("id does not match original")
	}
	profType, err := profile.Type()
	if err != nil {
		t.Error("type not properly set")
	} else if profType != profileType {
		t.Error("Voice Profile Type not properly set")
	}
}