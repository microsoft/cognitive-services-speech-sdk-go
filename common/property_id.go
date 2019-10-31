package common

type PropertyId int

const (
	SpeechServiceConnection_Key                                 PropertyId = 1000
	SpeechServiceConnection_Endpoint                            PropertyId = 1001
	SpeechServiceConnection_Region                              PropertyId = 1002
	SpeechServiceAuthorization_Token                            PropertyId = 1003
	SpeechServiceAuthorization_Type                             PropertyId = 1004
	SpeechServiceConnection_EndpointId                          PropertyId = 1005
	SpeechServiceConnection_Host                                PropertyId = 1006
	SpeechServiceConnection_ProxyHostName                       PropertyId = 1100
	SpeechServiceConnection_ProxyPort                           PropertyId = 1101
	SpeechServiceConnection_ProxyUserName                       PropertyId = 1102
	SpeechServiceConnection_ProxyPassword                       PropertyId = 1103
	SpeechServiceConnection_Url                                 PropertyId = 1104
	SpeechServiceConnection_TranslationToLanguages              PropertyId = 2000
	SpeechServiceConnection_TranslationVoice                    PropertyId = 2001
	SpeechServiceConnection_TranslationFeatures                 PropertyId = 2002
	SpeechServiceConnection_IntentRegion                        PropertyId = 2003
	SpeechServiceConnection_RecoMode                            PropertyId = 3000
	SpeechServiceConnection_RecoLanguage                        PropertyId = 3001
	Speech_SessionId                                            PropertyId = 3002
	SpeechServiceConnection_UserDefinedQueryParameters          PropertyId = 3003
	SpeechServiceConnection_SynthLanguage                       PropertyId = 3100
	SpeechServiceConnection_SynthVoice                          PropertyId = 3101
	SpeechServiceConnection_SynthOutputFormat                   PropertyId = 3102
	SpeechServiceConnection_InitialSilenceTimeoutMs             PropertyId = 3200
	SpeechServiceConnection_EndSilenceTimeoutMs                 PropertyId = 3201
	SpeechServiceConnection_EnableAudioLogging                  PropertyId = 3202
	SpeechServiceConnection_AutoDetectSourceLanguages           PropertyId = 3300
	SpeechServiceConnection_AutoDetectSourceLanguageResult      PropertyId = 3301
	SpeechServiceResponse_RequestDetailedResultTrueFalse        PropertyId = 4000
	SpeechServiceResponse_RequestProfanityFilterTrueFalse       PropertyId = 4001
	SpeechServiceResponse_ProfanityOption                       PropertyId = 4002
	SpeechServiceResponse_PostProcessingOption                  PropertyId = 4003
	SpeechServiceResponse_RequestWordLevelTimestamps            PropertyId = 4004
	SpeechServiceResponse_StablePartialResultThreshold          PropertyId = 4005
	SpeechServiceResponse_OutputFormatOption                    PropertyId = 4006
	SpeechServiceResponse_TranslationRequestStablePartialResult PropertyId = 4100
	SpeechServiceResponse_JsonResult                            PropertyId = 5000
	SpeechServiceResponse_JsonErrorDetails                      PropertyId = 5001
	SpeechServiceResponse_RecognitionLatencyMs                  PropertyId = 5002
	CancellationDetails_Reason                                  PropertyId = 6000
	CancellationDetails_ReasonText                              PropertyId = 6001
	CancellationDetails_ReasonDetailedText                      PropertyId = 6002
	LanguageUnderstandingServiceResponse_JsonResult             PropertyId = 7000
	AudioConfig_DeviceNameForCapture                            PropertyId = 8000
	AudioConfig_NumberOfChannelsForCapture                      PropertyId = 8001
	AudioConfig_SampleRateForCapture                            PropertyId = 8002
	AudioConfig_BitsPerSampleForCapture                         PropertyId = 8003
	AudioConfig_AudioSource                                     PropertyId = 8004
	Speech_LogFilename                                          PropertyId = 9001
	Conversation_ApplicationId                                  PropertyId = 10000
	Conversation_DialogType                                     PropertyId = 10001
	Conversation_Initial_Silence_Timeout                        PropertyId = 10002
	Conversation_From_Id                                        PropertyId = 10003
	Conversation_Conversation_Id                                PropertyId = 10004
	Conversation_Custom_Voice_Deployment_Ids                    PropertyId = 10005
	DataBuffer_TimeStamp                                        PropertyId = 11001
	DataBuffer_UserId                                                      = 11002
)
