// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"os"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

var timeout time.Duration = 20 * time.Second

func createSpeechSynthesizerFromSubscriptionRegionAndAudioConfig(t *testing.T, subscription string, region string, audioConfig *audio.AudioConfig) *SpeechSynthesizer {
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	defer config.Close()
	config.SetProperty(common.SpeechLogFilename, "go_synthesizer.log")
	synthesizer, err := NewSpeechSynthesizerFromConfig(config, audioConfig)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	return synthesizer
}

func createSpeechSynthesizerFromAudioConfig(t *testing.T, audioConfig *audio.AudioConfig) *SpeechSynthesizer {
	subscription := os.Getenv("SR_SUBSCRIPTION_KEY")
	region := os.Getenv("SR_SUBSCRIPTION_REGION")
	return createSpeechSynthesizerFromSubscriptionRegionAndAudioConfig(t, subscription, region, audioConfig)
}

func checkSynthesisResult(t *testing.T, result *SpeechSynthesisResult, reason common.ResultReason) {
	if result == nil {
		t.Error("Synthesis Result is nil.")
	}
	if result.Reason != reason {
		t.Errorf("Synthesis result reason mismatch. expected %v, got %v", reason, result.Reason)
	}
	if reason == common.SynthesizingAudioStarted {
		if len(result.AudioData) != 0 {
			t.Errorf("Synthesized audio should be empty for SynthesizingAudioStarted, got size %d.", len(result.AudioData))
		}
	} else {
		if len(result.AudioData) == 0 {
			t.Errorf("Synthesized audio is empty")
		}
	}
}

func TestSynthesizerEvents(t *testing.T) {
	synthesizer := createSpeechSynthesizerFromAudioConfig(t, nil)
	if synthesizer == nil {
		t.Error("synthesizer creation failed")
		return
	}
	defer synthesizer.Close()
	synthesisStartedFuture := make(chan bool)
	synthesizer.SynthesisStarted(func(event SpeechSynthesisEventArgs) {
		defer event.Close()
		t.Log("SynthesisStarted")
		checkSynthesisResult(t, &event.Result, common.SynthesizingAudioStarted)
		synthesisStartedFuture <- true
	})
	synthesizingFuture := make(chan string)
	synthesizer.Synthesizing(func(event SpeechSynthesisEventArgs) {
		defer event.Close()
		t.Logf("Synthesizing, audio chunk length %d", len(event.Result.AudioData))
		checkSynthesisResult(t, &event.Result, common.SynthesizingAudio)
		select {
		case synthesizingFuture <- "Synthesizing":
		default:
		}
	})
	synthesisCompletedFuture := make(chan string)
	synthesizer.SynthesisCompleted(func(event SpeechSynthesisEventArgs) {
		defer event.Close()
		t.Logf("SynthesisCompleted, audio length %d", len(event.Result.AudioData))
		checkSynthesisResult(t, &event.Result, common.SynthesizingAudioCompleted)
		synthesisCompletedFuture <- "synthesisCompletedFuture"
	})
	resultFuture := synthesizer.SpeakTextAsync("test")
	select {
	case <-synthesisStartedFuture:
	case <-time.After(timeout):
		t.Error("Timeout waiting for SynthesisStarted event.")
	}
	select {
	case <-synthesizingFuture:
		t.Logf("Received at least one Synthesizing event.")
	case <-time.After(timeout):
		t.Error("Timeout waiting for Synthesizing event.")
	}
	select {
	case <-synthesisCompletedFuture:
		t.Logf("Received synthesisCompletedFuture event.")
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesisCompletedFuture event.")
	}
	select {
	case result := <-resultFuture:
		defer result.Close()
		checkSynthesisResult(t, result.Result, common.SynthesizingAudioCompleted)
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}
}
