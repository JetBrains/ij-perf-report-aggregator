package setting

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

type testEditorMetricDef struct {
	test   []string
	metric []string
}

func GenerateEditorSettings() []detector.PerformanceSettings {
	testMetrics := []testEditorMetricDef{
		{test: []string{"intellij_commit/localInspection/java_file"}, metric: []string{"firstCodeAnalysis"}},
		{test: []string{"toolbox_enterprise/ultimateCase/SecurityTests"}, metric: []string{"typingCodeAnalyzing"}},
		{test: []string{"spring_boot_maven/inspection"}, metric: []string{"globalInspections"}},
		{test: []string{"intellij_commit/editor-highlighting"}, metric: []string{"typing_}_duration", "typing_EditorBackSpace_duration"}},
		{test: []string{"kotlin_coroutines/localInspection"}, metric: []string{"localInspections", "firstCodeAnalysis"}},
		{test: []string{"intellij_commit/editor-kotlin-highlighting"}, metric: []string{"typing_}_duration", "typing_EditorBackSpace_duration"}},
	}

	machines := []string{"intellij-linux-performance-aws-%"} // uncomment latter to cover all os
	// "intellij-windows-performance-aws-%",
	// "intellij-macos-perf-eqx-%",

	settings := make([]detector.PerformanceSettings, 0, 100)
	for _, testMetric := range testMetrics {
		for _, test := range testMetric.test {
			for _, metric := range testMetric.metric {
				for _, machine := range machines {
					settings = append(settings, detector.PerformanceSettings{
						Db:      "perfintDev",
						Table:   "idea",
						Project: test,
						BaseSettings: detector.BaseSettings{
							Machine: machine,
							Metric:  metric,
							Branch:  "master",
							SlackSettings: detector.SlackSettings{
								Channel:     "ij-highlighting-wg",
								ProductLink: "intellij",
							},
							AnalysisSettings: detector.AnalysisSettings{
								MinimumSegmentLength:      7,
								MedianDifferenceThreshold: 10,
							},
						},
					})
				}
			}
		}
	}
	return settings
}
