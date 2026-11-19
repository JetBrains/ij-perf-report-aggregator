package setting

import (
	"strings"
	"testing"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/stretchr/testify/assert"
)

func TestStaticGeneratorsProduceSettings(t *testing.T) {
	t.Parallel()
	for _, g := range Generators() {
		if !g.Static {
			continue
		}
		t.Run(g.Name, func(t *testing.T) {
			t.Parallel()
			assert.NotEmpty(t, g.Generate("", nil), "generator %s produced no settings", g.Name)
		})
	}
}

func TestKotlinSetting(t *testing.T) {
	t.Parallel()
	settings := make([]detector.PerformanceSettings, 0, 1000)
	settings = append(settings, GenerateKotlinSettings()...)
	for _, setting := range settings {
		assert.True(t, setting.AnalysisSettings.ReportType == detector.ImprovementEvent || setting.AnalysisSettings.ReportType == detector.DegradationEvent)
	}
}

func TestMavenSetting(t *testing.T) {
	t.Parallel()
	settings := make([]detector.PerformanceSettings, 0, 1000)
	settings = append(settings, GenerateMavenSettings()...)
	for _, setting := range settings {
		assert.Equal(t, detector.DegradationEvent, setting.AnalysisSettings.ReportType)
	}
}

func TestGenerateAIATestTokenSettings(t *testing.T) {
	t.Parallel()
	settings := make([]detector.PerformanceSettings, 0, 1000)
	settings = append(settings, GenerateAIATestTokenSettings()...)
	for _, setting := range settings {
		assert.Equal(t, detector.ThresholdAnalysis, setting.AnalysisSettings.AnalysisKind)
		assert.Equal(t, detector.ThresholdGreaterThan, setting.AnalysisSettings.ThresholdMode)
		assert.InEpsilon(t, 95.0, setting.AnalysisSettings.ThresholdValue, 0.0001)
	}
}

func TestAirSetting(t *testing.T) {
	t.Parallel()
	for _, s := range GenerateAirSettings() {
		analysis := s.AnalysisSettings
		assert.Equal(t, detector.DegradationEvent, analysis.ReportType, s.Metric)
		if analysis.AnalysisKind == detector.ThresholdAnalysis {
			assert.Equal(t, airOverWindow, analysis.ThresholdWindow, "%s compares the median of a few nights", s.Metric)
			assert.True(t, analysis.ThresholdOnCrossing, "%s reports the crossing only", s.Metric)
		}
		if strings.Contains(s.Metric, "streamEndToSettledMs") {
			assert.Equal(t, detector.ChangePointAnalysis, analysis.AnalysisKind, "%s is always over 300 ms, so it looks for a step", s.Metric)
		}
	}
}
