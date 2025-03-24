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
		assert.Equal(t, detector.AllEvent, setting.AnalysisSettings.ReportType)
	}
}
