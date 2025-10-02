// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

// Package main demonstrates usages for the speech recognizer and dialog service connector
package main

import (
	"fmt"
	"os"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/samples/dialog_service_connector"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/samples/recognizer"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/samples/synthesizer"
)

type functionMap = map[string]func(string, string, string)

func printHelp(executableName string, samples functionMap) {
	fmt.Println("Input not valid")
	fmt.Println("Usage: ")
	fmt.Println(executableName, " <subscription> <region> <file> <sample>")
	fmt.Println("Where sample is of the format <scenario>:<sample>")
	fmt.Println("Available samples:")
	for id, _ := range samples {
		fmt.Println(" -- ", id)
	}
}

func main() {
	samples := functionMap{
		"speech_recognizer:RecognizeOnceFromWavFile":        recognizer.RecognizeOnceFromWavFile,
		"speech_recognizer:RecognizeOnceFromCompressedFile": recognizer.RecognizeOnceFromCompressedFile,
		"speech_recognizer:RecognizeOnceFromALAWFile":       recognizer.RecognizeOnceFromALAWFile,
		"speech_recognizer:ContinuousFromMicrophone":        recognizer.ContinuousFromMicrophone,
		"speech_recognizer:RecognizeContinuousUsingWrapper": recognizer.RecognizeContinuousUsingWrapper,
		"dialog_service_connector:ListenOnce":               dialog_service_connector.ListenOnce,
		"dialog_service_connector:KWS":                      dialog_service_connector.KWS,
		"dialog_service_connector:ListenOnceFromStream":     dialog_service_connector.ListenOnceFromStream,
		"speech_synthesizer:SynthesisToSpeaker":             synthesizer.SynthesisToSpeaker,
		"speech_synthesizer:SynthesisToAudioDataStream":     synthesizer.SynthesisToAudioDataStream,
	}
	args := os.Args[1:]
	if len(args) != 4 {
		printHelp(os.Args[0], samples)
		return
	}
	subscription := args[0]
	region := args[1]
	file := args[2]
	sample := args[3]
	sampleFunction := samples[sample]
	if sampleFunction == nil {
		printHelp(os.Args[0], samples)
		return
	}
	sampleFunction(subscription, region, file)
}
