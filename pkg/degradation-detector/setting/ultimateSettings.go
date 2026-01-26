package setting

import (
	"net/http"
	"slices"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateUltimateSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	return slices.Concat(
		generateUltimateDevAnalysisSettings(backendUrl, client),
	)
}

func generateUltimateDevAnalysisSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	tests := []string{
		"keycloak_release_20/%", "train-ticket/%", "toolbox_enterprise/%", "json_schema_modes_comparison/%", "json_azure/%", "swagger_indexing/%",
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
	settings := make([]detector.PerformanceSettings, 0, 100)
	machines := []string{"intellij-linux-performance-aws-%", "intellij-windows-performance-%"}
	modes := []string{"split", ""}
	for _, mode := range modes {
		for _, machine := range machines {
			for _, test := range testsExpanded {
				metrics := getUltimateMetricsFromTestsNames(test)
				for _, metric := range metrics {
					settings = append(settings, detector.PerformanceSettings{
						Db:      baseSettings.Db,
						Table:   baseSettings.Table,
						Project: test,
						Mode:    mode,
						BaseSettings: detector.BaseSettings{
							Branch:  baseSettings.Branch,
							Machine: machine,
							Metric:  metric,
							SlackSettings: detector.SlackSettings{
								Channel:     "ij-u-team-performance-issues-check",
								ProductLink: "intellij",
							},
							AnalysisSettings: detector.AnalysisSettings{MinimumSegmentLength: 8},
						},
					})
				}
			}
		}
	}
	return settings
}

func getUltimateMetricsFromTestsNames(test string) []string {
	if strings.Contains(test, "/localInspection") {
		return []string{"localInspections", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/ultimateCase") {
		return []string{"localInspections", "firstCodeAnalysis", "completion", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"typingCodeAnalyzing", "test#average_awt_delay", "test#max_awt_delay", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion", "firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"}
	}
	if strings.Contains(test, "/indexing") {
		return []string{"indexingTimeWithoutPauses", "scanningTimeWithoutPauses"}
	}

	return []string{}
}
