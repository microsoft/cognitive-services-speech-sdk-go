// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package common

// PropertyID defines speech property ids.
// Changed in version 1.4.0.
type PropertyID int

const (
	// SpeechServiceConnectionKey is the Cognitive Services Speech Service subscription key. Under normal
	// circumstances, you shouldn't have to use this property directly.
	// Instead, use NewSpeechConfigFromSubscription.
	SpeechServiceConnectionKey PropertyID = 1000

	// SpeechServiceConnectionEndpoint is the Cognitive Services Speech Service endpoint (url).
	// Under normal circumstances, you shouldn't have to use this property directly.
	// Instead, use NewSpeechConfigFromEndpoint.
	// NOTE: This endpoint is not the same as the endpoint used to obtain an access token.
	SpeechServiceConnectionEndpoint PropertyID = 1001

	// SpeechServiceConnectionRegion is the Cognitive Services Speech Service region. Under normal circumstances,
	// you shouldn't have to use this property directly.
	// Instead, use NewSpeechConfigFromSubscription, NewSpeechConfigFromEndpoint, NewSpeechConfigFromHost,
	// NewSpeechConfigFromAuthorizationToken.
	SpeechServiceConnectionRegion PropertyID = 1002

	// SpeechServiceAuthorizationToken is the Cognitive Services Speech Service authorization token (aka access token).
	// Under normal circumstances, you shouldn't have to use this property directly.
	// Instead, use NewSpeechConfigFromAuthorizationToken,
	// Recognizer.SetAuthorizationToken
	SpeechServiceAuthorizationToken PropertyID = 1003

	// SpeechServiceAuthorizationType is the Cognitive Services Speech Service authorization type. Currently unused.
	SpeechServiceAuthorizationType PropertyID = 1004

	// SpeechServiceConnectionEndpointID is the Cognitive Services Custom Speech Service endpoint id. Under normal
	// circumstances, you shouldn't have to use this property directly.
	// Instead use SpeechConfig.SetEndpointId.
	// NOTE: The endpoint id is available in the Custom Speech Portal, listed under Endpoint Details.
	SpeechServiceConnectionEndpointID PropertyID = 1005

	// SpeechServiceConnectionHost is the Cognitive Services Speech Service host (url). Under normal circumstances,
	// you shouldn't have to use this property directly.
	// Instead, use NewSpeechConfigFromHost.
	SpeechServiceConnectionHost PropertyID = 1006

	// SpeechServiceConnectionProxyHostName is the host name of the proxy server used to connect to the Cognitive Services
	// Speech Service. Under normal circumstances, you shouldn't have to use this property directly.
	// Instead, use SpeechConfig.SetProxy.
	SpeechServiceConnectionProxyHostName PropertyID = 1100

	// SpeechServiceConnectionProxyPort is the port of the proxy server used to connect to the Cognitive Services Speech
	// Service. Under normal circumstances, you shouldn't have to use this property directly.
	// Instead, use SpeechConfig.SetProxy.
	SpeechServiceConnectionProxyPort PropertyID = 1101

	// SpeechServiceConnectionProxyUserName is the user name of the proxy server used to connect to the Cognitive Services
	// Speech Service. Under normal circumstances, you shouldn't have to use this property directly.
	// Instead, use SpeechConfig.SetProxy.
	SpeechServiceConnectionProxyUserName PropertyID = 1102

	// SpeechServiceConnectionProxyPassword is the password of the proxy server used to connect to the Cognitive Services
	// Speech Service. Under normal circumstances, you shouldn't have to use this property directly.
	// Instead, use SpeechConfig.SetProxy.
	SpeechServiceConnectionProxyPassword PropertyID = 1103

	// SpeechServiceConnectionURL is the URL string built from speech configuration. This property is intended to be read-only.
	// The SDK is using it internally.
	SpeechServiceConnectionURL PropertyID = 1104

	// Specifies the list of hosts for which proxies should not be used. This setting overrides all other configurations.
	// Hostnames are separated by commas and are matched in a case-insensitive manner. Wildcards are not supported.
	SpeechServiceConnectionProxyHostBypass PropertyID = 1105

	// SpeechServiceConnectionTranslationToLanguages is the list of comma separated languages used as target translation
	// languages. Under normal circumstances, you shouldn't have to use this property directly.
	// Instead use SpeechTranslationConfig.AddTargetLanguage and SpeechTranslationConfig.GetTargetLanguages.
	SpeechServiceConnectionTranslationToLanguages PropertyID = 2000

	// SpeechServiceConnectionTranslationVoice is the name of the Cognitive Service Text to Speech Service voice. Under normal
	// circumstances, you shouldn't have to use this property directly.
	// Instead use SpeechTranslationConfig.SetVoiceName.
	// NOTE: Valid voice names can be found at https://aka.ms/csspeech/voicenames.
	SpeechServiceConnectionTranslationVoice PropertyID = 2001

	// SpeechServiceConnectionTranslationFeatures is the translation features. For internal use.
	SpeechServiceConnectionTranslationFeatures PropertyID = 2002

	// SpeechServiceConnectionRecoMode is the Cognitive Services Speech Service recognition mode. Can be "INTERACTIVE",
	// "CONVERSATION" or "DICTATION".
	// This property is intended to be read-only. The SDK is using it internally.
	SpeechServiceConnectionRecoMode PropertyID = 3000

	// SpeechServiceConnectionRecoLanguage is the spoken language to be recognized (in BCP-47 format). Under normal
	// circumstances, you shouldn't have to use this property directly.
	// Instead, use SpeechConfig.SetSpeechRecognitionLanguage.
	SpeechServiceConnectionRecoLanguage PropertyID = 3001

	// SpeechSessionID is the session id. This id is a universally unique identifier (aka UUID) representing a specific
	// binding of an audio input stream and the underlying speech recognition instance to which it is bound. Under normal
	// circumstances, you shouldn't have to use this property directly.
	// Instead use SessionEventArgs.SessionId.
	SpeechSessionID PropertyID = 3002

	// SpeechServiceConnectionUserDefinedQueryParameters are the query parameters provided by users. They will be passed
	// to the service as URL query parameters.
	SpeechServiceConnectionUserDefinedQueryParameters PropertyID = 3003

	// SpeechServiceConnectionSynthLanguage is the spoken language to be synthesized (e.g. en-US)
	SpeechServiceConnectionSynthLanguage PropertyID = 3100

	// SpeechServiceConnectionSynthVoice is the name of the TTS voice to be used for speech synthesis
	SpeechServiceConnectionSynthVoice PropertyID = 3101

	// SpeechServiceConnectionSynthOutputFormat is the string to specify TTS output audio format.
	SpeechServiceConnectionSynthOutputFormat PropertyID = 3102

	// SpeechServiceConnectionSynthEnableCompressedAudioTransmission indicates if use compressed audio format
	// for speech synthesis audio transmission.
	// This property only affects when SpeechServiceConnectionSynthOutputFormat is set to a pcm format.
	// If this property is not set and GStreamer is available, SDK will use compressed format for synthesized audio transmission,
	// and decode it. You can set this property to "false" to use raw pcm format for transmission on wire.
	// Added in version 1.17.0
	SpeechServiceConnectionSynthEnableCompressedAudioTransmission PropertyID = 3103

	// SpeechServiceConnectionInitialSilenceTimeoutMs is the initial silence timeout value (in milliseconds) used by the
	// service.
	SpeechServiceConnectionInitialSilenceTimeoutMs PropertyID = 3200

	// This property is deprecated.
	// For current information about silence timeouts, please visit https://aka.ms/csspeech/timeouts.
	SpeechServiceConnectionEndSilenceTimeoutMs PropertyID = 3201

	// SpeechServiceConnectionEnableAudioLogging is a boolean value specifying whether audio logging is enabled in the service
	// or not.
	SpeechServiceConnectionEnableAudioLogging PropertyID = 3202

	// SpeechServiceConnectionAutoDetectSourceLanguages is the auto detect source languages.
	SpeechServiceConnectionAutoDetectSourceLanguages PropertyID = 3300

	// SpeechServiceConnectionAutoDetectSourceLanguageResult is the auto detect source language result.
	SpeechServiceConnectionAutoDetectSourceLanguageResult PropertyID = 3301

	// SpeechServiceResponseRequestDetailedResultTrueFalse the requested Cognitive Services Speech Service response output
	// format (simple or detailed). Under normal circumstances, you shouldn't have to use this property directly.
	// Instead use SpeechConfig.SetOutputFormat.
	SpeechServiceResponseRequestDetailedResultTrueFalse PropertyID = 4000

	// SpeechServiceResponseRequestProfanityFilterTrueFalse is the requested Cognitive Services Speech Service response
	// output profanity level. Currently unused.
	SpeechServiceResponseRequestProfanityFilterTrueFalse PropertyID = 4001

	// SpeechServiceResponseProfanityOption is the requested Cognitive Services Speech Service response output profanity
	// setting.
	// Allowed values are "masked", "removed", and "raw".
	SpeechServiceResponseProfanityOption PropertyID = 4002

	// SpeechServiceResponsePostProcessingOption a string value specifying which post processing option should be used
	// by the service.
	// Allowed values are "TrueText".
	SpeechServiceResponsePostProcessingOption PropertyID = 4003

	// SpeechServiceResponseRequestWordLevelTimestamps is a boolean value specifying whether to include word-level
	// timestamps in the response result.
	SpeechServiceResponseRequestWordLevelTimestamps PropertyID = 4004

	// SpeechServiceResponseStablePartialResultThreshold is the number of times a word has to be in partial results
	// to be returned.
	SpeechServiceResponseStablePartialResultThreshold PropertyID = 4005

	// SpeechServiceResponseOutputFormatOption is a string value specifying the output format option in the response
	// result. Internal use only.
	SpeechServiceResponseOutputFormatOption PropertyID = 4006

	// SpeechServiceResponseTranslationRequestStablePartialResult is a boolean value to request for stabilizing translation
	// partial results by omitting words in the end.
	SpeechServiceResponseTranslationRequestStablePartialResult PropertyID = 4100

	// SpeechServiceResponseRequestWordBoundary is a boolean value specifying whether to request WordBoundary events.
	// Added in version 1.21.0.
	SpeechServiceResponseRequestWordBoundary PropertyID = 4200

	// SpeechServiceResponseRequestPunctuationBoundary is a boolean value specifying whether to request punctuation boundary
	// in WordBoundary Events. Default is true.
	// Added in version 1.21.0.
	SpeechServiceResponseRequestPunctuationBoundary PropertyID = 4201

	// SpeechServiceResponseRequestSentenceBoundary ia a boolean value specifying whether to request sentence boundary
	// in WordBoundary Events. Default is false.
	// Added in version 1.21.0.
	SpeechServiceResponseRequestSentenceBoundary PropertyID = 4202

	// SpeechServiceResponseJSONResult is the Cognitive Services Speech Service response output (in JSON format). This
	// property is available on recognition result objects only.
	SpeechServiceResponseJSONResult PropertyID = 5000

	// SpeechServiceResponseJSONErrorDetails is the Cognitive Services Speech Service error details (in JSON format).
	// Under normal circumstances, you shouldn't have to use this property directly.
	// Instead, use CancellationDetails.ErrorDetails.
	SpeechServiceResponseJSONErrorDetails PropertyID = 5001

	// SpeechServiceResponseRecognitionLatencyMs is the recognition latency in milliseconds. Read-only, available on final
	// speech/translation results. This measures the latency between when an audio input is received by the SDK, and
	// the moment the final result is received from the service. The SDK computes the time difference between the last audio
	// fragment from the audio input that is contributing to the final result, and the time the final result is received from
	// the speech service.
	SpeechServiceResponseRecognitionLatencyMs PropertyID = 5002

	// SpeechServiceResponseSynthesisFirstByteLatencyMs is the speech synthesis first byte latency in milliseconds.
	// Read-only, available on final speech synthesis results.
	// This measures the latency between when the synthesis is started to be processed, and the moment the first byte audio is available.
	// Added in version 1.17.0.
	SpeechServiceResponseSynthesisFirstByteLatencyMs PropertyID = 5010

	// SpeechServiceResponseSynthesisFinishLatencyMs is the speech synthesis all bytes latency in milliseconds.
	// Read-only, available on final speech synthesis results.
	// This measures the latency between when the synthesis is started to be processed, and the moment the whole audio is synthesized.
	// Added in version 1.17.0.
	SpeechServiceResponseSynthesisFinishLatencyMs PropertyID = 5011

	// SpeechServiceResponseSynthesisUnderrunTimeMs is the underrun time for speech synthesis in milliseconds.
	// Read-only, available on results in SynthesisCompleted events.
	// This measures the total underrun time from AudioConfigPlaybackBufferLengthInMs is filled to synthesis completed.
	// Added in version 1.17.0.
	SpeechServiceResponseSynthesisUnderrunTimeMs PropertyID = 5012

	// SpeechServiceResponseSynthesisConnectionLatencyMs is the speech synthesis connection latency in milliseconds.
	// Read-only, available on final speech synthesis results.
	// This measures the latency between when the synthesis is started to be processed, and the moment the HTTP/WebSocket connection is established.
	// Added in version 1.31.0
	SpeechServiceResponseSynthesisConnectionLatencyMs PropertyID = 5013

	// SpeechServiceResponseSynthesisNetworkLatencyMs is the speech synthesis network latency in milliseconds.
	// Read-only, available on final speech synthesis results.
	// This measures the network round trip time.
	// Added in version 1.31.0
	SpeechServiceResponseSynthesisNetworkLatencyMs PropertyID = 5014

	// SpeechServiceResponseSynthesisServiceLatencyMs is the speech synthesis service latency in milliseconds.
	// Read-only, available on final speech synthesis results.
	// This measures the service processing time to synthesize the first byte of audio.
	// Added in version 1.31.0
	SpeechServiceResponseSynthesisServiceLatencyMs PropertyID = 5015

	// SpeechServiceResponseSynthesisBackend indicates which backend the synthesis is finished by.
	// Read-only, available on speech synthesis results, except for the result in SynthesisStarted event
	// Added in version 1.17.0.
	SpeechServiceResponseSynthesisBackend PropertyID = 5020

	// CancellationDetailsReason is the cancellation reason. Currently unused.
	CancellationDetailsReason PropertyID = 6000

	// CancellationDetailsReasonText the cancellation text. Currently unused.
	CancellationDetailsReasonText PropertyID = 6001

	// CancellationDetailsReasonDetailedText is the cancellation detailed text. Currently unused.
	CancellationDetailsReasonDetailedText PropertyID = 6002

	// AudioConfigDeviceNameForCapture is the device name for audio capture. Under normal circumstances, you shouldn't have
	// to use this property directly.
	// Instead, use AudioConfig.FromMicrophoneInput.
	AudioConfigDeviceNameForCapture PropertyID = 8000

	// AudioConfigNumberOfChannelsForCapture is the number of channels for audio capture. Internal use only.
	AudioConfigNumberOfChannelsForCapture PropertyID = 8001

	// AudioConfigSampleRateForCapture is the sample rate (in Hz) for audio capture. Internal use only.
	AudioConfigSampleRateForCapture PropertyID = 8002

	// AudioConfigBitsPerSampleForCapture is the number of bits of each sample for audio capture. Internal use only.
	AudioConfigBitsPerSampleForCapture PropertyID = 8003

	// AudioConfigAudioSource is the audio source. Allowed values are "Microphones", "File", and "Stream".
	AudioConfigAudioSource PropertyID = 8004

	// AudioConfigDeviceNameForRender indicates the device name for audio render. Under normal circumstances,
	// you shouldn't have to use this property directly. Instead, use NewAudioConfigFromDefaultSpeakerOutput.
	// Added in version 1.17.0
	AudioConfigDeviceNameForRender PropertyID = 8005

	// AudioConfigPlaybackBufferLengthInMs indicates the playback buffer length in milliseconds, default is 50 milliseconds.
	AudioConfigPlaybackBufferLengthInMs PropertyID = 8006

	// AudioProcessingOptions provides advanced configuration for audio input for features like Voice Activity Detection
	// and is provided in the form of a JSON string.
	AudioProcessingOptions PropertyID = 8007

	// SpeechLogFilename is the file name to write logs.
	SpeechLogFilename PropertyID = 9001

	// SegmentationSilenceTimeoutMs specifies a duration of detected silence, measured in milliseconds, after which
	// speech-to-text will determine a spoken phrase has ended and generate a final Recognized result. Configuring
	// this timeout may be helpful in situations where spoken input is significantly faster or slower than usual and
	// default segmentation behavior consistently yields results that are too long or too short. Segmentation timeout
	// values that are inappropriately high or low can negatively affect speech-to-text accuracy; this property should
	// be carefully configured and the resulting behavior should be thoroughly validated as intended.
	//
	// For more information about timeout configuration that includes discussion of default behaviors, please visit
	// https://aka.ms/csspeech/timeouts.
	SegmentationSilenceTimeoutMs PropertyID = 9002

	// SegmentationMaximumTimeMs represents the maximum length of a spoken phrase when using the Time segmentation strategy.
	// As the length of a spoken phrase approaches this value, the SegmentationSilenceTimeoutMs will be reduced until either 
	// the phrase silence timeout is reached or the phrase reaches the maximum length.
	SegmentationMaximumTimeMs PropertyID = 9003

	// SegmentationStrategy defines the strategy used to determine when a spoken phrase has ended and a final Recognized result should be generated.
	// Allowed values are "Default", "Time", and "Semantic".
	//
	// Valid values:
	// - "Default": Uses the default strategy and settings as determined by the Speech Service. Suitable for most situations.
	// - "Time": Uses a time-based strategy where the amount of silence between speech determines when to generate a final result.
	// - "Semantic": Uses an AI model to determine the end of a spoken phrase based on the phrase's content.
	//
	// Additional Notes:
	// - When using the Time strategy, SegmentationSilenceTimeoutMs can be adjusted to modify the required silence duration for ending a phrase, 
	//   and SegmentationMaximumTimeMs can be adjusted to set the maximum length of a spoken phrase.
	// - The Semantic strategy does not have any adjustable properties.
	SegmentationStrategy PropertyID = 9004

	// ConversationApplicationID is the identifier used to connect to the backend service.
	ConversationApplicationID PropertyID = 10000

	// ConversationDialogType is the type of dialog backend to connect to.
	ConversationDialogType PropertyID = 10001

	// ConversationInitialSilenceTimeout is the silence timeout for listening.
	ConversationInitialSilenceTimeout PropertyID = 10002

	// ConversationFromID is the FromId to be used on speech recognition activities.
	ConversationFromID PropertyID = 10003

	// ConversationConversationID is the ConversationId for the session.
	ConversationConversationID PropertyID = 10004

	// ConversationCustomVoiceDeploymentIDs is a comma separated list of custom voice deployment ids.
	ConversationCustomVoiceDeploymentIDs PropertyID = 10005

	// ConversationSpeechActivityTemplate is use to stamp properties in the template on the activity generated by the service for speech.
	ConversationSpeechActivityTemplate PropertyID = 10006

	// DataBufferTimeStamp is the time stamp associated to data buffer written by client when using Pull/Push
	// audio input streams.
	// The time stamp is a 64-bit value with a resolution of 90 kHz. It is the same as the presentation timestamp
	// in an MPEG transport stream. See https://en.wikipedia.org/wiki/Presentation_timestamp
	DataBufferTimeStamp PropertyID = 11001

	// DataBufferUserID is the user id associated to data buffer written by client when using Pull/Push audio
	// input streams.
	DataBufferUserID PropertyID = 11002
)
