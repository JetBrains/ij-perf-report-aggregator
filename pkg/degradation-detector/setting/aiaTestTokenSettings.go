package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateAIATestTokenSettings() []detector.PerformanceSettings {
	metrics := map[string][]string{
		"stage: tokenQuota": {"currentPercents"},
		"prod: tokenQuota":  {"currentPercents"},
	}

	settings := make([]detector.PerformanceSettings, 0, 10)

	for test, metrics := range metrics {
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:             "perfintDev",
				Table:          "ml",
				Project:        test,
				Machine:        "intellij-linux-%-aws-%",
				Metric:         metric,
				Branch:         "master",
				Channel:        "ai-assistant-autotest-notifications",
				ProductLink:    "ml/dev",
				ReportType:     detector.AllEvent,
				AnalysisKind:   detector.ThresholdAnalysis,
				ThresholdMode:  detector.ThresholdGreaterThan,
				ThresholdValue: 95,
			})
		}
	}
	return settings
}
