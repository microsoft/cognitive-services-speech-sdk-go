package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func main() {
	// 1. Configuration
	key := os.Getenv("AZURE_TTS_API_KEY")
	region := os.Getenv("AZURE_TTS_REGION")
	if key == "" || region == "" {
		log.Println("Please set AZURE_TTS_API_KEY and AZURE_TTS_REGION environment variables.")
		return
	}

	// MUST use WebSocket v2 Endpoint
	endpoint := fmt.Sprintf("wss://%s.tts.speech.microsoft.com/cognitiveservices/websocket/v2", region)
	config, err := speech.NewSpeechConfigFromEndpointWithSubscription(endpoint, key)
	if err != nil {
		log.Printf("Failed to create config: %v\n", err)
		return
	}
	defer config.Close()

	// Explicitly set output format (MP3 24kHz)
	config.SetSpeechSynthesisOutputFormat(common.Audio24Khz48KBitRateMonoMp3)
	// Set timeouts to prevent server disconnection
	config.SetProperty(common.PropertyID(14101), "100000000") // FrameTimeoutInterval
	config.SetProperty(common.PropertyID(14102), "10")        // RtfTimeoutThreshold

	// 2. Create Synthesizer
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(config, nil)
	if err != nil {
		log.Printf("Failed to create synthesizer: %v\n", err)
		return
	}
	defer synthesizer.Close()

	// Create output file
	outputFile, err := os.Create("output.mp3")
	if err != nil {
		log.Printf("Failed to create output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	var startTime time.Time
	var firstByteReceived bool

	// 3. Bind callbacks (Receive audio)
	synthesizer.Synthesizing(func(event speech.SpeechSynthesisEventArgs) {
		defer event.Close()
		data := event.Result.AudioData
		if len(data) > 0 {
			if !firstByteReceived {
				latency := time.Since(startTime)
				log.Printf("First byte received. Latency (TTFB): %v\n", latency)
				firstByteReceived = true
			}
			log.Printf("Received audio chunk: %d bytes\n", len(data))
			if _, err := outputFile.Write(data); err != nil {
				log.Printf("Failed to write audio data: %v\n", err)
			}
		}
	})

	synthesizer.SynthesisCanceled(func(event speech.SpeechSynthesisEventArgs) {
		defer event.Close()
		details, _ := speech.NewCancellationDetailsFromSpeechSynthesisResult(&event.Result)
		log.Printf("CANCELED: Reason=%d, ErrorDetails=%s\n", details.Reason, details.ErrorDetails)
	})

	// 4. Create TextStream Request
	req, err := speech.NewSpeechSynthesisRequest(speech.SpeechSynthesisRequestInputType_TextStream)
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		return
	}
	defer req.Close()

	// Set Voice
	req.SetVoice("en-US-JennyNeural")

	// 5. Start Async Synthesis
	log.Println("Starting synthesis...")
	startTime = time.Now()
	outcomeChan := synthesizer.SpeakRequestAsync(req)

	// 6. Stream Text
	stream := req.InputStream()
	text := "Here is a longer text to simulate LLM streaming output. " +
		"We want to verify that the audio is generated in real-time as we send text chunks. " +
		"The latency should be low, and the audio should be smooth. " +
		"Azure TTS supports text streaming via WebSocket v2, which allows us to send text piece by piece. " +
		"This is crucial for conversational AI applications where we want to start speaking as soon as the first sentence is generated. " +
		"Let's see how it performs with this longer paragraph."
	words := strings.Split(text, " ")

	for _, word := range words {
		chunk := word + " "
		log.Printf("Sending: %s\n", chunk)
		if err := stream.Write(chunk); err != nil {
			log.Printf("Write failed: %v\n", err)
			return
		}
		// Simulate LLM generation latency
		time.Sleep(100 * time.Millisecond)
	}

	// 7. End Input
	log.Println("Closing input stream...")
	if err := stream.Close(); err != nil {
		log.Printf("Close failed: %v\n", err)
		return
	}

	// 8. Wait for Completion
	outcome := <-outcomeChan
	if outcome.Error != nil {
		log.Printf("Synthesis error: %v\n", outcome.Error)
	} else if outcome.Result.Reason == common.Canceled {
		log.Println("Synthesis was canceled.")
	} else {
		log.Println("Synthesis completed successfully. Audio saved to output.mp3")
	}
}
