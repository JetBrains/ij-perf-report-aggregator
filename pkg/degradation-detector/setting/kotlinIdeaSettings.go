package setting

import (
	"net/http"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateKotlinIdeaSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	tests := []string{
		"intellij_commit/createKotlinClass",
		"intellij_commit/editor-kotlin-highlighting",
		"intellij_commit/red-code-kotlin",
		"kotlin/%",
		"kotlin_coroutines/%",
	}

	baseSettings := detector.PerformanceSettings{
		Db:    "perfintDev",
		Table: "idea",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-performance-aws-%",
		},
	}
	testsExpanded := detector.ExpandTestsByPattern(backendUrl, client, tests, baseSettings)
	testsWithoutIndexingScanning := filterIndexingScanningTests(testsExpanded)

	settings := make([]detector.PerformanceSettings, 0, 100)
	machines := []string{"intellij-linux-performance-aws-%", "intellij-windows-performance-%"}
	for _, machine := range machines {
		for _, test := range testsWithoutIndexingScanning {
			metrics := getKotlinMetricsFromTestName(test)
			for _, metric := range metrics {
				settings = append(settings, detector.PerformanceSettings{
					Db:      baseSettings.Db,
					Table:   baseSettings.Table,
					Project: test,
					BaseSettings: detector.BaseSettings{
						Branch:  baseSettings.Branch,
						Machine: machine,
						Metric:  metric,
						SlackSettings: detector.SlackSettings{
							Channel:     "kotlin-plugin-perf-tests",
							ProductLink: "intellij",
						},
						AnalysisSettings: detector.AnalysisSettings{MinimumSegmentLength: 8},
					},
				})
			}
		}
	}
	return settings
}

func getKotlinMetricsFromTestName(test string) []string {
	if strings.Contains(test, "/editor-kotlin-highlighting") {
		return []string{"typing_EditorBackSpace_duration", "typing_EditorBackSpace_warmup_duration", "typing_}_duration", "typing_}_warmup_duration"}
	}
	if strings.Contains(test, "/red-code-kotlin") {
		return []string{"replaceTextCodeAnalysis"}
	}
	if strings.Contains(test, "/createKotlinClass") {
		return []string{"createKotlinFile"}
	}
	if strings.Contains(test, "kotlin/localInspection") || strings.Contains(test, "kotlin_coroutines/localInspection") {
		return []string{"firstCodeAnalysis", "localInspections", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "kotlin/inspection") || strings.Contains(test, "kotlin_coroutines/inspection") {
		return []string{"globalInspections", "fullGCPause", "JVM.GC.collectionTimesMs"}
	}
	return []string{}
}
