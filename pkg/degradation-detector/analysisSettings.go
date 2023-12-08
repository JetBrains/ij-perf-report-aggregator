package degradation_detector

type AnalysisSettings struct {
  DoNotReportImprovement    bool
  MinimumSegmentLength      int
  MedianDifferenceThreshold float64
  EffectSizeThreshold       float64
}

func (s AnalysisSettings) GetDoNotReportImprovement() bool {
  return s.DoNotReportImprovement
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
