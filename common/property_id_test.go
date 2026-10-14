// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

import "testing"

func TestSpeechServiceConnectionEnableIPv6PropertyID(t *testing.T) {
	const expected PropertyID = 1106
	if SpeechServiceConnectionEnableIPv6 != expected {
		t.Fatalf("SpeechServiceConnectionEnableIPv6 = %d, expected %d", SpeechServiceConnectionEnableIPv6, expected)
	}
}
