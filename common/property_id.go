//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package common

// PropertyID defines speech property ids.
// Changed in version 1.4.0.
type PropertyID int

const (
	// SpeechServiceConnectionKey is the Cognitive Services Speech Service subscription key. If you are using an
	// intent recognizer, you need to specify the LUIS endpoint key for your particular LUIS app. Under normal
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

	// SpeechServiceConnectionIntentRegion is the Language Understanding Service region. Under normal circumstances, you
	// shouldn't have to use this property directly.
	// Instead use LanguageUnderstandingModel.
	SpeechServiceConnectionIntentRegion PropertyID = 2003

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
	/// Instead use SessionEventArgs.SessionId.
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

	// SpeechServiceConnectionInitialSilenceTimeoutMs is the initial silence timeout value (in milliseconds) used by the
	// service.
	SpeechServiceConnectionInitialSilenceTimeoutMs PropertyID = 3200

	// SpeechServiceConnectionEndSilenceTimeoutMs is the end silence timeout value (in milliseconds) used by the service.
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

	// SpeechServiceResponseJSONResult is the Cognitive Services Speech Service response output (in JSON format). This
	// property is available on recognition result objects only.
	SpeechServiceResponseJSONResult PropertyID = 5000

	// SpeechServiceResponseJSONErrorDetails is the Cognitive Services Speech Service error details (in JSON format).
	// Under normal circumstances, you shouldn't have to use this property directly.
	// Instead, use CancellationDetails.ErrorDetails.
	SpeechServiceResponseJSONErrorDetails PropertyID = 5001

	// SpeechServiceResponseRecognitionLatencyMs is the recognition latency in milliseconds. Read-only, available on final
	// speech/translation/intent results. This measures the latency between when an audio input is received by the SDK, and
	// the moment the final result is received from the service. The SDK computes the time difference between the last audio
	// fragment from the audio input that is contributing to the final result, and the time the final result is received from
	// the speech service.
	SpeechServiceResponseRecognitionLatencyMs PropertyID = 5002

	// CancellationDetailsReason is the cancellation reason. Currently unused.
	CancellationDetailsReason PropertyID = 6000

	// CancellationDetailsReasonText the cancellation text. Currently unused.
	CancellationDetailsReasonText PropertyID = 6001

	// CancellationDetailsReasonDetailedText is the cancellation detailed text. Currently unused.
	CancellationDetailsReasonDetailedText PropertyID = 6002

	// LanguageUnderstandingServiceResponseJSONResult is the Language Understanding Service response output (in JSON format).
	// Available via IntentRecognitionResult.Properties.
	LanguageUnderstandingServiceResponseJSONResult PropertyID = 7000

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

	// SpeechLogFilename is the file name to write logs.
	SpeechLogFilename PropertyID = 9001

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
