package degradation_detector

type CommonAnalysisSettings struct {
  DoNotReportImprovement    bool
  MinimumSegmentLength      int
  MedianDifferenceThreshold float64
}

func (s CommonAnalysisSettings) GetDoNotReportImprovement() bool {
  return s.DoNotReportImprovement
}

func (s CommonAnalysisSettings) GetMinimumSegmentLength() int {
  return s.MinimumSegmentLength
}

func (s CommonAnalysisSettings) GetMedianDifferenceThreshold() float64 {
  return s.MedianDifferenceThreshold
}
