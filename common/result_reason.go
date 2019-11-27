package common

// OutputFormat specifies the possible reasons a recognition result might be generated.
type ResultReason int

const (
	// Indicates speech could not be recognized. More details can be found in the NoMatchDetails object.
	NoMatch                    ResultReason = 0

	// Indicates that the recognition was canceled. More details can be found using the CancellationDetails object.
	Canceled                   ResultReason = 1

	// Indicates the speech result contains hypothesis text.
	RecognizingSpeech          ResultReason = 2

	// Indicates the speech result contains final text that has been recognized.
	// Speech Recognition is now complete for this phrase.
	RecognizedSpeech           ResultReason = 3

	// Indicates the intent result contains hypothesis text and intent.
    /// </summary>
    RecognizingIntent          ResultReason = 4


    // Indicates the intent result contains final text and intent.
    // Speech Recognition and Intent determination are now complete for this phrase.
    RecognizedIntent           ResultReason = 5

    // Indicates the translation result contains hypothesis text and its translation(s).
    TranslatingSpeech          ResultReason = 6

    // Indicates the translation result contains final text and corresponding translation(s).
    // Speech Recognition and Translation are now complete for this phrase.
    TranslatedSpeech           ResultReason = 7

    // Indicates the synthesized audio result contains a non-zero amount of audio data
    SynthesizingAudio          ResultReason = 8

    // Indicates the synthesized audio is now complete for this phrase.
    SynthesizingAudioCompleted ResultReason = 9

    // Indicates the speech result contains (unverified) keyword text.
    // Added in version 1.3.0
    RecognizingKeyword         ResultReason = 10

    // Indicates that keyword recognition completed recognizing the given keyword.
    // Added in version 1.3.0
    RecognizedKeyword          ResultReason = 11

    // Indicates the speech synthesis is now started
    // Added in version 1.4.0
    SynthesizingAudioStarted   ResultReason = 12

)
