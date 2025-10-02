// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"bytes"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

var timeout time.Duration = 10 * time.Second

func createSpeechSynthesizerFromSpeechConfigAndAudioConfig(t *testing.T, speechConfig *SpeechConfig, audioConfig *audio.AudioConfig) *SpeechSynthesizer {
	speechConfig.SetProperty(common.SpeechLogFilename, "go_synthesizer.log")
	synthesizer, err := NewSpeechSynthesizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	if synthesizer == nil {
		t.Error("synthesizer creation failed")
	}
	return synthesizer
}

func createSpeechConfig(t *testing.T) *SpeechConfig {
	subscription := os.Getenv("SPEECH_SUBSCRIPTION_KEY")
	region := os.Getenv("SPEECH_SUBSCRIPTION_REGION")
	config, err := NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		t.Error("Got an error: ", err)
		return nil
	}
	return config
}

func createSpeechSynthesizerFromAudioConfig(t *testing.T, audioConfig *audio.AudioConfig) *SpeechSynthesizer {
	config := createSpeechConfig(t)
	defer config.Close()
	return createSpeechSynthesizerFromSpeechConfigAndAudioConfig(t, config, audioConfig)
}

func checkSynthesisResult(t *testing.T, result *SpeechSynthesisResult, reason common.ResultReason) {
	if result == nil {
		t.Error("Synthesis Result is nil.")
	}
	t.Logf("checking synthesis result with result id of %v", result.ResultID)
	if result.Reason != reason {
		t.Errorf("Synthesis result reason mismatch. expected %v, got %v", reason, result.Reason)
		if result.Reason == common.Canceled {
			cancellation, _ := NewCancellationDetailsFromSpeechSynthesisResult(result)
			t.Errorf("CANCELED: Reason=%v", cancellation.Reason)
			t.Errorf("CANCELED: ErrorCode=%v ErrorDetails=%s",
				cancellation.ErrorCode,
				cancellation.ErrorDetails)
		}
	}
	if result.Reason == common.Canceled {
		return
	}
	if result.Reason == common.SynthesizingAudioStarted {
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
		return
	}
	if result2 == nil {
		t.Error("result1 is nil.")
		return
	}
	if !bytes.Equal(result1.AudioData, result2.AudioData) {
		t.Error("result1 is not binary equal with result2.")
	}
}

func almostEqual(expected, actual, threshold float64) bool {
	return math.Abs(expected-actual) <= threshold
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
		durationFromProperty := (float64)(event.Result.AudioDuration/time.Millisecond)
		durationFromAudioBuffer := (float64)(len(event.Result.AudioData)/32)
		if !almostEqual(durationFromProperty, durationFromAudioBuffer, 150) {
			t.Errorf("Synthesis duration incorrect (%.2f vs %.2f)", durationFromProperty, durationFromAudioBuffer)
		}
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
	synthesizer.Properties.SetProperty(common.SpeechServiceConnectionSynthVoice, "en-GB-SoniaNeural")
	textResultFuture := synthesizer.SpeakTextAsync("text")

	var textResult SpeechSynthesisOutcome
	select {
	case textResult = <-textResultFuture:
		defer textResult.Close()
		checkSynthesisResult(t, textResult.Result, common.SynthesizingAudioCompleted)
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}

	ssml := "<speak xmlns='http://www.w3.org/2001/10/synthesis' xmlns:mstts='http://www.w3.org/2001/mstts' xmlns:emo='http://www.w3.org/2009/10/emotionml' version='1.0' xml:lang='en-US'><voice name='en-GB-SoniaNeural'>text</voice></speak>"
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

func TestSynthesisToAudioDataStream(t *testing.T) {
	config := createSpeechConfig(t)
	config.SetSpeechSynthesisOutputFormat(common.Audio24Khz48KBitRateMonoMp3)
	defer config.Close()
	synthesizer := createSpeechSynthesizerFromSpeechConfigAndAudioConfig(t, config, nil)
	defer synthesizer.Close()

	textResultFuture := synthesizer.SpeakTextAsync("text")
	var textResult SpeechSynthesisOutcome
	var stream *AudioDataStream
	var err error
	select {
	case textResult = <-textResultFuture:
		defer textResult.Close()
		checkSynthesisResult(t, textResult.Result, common.SynthesizingAudioCompleted)
		stream, err = NewAudioDataStreamFromSpeechSynthesisResult(textResult.Result)
		if err != nil {
			t.Error("crate audio data stream failed")
		}
		defer stream.Close()
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}

	var status common.StreamStatus
	status, err = stream.GetStatus()
	if err != nil {
		t.Error("stream status")
	}
	if status != common.StreamStatusAllData {
		t.Error("stream status")
	}

	audioData1 := textResult.Result.AudioData
	audioData2 := make([]byte, len(audioData1))
	stream.Read(audioData2)
	if !bytes.Equal(audioData1, audioData2) {
		t.Error("audio data in result and audio data stream are not equal.")
	}
	off, e := stream.GetOffset()
	if e != nil {
		t.Error("audio data stream get offset error: ", e)
	}
	if off != len(audioData1) {
		t.Error("audio data stream get offset incorrect.")
	}
	// set offset to 0 and read again
	e = stream.SetOffset(0)
	if e != nil {
		t.Error("audio data stream set offset error: ", e)
	}
	audioData3 := make([]byte, len(audioData1))
	stream.Read(audioData3)
	if !bytes.Equal(audioData2, audioData3) {
		t.Error("audio data is not equal.")
	}

	saveOutcome := stream.SaveToWavFileAsync("tmp_synthesis.mp3")
	select {
	case err = <-saveOutcome:
		if err != nil {
			t.Error("audio data stream save to file failed")
		}
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}

	file, _ := os.Open("tmp_synthesis.mp3")
	defer file.Close()
	audioData4 := make([]byte, len(audioData1))
	file.Read(audioData4)
	if !bytes.Equal(audioData2, audioData4) {
		t.Error("audio data is not equal.")
	}
}

func TestSynthesisWithInvalidVoice(t *testing.T) {
	config := createSpeechConfig(t)
	config.SetSpeechSynthesisVoiceName("invalid")
	defer config.Close()
	synthesizer := createSpeechSynthesizerFromSpeechConfigAndAudioConfig(t, config, nil)
	defer synthesizer.Close()

	textResultFuture := synthesizer.SpeakTextAsync("text")
	var textResult SpeechSynthesisOutcome
	select {
	case textResult = <-textResultFuture:
		defer textResult.Close()
		checkSynthesisResult(t, textResult.Result, common.Canceled)
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}

	details, err := NewCancellationDetailsFromSpeechSynthesisResult(textResult.Result)
	if err != nil {
		t.Error("create cancellation details from synthesis result failed.")
	}
	if details.Reason != common.Error {
		t.Error("reason")
	}
	if details.ErrorCode != common.BadRequest {
		t.Error("error code")
	}
	if !strings.Contains(details.ErrorDetails, "invalid") {
		t.Error("error details")
	}
}

func TestSynthesisToPullAudioOutputStream(t *testing.T) {
	stream, err := audio.CreatePullAudioOutputStream()
	if err != nil {
		t.Error("create pull audio output stream error: ", err)
	}
	defer stream.Close()
	var audioConfig *audio.AudioConfig
	audioConfig, err = audio.NewAudioConfigFromStreamOutput(stream)
	if err != nil {
		t.Error("new audio config from stream output error: ", err)
	}
	defer audioConfig.Close()
	synthesizer := createSpeechSynthesizerFromAudioConfig(t, audioConfig)

	textResultFuture := synthesizer.SpeakTextAsync("text")
	select {
	case textResult := <-textResultFuture:
		defer textResult.Close()
		// checkSynthesisResult(t, textResult.Result, common.SynthesizingAudioStarted)
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}

	synthesizer.Close()
	var bytes []byte
	for {
		buf, err := stream.Read(3200)
		if err != nil || len(buf) == 0 {
			break
		}
		bytes = append(bytes, buf...)
	}

	if len(bytes) == 0 {
		t.Error("error reading data from pull audio output stream.")
	}
}

// viseme received
func TestSynthesizerVisemeEvents(t *testing.T) {
	synthesizer := createSpeechSynthesizerFromAudioConfig(t, nil)
	defer synthesizer.Close()
	visemeReceivedFuture := make(chan string)
	synthesizer.VisemeReceived(func(event SpeechSynthesisVisemeEventArgs) {
		defer event.Close()
		t.Logf("viseme received event, audio offset [%d], viseme ID [%d]", event.AudioOffset, event.VisemeID)
		if event.AudioOffset <= 0 {
			t.Error("viseme received audio offset")
		}
		select {
		case visemeReceivedFuture <- "visemeReceivedFuture":
		default:
		}
	})
	resultFuture := synthesizer.SpeakSsmlAsync("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xmlns:mstts='http://www.w3.org/2001/mstts' xmlns:emo='http://www.w3.org/2009/10/emotionml' xml:lang='en-US'><voice name='en-US-AriaNeural'><mstts:viseme type='redlips_front'>yet</mstts:viseme></voice></speak>")
	select {
	case <-visemeReceivedFuture:
	case <-time.After(timeout):
		t.Error("Timeout waiting for VisemeReceived event.")
	}
	select {
	case result := <-resultFuture:
		defer result.Close()
		checkSynthesisResult(t, result.Result, common.SynthesizingAudioCompleted)
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}
}

// word boundary and bookmark reached
func TestSynthesizerEvents2(t *testing.T) {
	synthesizer := createSpeechSynthesizerFromAudioConfig(t, nil)
	defer synthesizer.Close()
	wordBoundaryFuture := make(chan bool)
	synthesizer.WordBoundary(func(event SpeechSynthesisWordBoundaryEventArgs) {
		defer event.Close()
		t.Logf("word boundary event, audio offset [%d], text offset [%d], word length [%d]", event.AudioOffset, event.TextOffset, event.WordLength)
		if event.AudioOffset <= 0 {
			t.Error("word boundary audio offset")
		}
		if event.Duration <= 0 {
			t.Error("word boundary duration")
		}
		if event.TextOffset <= 0 {
			t.Error("word boundary text offset")
		}
		if event.WordLength <= 0 {
			t.Error("word boundary word length")
		}
		select {
		case wordBoundaryFuture <- true:
		default:
		}
	})
	bookmarkReachedFuture := make(chan string)
	synthesizer.BookmarkReached(func(event SpeechSynthesisBookmarkEventArgs) {
		defer event.Close()
		t.Logf("Bookmark reached event, audio offset [%d], text [%s]", event.AudioOffset, event.Text)
		if event.AudioOffset <= 0 {
			t.Error("bookmark audio offset error")
		}
		if event.Text != "mark" {
			t.Error("bookmark text error")
		}
		bookmarkReachedFuture <- "bookmarkReachedFuture"
	})
	resultFuture := synthesizer.SpeakSsmlAsync("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xmlns:mstts='http://www.w3.org/2001/mstts' xmlns:emo='http://www.w3.org/2009/10/emotionml' xml:lang='en-US'><voice name='en-US-AriaNeural'>hello<bookmark mark='mark'/></voice></speak>")
	select {
	case <-bookmarkReachedFuture:
		t.Logf("Received BookmarkReached event.")
	case <-time.After(timeout):
		t.Error("Timeout waiting for BookmarkReached event.")
	}
	select {
	case <-wordBoundaryFuture:
	case <-time.After(timeout):
		t.Error("Timeout waiting for WordBoundary event.")
	}
	select {
	case result := <-resultFuture:
		defer result.Close()
		checkSynthesisResult(t, result.Result, common.SynthesizingAudioCompleted)
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}
}

func TestSynthesisGetAvailableVoices(t *testing.T) {
	synthesizer := createSpeechSynthesizerFromAudioConfig(t, nil)
	defer synthesizer.Close()
	resultFuture := synthesizer.GetVoicesAsync("en-US")
	select {
	case outcome := <-resultFuture:
		defer outcome.Close()
		if outcome.Result.Reason != common.VoicesListRetrieved {
			t.Error("synthesizer get voices failed")
		}
		if len(outcome.Result.Voices) == 0 {
			t.Error("no voice")
		}
		for _, e := range outcome.Result.Voices {
			if len(e.Name) == 0 {
				t.Error("voice name error")
			}
		}
		voice := outcome.Result.Voices[0]
		if voice.VoiceType != common.OnlineNeural {
			t.Error("Voice type is incorrect.")
		}
		if voice.Gender != common.Female {
			t.Error("Voice gender error.")
		}
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis voices result.")
	}
}

func TestSynthesisWithLanguageAutoDetection(t *testing.T) {
	config := createSpeechConfig(t)
	defer config.Close()
	languageConfig, err := NewAutoDetectSourceLanguageConfigFromOpenRange()
	if err != nil {
		t.Error("Got an error: ", err)
	}
	defer languageConfig.Close()
	synthesizer, err := NewSpeechSynthesizerFomAutoDetectSourceLangConfig(config, languageConfig, nil)
	if err != nil {
		t.Error("Got an error: ", err)
	}
	if synthesizer == nil {
		t.Error("synthesizer creation failed")
	}
	defer synthesizer.Close()
	textResultFuture := synthesizer.SpeakTextAsync("你好，世界。")

	var textResult SpeechSynthesisOutcome
	select {
	case textResult = <-textResultFuture:
		defer textResult.Close()
		checkSynthesisResult(t, textResult.Result, common.SynthesizingAudioCompleted)
		if len(textResult.Result.AudioData) < 32000 {
			t.Error("audio should longer than 1s.")
		}
	case <-time.After(timeout):
		t.Error("Timeout waiting for synthesis result.")
	}
}
