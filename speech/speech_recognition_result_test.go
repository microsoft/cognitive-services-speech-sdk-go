// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestReadRecognitionResultText(t *testing.T) {
	testCases := []struct {
		name          string
		text          string
		expectedCalls int
	}{
		{
			name:          "short text",
			text:          strings.Repeat("a", initialRecognitionResultTextBufferSize-utf8.UTFMax-1),
			expectedCalls: 1,
		},
		{
			name:          "text at initial buffer boundary",
			text:          strings.Repeat("a", initialRecognitionResultTextBufferSize-1),
			expectedCalls: 2,
		},
		{
			name:          "long UTF-8 text",
			text:          strings.Repeat("ก", 400),
			expectedCalls: 2,
		},
		{
			name:          "long four-byte UTF-8 text",
			text:          strings.Repeat("😀", 300),
			expectedCalls: 2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			calls := 0
			text, ret := readRecognitionResultText(func(buffer []byte) uintptr {
				calls++
				length := len(testCase.text)
				if length >= len(buffer) {
					length = len(buffer) - 1
				}
				for length > 0 && !utf8.ValidString(testCase.text[:length]) {
					length--
				}
				copy(buffer, testCase.text[:length])
				buffer[length] = 0
				return 0
			})

			if ret != 0 {
				t.Fatalf("readRecognitionResultText returned error code %#x", ret)
			}
			if text != testCase.text {
				t.Fatalf("readRecognitionResultText returned %d bytes, want %d", len(text), len(testCase.text))
			}
			if calls != testCase.expectedCalls {
				t.Fatalf("readRecognitionResultText called getter %d times, want %d", calls, testCase.expectedCalls)
			}
		})
	}
}

func TestReadRecognitionResultTextPropagatesError(t *testing.T) {
	const expected = uintptr(0x1234)

	text, ret := readRecognitionResultText(func(buffer []byte) uintptr {
		return expected
	})

	if ret != expected {
		t.Fatalf("readRecognitionResultText returned error code %#x, want %#x", ret, expected)
	}
	if text != "" {
		t.Fatalf("readRecognitionResultText returned text %q after an error", text)
	}
}
