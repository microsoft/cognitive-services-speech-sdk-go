// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// SpeechSynthesisOutputFormat defines the possible speech synthesis output audio formats.
type SpeechSynthesisOutputFormat int

const (
	// Raw8Khz8BitMonoMULaw stands for raw-8khz-8bit-mono-mulaw
	Raw8Khz8BitMonoMULaw SpeechSynthesisOutputFormat = 1

	// Riff16Khz16KbpsMonoSiren stands for riff-16khz-16kbps-mono-siren
	// Unsupported by the service. Do not use this value.
	Riff16Khz16KbpsMonoSiren SpeechSynthesisOutputFormat = 2

	// Audio16Khz16KbpsMonoSiren stands for audio-16khz-16kbps-mono-siren
	// Unsupported by the service. Do not use this value.
	Audio16Khz16KbpsMonoSiren SpeechSynthesisOutputFormat = 3

	// Audio16Khz32KBitRateMonoMp3 stands for audio-16khz-32kbitrate-mono-mp3
	Audio16Khz32KBitRateMonoMp3 SpeechSynthesisOutputFormat = 4

	// Audio16Khz128KBitRateMonoMp3 stands for audio-16khz-128kbitrate-mono-mp3
	Audio16Khz128KBitRateMonoMp3 SpeechSynthesisOutputFormat = 5

	// Audio16Khz64KBitRateMonoMp3 stands for audio-16khz-64kbitrate-mono-mp3
	Audio16Khz64KBitRateMonoMp3 SpeechSynthesisOutputFormat = 6

	// Audio24Khz48KBitRateMonoMp3 stands for audio-24khz-48kbitrate-mono-mp3
	Audio24Khz48KBitRateMonoMp3 SpeechSynthesisOutputFormat = 7

	// Audio24Khz96KBitRateMonoMp3 stands for audio-24khz-96kbitrate-mono-mp3
	Audio24Khz96KBitRateMonoMp3 SpeechSynthesisOutputFormat = 8

	// Audio24Khz160KBitRateMonoMp3 stands for audio-24khz-160kbitrate-mono-mp3
	Audio24Khz160KBitRateMonoMp3 SpeechSynthesisOutputFormat = 9

	// Raw16Khz16BitMonoTrueSilk stands for raw-16khz-16bit-mono-truesilk
	Raw16Khz16BitMonoTrueSilk SpeechSynthesisOutputFormat = 10

	// Riff16Khz16BitMonoPcm stands for riff-16khz-16bit-mono-pcm
	Riff16Khz16BitMonoPcm SpeechSynthesisOutputFormat = 11

	// Riff8Khz16BitMonoPcm stands for riff-8khz-16bit-mono-pcm
	Riff8Khz16BitMonoPcm SpeechSynthesisOutputFormat = 12

	// Riff24Khz16BitMonoPcm stands for riff-24khz-16bit-mono-pcm
	Riff24Khz16BitMonoPcm SpeechSynthesisOutputFormat = 13

	// Riff8Khz8BitMonoMULaw stands for riff-8khz-8bit-mono-mulaw
	Riff8Khz8BitMonoMULaw SpeechSynthesisOutputFormat = 14

	// Raw16Khz16BitMonoPcm stands for raw-16khz-16bit-mono-pcm
	Raw16Khz16BitMonoPcm SpeechSynthesisOutputFormat = 15

	// Raw24Khz16BitMonoPcm stands for raw-24khz-16bit-mono-pcm
	Raw24Khz16BitMonoPcm SpeechSynthesisOutputFormat = 16

	// Raw8Khz16BitMonoPcm stands for raw-8khz-16bit-mono-pcm
	Raw8Khz16BitMonoPcm SpeechSynthesisOutputFormat = 17

	// Ogg16Khz16BitMonoOpus stands for ogg-16khz-16bit-mono-opus
	Ogg16Khz16BitMonoOpus SpeechSynthesisOutputFormat = 18

	// Ogg24Khz16BitMonoOpus stands for ogg-24khz-16bit-mono-opus
	Ogg24Khz16BitMonoOpus SpeechSynthesisOutputFormat = 19

	// Raw48Khz16BitMonoPcm stands for raw-48khz-16bit-mono-pcm
	Raw48Khz16BitMonoPcm SpeechSynthesisOutputFormat = 20

	// Riff48Khz16BitMonoPcm stands for riff-48khz-16bit-mono-pcm
	Riff48Khz16BitMonoPcm SpeechSynthesisOutputFormat = 21

	// Audio48Khz96KBitRateMonoMp3 stands for audio-48khz-96kbitrate-mono-mp3
	Audio48Khz96KBitRateMonoMp3 SpeechSynthesisOutputFormat = 22

	// Audio48Khz192KBitRateMonoMp3 stands for audio-48khz-192kbitrate-mono-mp3
	Audio48Khz192KBitRateMonoMp3 SpeechSynthesisOutputFormat = 23

	// Ogg48Khz16BitMonoOpus stands for ogg-48khz-16bit-mono-opus
	Ogg48Khz16BitMonoOpus SpeechSynthesisOutputFormat = 24

	// Webm16Khz16BitMonoOpus stands for webm-16khz-16bit-mono-opus
	Webm16Khz16BitMonoOpus SpeechSynthesisOutputFormat = 25

	// Webm24Khz16BitMonoOpus stands for webm-24khz-16bit-mono-opus
	Webm24Khz16BitMonoOpus SpeechSynthesisOutputFormat = 26

	// Raw24Khz16BitMonoTrueSilk stands for raw-24khz-16bit-mono-truesilk
	Raw24Khz16BitMonoTrueSilk SpeechSynthesisOutputFormat = 27

	// Raw8Khz8BitMonoALaw stands for raw-8khz-8bit-mono-alaw
	Raw8Khz8BitMonoALaw SpeechSynthesisOutputFormat = 28

	// Riff8Khz8BitMonoALaw stands for riff-8khz-8bit-mono-alaw
	Riff8Khz8BitMonoALaw SpeechSynthesisOutputFormat = 29

	// Webm24Khz16Bit24KbpsMonoOpus stands for webm-24khz-16bit-24kbps-mono-opus
	// Audio compressed by OPUS codec in a WebM container, with bitrate of 24kbps, optimized for IoT scenario.
	Webm24Khz16Bit24KbpsMonoOpus SpeechSynthesisOutputFormat = 30

	// Audio16Khz16Bit32KbpsMonoOpus stands for audio-16khz-16bit-32kbps-mono-opus
	// Audio compressed by OPUS codec without container, with bitrate of 32kbps.
	Audio16Khz16Bit32KbpsMonoOpus SpeechSynthesisOutputFormat = 31

	// Audio24Khz16Bit48KbpsMonoOpus stands for audio-24khz-16bit-48kbps-mono-opus
	// Audio compressed by OPUS codec without container, with bitrate of 48kbps.
	Audio24Khz16Bit48KbpsMonoOpus SpeechSynthesisOutputFormat = 32

	// Audio24Khz16Bit24KbpsMonoOpus stands for audio-24khz-16bit-24kbps-mono-opus
	// Audio compressed by OPUS codec without container, with bitrate of 24kbps.
	Audio24Khz16Bit24KbpsMonoOpus SpeechSynthesisOutputFormat = 33

	// Raw22050Hz16BitMonoPcm stands for raw-22050hz-16bit-mono-pcm
	// Raw PCM audio at 22050Hz sampling rate and 16-bit depth.
	Raw22050Hz16BitMonoPcm SpeechSynthesisOutputFormat = 34

	// Riff22050Hz16BitMonoPcm stands for riff-22050hz-16bit-mono-pcm
	// PCM audio at 22050Hz sampling rate and 16-bit depth, with RIFF header.
	Riff22050Hz16BitMonoPcm SpeechSynthesisOutputFormat = 35

	// Raw44100Hz16BitMonoPcm stands for raw-44100hz-16bit-mono-pcm
	// Raw PCM audio at 44100Hz sampling rate and 16-bit depth.
	Raw44100Hz16BitMonoPcm SpeechSynthesisOutputFormat = 36

	// Riff44100Hz16BitMonoPcm stands for riff-44100hz-16bit-mono-pcm
	// PCM audio at 44100Hz sampling rate and 16-bit depth, with RIFF header.
	Riff44100Hz16BitMonoPcm SpeechSynthesisOutputFormat = 37

	// AmrWb16000Hz stands for amr-wb-16000hz
	// AMR-WB audio at 16kHz sampling rate.
	AmrWb16000Hz SpeechSynthesisOutputFormat = 38
)
