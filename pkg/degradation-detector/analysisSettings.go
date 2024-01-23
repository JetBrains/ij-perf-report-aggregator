package degradation_detector

// enum for report type
type ReportType int

const (
	AllEvent         ReportType = 0
	DegradationEvent ReportType = 1
	ImprovementEvent ReportType = 2
)

type AnalysisSettings struct {
	ReportType                ReportType
	MinimumSegmentLength      int
	MedianDifferenceThreshold float64
	EffectSizeThreshold       float64
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
