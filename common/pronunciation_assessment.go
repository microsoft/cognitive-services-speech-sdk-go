package common

type PronunciationAssessment_GradingSystem int

const (
	PronunciationAssessmentGradingSystemFivePoint   PronunciationAssessment_GradingSystem = 1
	PronunciationAssessmentGradingSystemHundredMark PronunciationAssessment_GradingSystem = 2
)

type PronunciationAssessment_Granularity int

const (
	PronunciationAssessmentGranularityPhoneme  PronunciationAssessment_Granularity = 1
	PronunciationAssessmentGranularityWord     PronunciationAssessment_Granularity = 2
	PronunciationAssessmentGranularityFullText PronunciationAssessment_Granularity = 3
)
