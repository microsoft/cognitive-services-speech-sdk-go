// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package speech

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
)

// Inline commit tests. The Go bindings only add API surface on top of the
// native SDK (PushAudioInputStream.Commit/CommitChannel and the result
// CommitToken field), so these tests verify that each API works for each
// recognizer type rather than repeating the native SDK test suite.

const inlineCommitAckTimeout = 20 * time.Second

// enableInlineCommit sets the service flight flag needed for inline commit support.
func enableInlineCommit(t *testing.T, config *SpeechConfig) {
	err := config.SetServiceProperty("setfeature", "forcecommit", common.URIQueryParameter)
	if err != nil {
		t.Fatal("Got an error setting service property: ", err)
	}
}

// readWavData returns the format and audio samples (data chunk) of a WAV file.
func readWavData(t *testing.T, filename string) (sampleRate uint32, bitsPerSample uint8, channels uint8, data []byte) {
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal("Error reading file: ", err)
	}
	if len(content) < 12 || !bytes.Equal(content[0:4], []byte("RIFF")) || !bytes.Equal(content[8:12], []byte("WAVE")) {
		t.Fatal("Not a WAV file: ", filename)
	}
	pos := 12
	for pos+8 <= len(content) {
		id := string(content[pos : pos+4])
		size := int(binary.LittleEndian.Uint32(content[pos+4 : pos+8]))
		body := pos + 8
		if body+size > len(content) {
			size = len(content) - body
		}
		switch id {
		case "fmt ":
			channels = uint8(binary.LittleEndian.Uint16(content[body+2 : body+4]))
			sampleRate = binary.LittleEndian.Uint32(content[body+4 : body+8])
			bitsPerSample = uint8(binary.LittleEndian.Uint16(content[body+14 : body+16]))
		case "data":
			return sampleRate, bitsPerSample, channels, content[body : body+size]
		}
		pos = body + size + size%2
	}
	t.Fatal("No data chunk in WAV file: ", filename)
	return
}

// createPushStreamForWav creates a push stream with the format of the given WAV file
// and returns it together with the audio data to write.
func createPushStreamForWav(t *testing.T, filename string) (*audio.PushAudioInputStream, []byte) {
	sampleRate, bitsPerSample, channels, data := readWavData(t, filename)
	format, err := audio.GetWaveFormatPCM(sampleRate, bitsPerSample, channels)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer format.Close()
	stream, err := audio.CreatePushAudioInputStreamFromFormat(format)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	return stream, data
}

func writeInChunks(t *testing.T, stream *audio.PushAudioInputStream, data []byte, commitAtBytes int,
	paced bool, commitHere func()) {
	const chunkSize = 3200
	committed := false
	for start := 0; start < len(data); start += chunkSize {
		end := start + chunkSize
		if end > len(data) {
			end = len(data)
		}
		if err := stream.Write(data[start:end]); err != nil {
			t.Fatal("Error writing to the stream: ", err)
		}
		if !committed && commitAtBytes >= 0 && end >= commitAtBytes {
			commitHere()
			committed = true
		}
		if paced {
			time.Sleep(100 * time.Millisecond)
		}
	}
	if !committed {
		commitHere()
	}
}

// commitResult holds the fields of a recognized result that the tests check.
type commitResult struct {
	commitToken uint32
	channel     uint32
	text        string
}

// commitAndWait writes the audio, commits after commitAtBytes bytes (at the end
// if negative), checks the rate limit and waits for a result carrying the issued
// token. If paced, audio is written at roughly real-time pace for 16-bit audio at
// 16 kHz mono or 8 kHz stereo. The stream is closed after that.
func commitAndWait(t *testing.T, stream *audio.PushAudioInputStream, data []byte, commitAtBytes int, paced bool,
	results <-chan commitResult, commit func() (uint32, error)) commitResult {
	var token uint32
	writeInChunks(t, stream, data, commitAtBytes, paced, func() {
		var err error
		token, err = commit()
		if err != nil {
			t.Fatal("Commit returned an error: ", err)
		}
		if token == 0 {
			t.Fatal("Expected a non-zero commit token")
		}
		t.Logf("Commit token: %d", token)

		// A second commit immediately afterwards is rejected by the 100 ms rate limit.
		rejected, err := commit()
		if err != nil {
			t.Fatal("Commit returned an error: ", err)
		}
		if rejected != 0 {
			t.Errorf("Expected commit within 100 ms to be rejected with 0, got %d", rejected)
		}
	})

	timeout := time.After(inlineCommitAckTimeout)
	for {
		select {
		case r := <-results:
			t.Logf("Recognized (channel %d, commit token %d): %s", r.channel, r.commitToken, r.text)
			if r.commitToken == 0 {
				continue
			}
			if r.commitToken != token {
				t.Fatalf("Expected commit token %d, got %d", token, r.commitToken)
			}
			stream.CloseStream()
			return r
		case <-timeout:
			stream.CloseStream()
			t.Fatal("Timeout waiting for commit acknowledgment")
			return commitResult{}
		}
	}
}

func TestInlineCommitSpeechRecognizer(t *testing.T) {
	config, err := NewSpeechConfigFromSubscription(os.Getenv("SPEECH_SUBSCRIPTION_KEY"), os.Getenv("SPEECH_SUBSCRIPTION_REGION"))
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer config.Close()
	enableInlineCommit(t, config)

	stream, data := createPushStreamForWav(t, "../test_files/TalkForAFewSeconds16.wav")
	defer stream.Close()
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer audioConfig.Close()
	recognizer, err := NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer recognizer.Close()

	results := make(chan commitResult, 16)
	recognizer.Recognized(func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		results <- commitResult{event.Result.CommitToken, event.Result.Channel, event.Result.Text}
	})
	recognizer.Canceled(func(event SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Logf("Canceled: reason %d, details: %s", event.Reason, event.ErrorDetails)
		if event.Reason == common.Error {
			t.Errorf("Recognition canceled with an error: %s", event.ErrorDetails)
		}
	})

	if err := <-recognizer.StartContinuousRecognitionAsync(); err != nil {
		t.Fatal("Got an error: ", err)
	}
	r := commitAndWait(t, stream, data, 90000, false, results, stream.Commit)
	if r.text == "" {
		t.Error("Expected recognized text in the commit acknowledgment")
	}
	if err := <-recognizer.StopContinuousRecognitionAsync(); err != nil {
		t.Error("Got an error: ", err)
	}
}

func TestInlineCommitChannelScopedMultichannel(t *testing.T) {
	config, err := NewSpeechConfigFromSubscription(os.Getenv("SPEECH_SUBSCRIPTION_KEY"), os.Getenv("SPEECH_SUBSCRIPTION_REGION"))
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer config.Close()
	enableInlineCommit(t, config)
	if err := config.SetProperty(common.EnableMultiChannelProcessing, "true"); err != nil {
		t.Fatal("Got an error: ", err)
	}

	stream, data := createPushStreamForWav(t, "../test_files/whatstheweatherlike_8khz_2ch.wav")
	defer stream.Close()
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer audioConfig.Close()
	recognizer, err := NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer recognizer.Close()

	results := make(chan commitResult, 16)
	recognizer.Recognized(func(event SpeechRecognitionEventArgs) {
		defer event.Close()
		results <- commitResult{event.Result.CommitToken, event.Result.Channel, event.Result.Text}
	})
	recognizer.Canceled(func(event SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Logf("Canceled: reason %d, details: %s", event.Reason, event.ErrorDetails)
		if event.Reason == common.Error {
			t.Errorf("Recognition canceled with an error: %s", event.ErrorDetails)
		}
	})

	if err := <-recognizer.StartContinuousRecognitionAsync(); err != nil {
		t.Fatal("Got an error: ", err)
	}
	const committedChannel = 1
	r := commitAndWait(t, stream, data, 85000, true, results, func() (uint32, error) { return stream.CommitChannel(committedChannel) })
	if r.channel != committedChannel {
		t.Errorf("Expected acknowledgment on channel %d, got %d", committedChannel, r.channel)
	}
	if err := <-recognizer.StopContinuousRecognitionAsync(); err != nil {
		t.Error("Got an error: ", err)
	}
}

func TestInlineCommitTranslationRecognizer(t *testing.T) {
	config, err := NewSpeechTranslationConfigFromSubscription(os.Getenv("SPEECH_SUBSCRIPTION_KEY"), os.Getenv("SPEECH_SUBSCRIPTION_REGION"))
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer config.Close()
	enableInlineCommit(t, &config.SpeechConfig)
	if err := config.SetSpeechRecognitionLanguage("en-us"); err != nil {
		t.Fatal("Got an error: ", err)
	}
	if err := config.AddTargetLanguage("es"); err != nil {
		t.Fatal("Got an error: ", err)
	}
	if err := config.AddTargetLanguage("fr"); err != nil {
		t.Fatal("Got an error: ", err)
	}

	stream, data := createPushStreamForWav(t, "../test_files/TalkForAFewSeconds16.wav")
	defer stream.Close()
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer audioConfig.Close()
	recognizer, err := NewTranslationRecognizerFromConfig(config, audioConfig)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer recognizer.Close()

	results := make(chan commitResult, 16)
	recognizer.Recognized(func(event TranslationRecognitionEventArgs) {
		defer event.Close()
		t.Logf("Translations: %v", event.Result.GetTranslations())
		results <- commitResult{event.Result.CommitToken, event.Result.Channel, event.Result.Text}
	})
	recognizer.Canceled(func(event TranslationRecognitionCanceledEventArgs) {
		defer event.Close()
		t.Logf("Canceled: reason %d, details: %s", event.Reason, event.ErrorDetails)
		if event.Reason == common.Error {
			t.Errorf("Recognition canceled with an error: %s", event.ErrorDetails)
		}
	})

	if err := <-recognizer.StartContinuousRecognitionAsync(); err != nil {
		t.Fatal("Got an error: ", err)
	}
	r := commitAndWait(t, stream, data, 90000, false, results, stream.Commit)
	if r.text == "" {
		t.Error("Expected recognized text in the commit acknowledgment")
	}
	if err := <-recognizer.StopContinuousRecognitionAsync(); err != nil {
		t.Error("Got an error: ", err)
	}
}

func TestInlineCommitConversationTranscriber(t *testing.T) {
	config, err := NewSpeechConfigFromSubscription(os.Getenv("SPEECH_SUBSCRIPTION_KEY"), os.Getenv("SPEECH_SUBSCRIPTION_REGION"))
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer config.Close()
	enableInlineCommit(t, config)

	stream, data := createPushStreamForWav(t, "../test_files/TalkForAFewSeconds16.wav")
	defer stream.Close()
	audioConfig, err := audio.NewAudioConfigFromStreamInput(stream)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer audioConfig.Close()
	transcriber, err := NewConversationTranscriberFromConfig(config, audioConfig)
	if err != nil {
		t.Fatal("Got an error: ", err)
	}
	defer transcriber.Close()

	results := make(chan commitResult, 16)
	transcriber.Transcribed(func(event ConversationTranscriptionEventArgs) {
		defer event.Close()
		results <- commitResult{event.Result.CommitToken, event.Result.Channel, event.Result.Text}
	})
	transcriber.Canceled(func(event ConversationTranscriptionCanceledEventArgs) {
		defer event.Close()
		t.Logf("Canceled: reason %d, details: %s", event.Reason, event.ErrorDetails)
		if event.Reason == common.Error {
			t.Errorf("Recognition canceled with an error: %s", event.ErrorDetails)
		}
	})

	if err := <-transcriber.StartTranscribingAsync(); err != nil {
		t.Fatal("Got an error: ", err)
	}
	r := commitAndWait(t, stream, data, 90000, false, results, stream.Commit)
	if r.text == "" {
		t.Error("Expected transcribed text in the commit acknowledgment")
	}
	if err := <-transcriber.StopTranscribingAsync(); err != nil {
		t.Error("Got an error: ", err)
	}
}
