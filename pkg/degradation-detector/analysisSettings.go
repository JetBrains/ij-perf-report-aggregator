package degradation_detector

type ReportType int

const (
	AllEvent ReportType = iota
	DegradationEvent
	ImprovementEvent
)

// AnalysisKind selects which algorithm to use for detection.
// Default (zero) keeps backward-compatible change-point detection.
type AnalysisKind int

const (
	ChangePointAnalysis AnalysisKind = iota
	ThresholdAnalysis
)

// ThresholdMode specifies how the latest value should be compared to the threshold.
type ThresholdMode int

const (
	ThresholdGreaterThan ThresholdMode = iota
	ThresholdLessThan
)

type AnalysisSettings struct {
	ReportType ReportType
	// Determines the minimum length of a segment the larger the segment the more accurate the analysis but it will take more time to detect degradation
	// default value is 5 which means if you have 2 builds per day the delay will be about 3 days
	MinimumSegmentLength int
	// Determines the threshold for the median difference between two segments to be considered as a change.
	// The default value is 10 %.
	MedianDifferenceThreshold float64
	// See: https://en.wikipedia.org/wiki/Effect_size for details.
	// The setting affects how noise affects the analysis. The larger the value the more noise is ignored.
	// The default value is 2.
	EffectSizeThreshold float64
	// Number of days to check for missing data.
	// The default value is -3 (3 days ago).
	DaysToCheckMissing int

	AnalysisKind   AnalysisKind
	ThresholdMode  ThresholdMode
	ThresholdValue float64
}

func (s AnalysisSettings) GetReportType() ReportType {
	return s.ReportType
}

func (s AnalysisSettings) GetMinimumSegmentLength() int {
	return s.MinimumSegmentLength
}

func (s AnalysisSettings) GetMedianDifferenceThreshold() float64 {
	return s.MedianDifferenceThreshold
}

func (s AnalysisSettings) GetEffectSizeThreshold() float64 {
	return s.EffectSizeThreshold
}

func (s AnalysisSettings) GetDaysToCheckMissing() int {
	if s.DaysToCheckMissing == 0 {
		return -3
	}
	return s.DaysToCheckMissing
}

// GetAnalysisKind returns the analysis kind setting.
func (s AnalysisSettings) GetAnalysisKind() AnalysisKind { return s.AnalysisKind }

// GetThresholdMode returns the threshold comparison mode.
func (s AnalysisSettings) GetThresholdMode() ThresholdMode { return s.ThresholdMode }

// GetThresholdValue returns the threshold value.
func (s AnalysisSettings) GetThresholdValue() float64 { return s.ThresholdValue }
