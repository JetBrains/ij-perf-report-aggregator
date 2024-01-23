package setting

import (
	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKotlinSetting(t *testing.T) {
	settings := make([]detector.PerformanceSettings, 0, 1000)
	settings = append(settings, GenerateKotlinSettings()...)
	for _, setting := range settings {
		assert.True(t, setting.AnalysisSettings.ReportType == detector.ImprovementEvent || setting.AnalysisSettings.ReportType == detector.DegradationEvent)
	}
}

func TestMavenSetting(t *testing.T) {
	settings := make([]detector.PerformanceSettings, 0, 1000)
	settings = append(settings, GenerateMavenSettings()...)
	for _, setting := range settings {
		assert.True(t, setting.AnalysisSettings.ReportType == detector.AllEvent)
	}
}
