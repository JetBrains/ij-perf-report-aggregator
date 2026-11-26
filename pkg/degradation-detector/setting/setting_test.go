package setting

import (
	"strings"
	"testing"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/stretchr/testify/assert"
)

func TestAirSetting(t *testing.T) {
	t.Parallel()
	for _, s := range GenerateAirSettings() {
		analysis := s.AnalysisSettings
		assert.Equal(t, detector.DegradationEvent, analysis.ReportType, s.Metric)
		if analysis.AnalysisKind == detector.ThresholdAnalysis {
			assert.Equal(t, airOverWindow, analysis.ThresholdWindow, "%s compares the median of a few nights", s.Metric)
			assert.True(t, analysis.ThresholdOnCrossing, "%s reports the crossing only", s.Metric)
		}
		if strings.Contains(s.Metric, "streamEndToSettledMs") || strings.HasPrefix(s.Metric, "open.") {
			assert.Equal(t, detector.StepAnalysis, analysis.AnalysisKind, "%s reports a step on its first night", s.Metric)
		}
	}
}

func TestAirSettingWatchesEveryOpenMetric(t *testing.T) {
	t.Parallel()
	metrics := map[string]bool{}
	for _, s := range GenerateAirSettings() {
		metrics[s.Metric] = true
	}
	for _, size := range []string{"M", "XL", "XXL"} {
		assert.True(t, metrics["open."+size+".gestureToTabMs"], "the time from the click to the open tab of %s", size)
		assert.True(t, metrics["open."+size+".coverRemovedMs"], "the time from the click to the first text of %s", size)
	}
}
