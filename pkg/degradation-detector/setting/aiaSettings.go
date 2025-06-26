package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GenerateAIASettings() []detector.PerformanceSettings {
	metrics := map[string][]string{
		"gradle-calculator_SimpleInlineCompletionTest/simple cloud completion": {"callInlineCompletionOnCompletion#mean_value"},
		"gradle-calculator_CodeGenerationPerformanceTest/generate code":        {"ai-generate-code#mean_value", "doGenerate#mean_value"},
		"kotlinx_coroutines_k2_dev_ContextPrivacyFiltersTest/coroutinesProject-context-privacy-filters-all-files": {
			"contextPrivacyFilter.AiIgnoreContextPrivacyFilter.sum.ms",
			"contextPrivacyFilter.VcsUnversionedContextPrivacyFilter.sum.ms",
			"contextPrivacyFilter.VcsIgnoredContextPrivacyFilter.sum.ms",
			"ai-ignore.sum.ms",
		},
	}

	settings := make([]detector.PerformanceSettings, 0, 10)

	for test, metrics := range metrics {
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:      "perfintDev",
				Table:   "ml",
				Project: test,
				BaseSettings: detector.BaseSettings{
					Machine: "intellij-linux-performance-aws-%",
					Metric:  metric,
					Branch:  "master",
					SlackSettings: detector.SlackSettings{
						Channel:     "ai-assistant-autotest-notifications-all",
						ProductLink: "ml/dev",
					},
					AnalysisSettings: detector.AnalysisSettings{
						ReportType: detector.AllEvent,
					},
				},
			})
		}
	}
	return settings
}
