package setting

import (
	"testing"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/stretchr/testify/assert"
)

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
