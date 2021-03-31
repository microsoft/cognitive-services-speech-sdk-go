// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"bytes"
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
	t.Logf("checking synthesis result with result id of %v", result.ResultId)
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

func checkBinaryEqual(t *testing.T, result1 *SpeechSynthesisResult, result2 *SpeechSynthesisResult) {
	if result1 == nil {
		t.Error("result1 is nil.")
	}
	if result2 == nil {
		t.Error("result1 is nil.")
	}
	if !bytes.Equal(result1.AudioData, result2.AudioData) {
		t.Error("result1 is not binary equal with result2.")
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

func TestSynthesizerSpeakingSsml(t *testing.T) {
	synthesizer := createSpeechSynthesizerFromAudioConfig(t, nil)
	if synthesizer == nil {
		t.Error("synthesizer creation failed")
		return
	}
	defer synthesizer.Close()
	synthesizer.Properties.SetProperty(common.SpeechServiceConnectionSynthVoice, "en-GB-George")
	textResultFuture := synthesizer.SpeakTextAsync("text")

	var textResult SpeechSynthesisOutcome
	select {
	case textResult = <-textResultFuture:
		defer textResult.Close()
		checkSynthesisResult(t, textResult.Result, common.SynthesizingAudioCompleted)
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}

	ssml := "<speak xmlns='http://www.w3.org/2001/10/synthesis' xmlns:mstts='http://www.w3.org/2001/mstts' xmlns:emo='http://www.w3.org/2009/10/emotionml' version='1.0' xml:lang='en-US'><voice name='en-GB-George'>text</voice></speak>"
	ssmlResultFuture := synthesizer.SpeakSsmlAsync(ssml)

	select {
	case ssmlResult := <-ssmlResultFuture:
		defer ssmlResult.Close()
		checkSynthesisResult(t, ssmlResult.Result, common.SynthesizingAudioCompleted)
		checkBinaryEqual(t, textResult.Result, ssmlResult.Result)
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}
}
