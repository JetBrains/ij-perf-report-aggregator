package setting

import (
	"log/slog"
	"net/http"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateRustPerfSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	machines := []string{"intellij-linux-performance-aws-%", "intellij-windows-performance-aws-%"}
	baseSettings := detector.PerformanceSettings{
		Db:    "perfint",
		Table: "rust",
		BaseSettings: detector.BaseSettings{
			Branch: "master",
		},
	}
	tests, err := detector.FetchAllTests(backendUrl, client, baseSettings)
	settings := make([]detector.PerformanceSettings, 0, 100)
	if err != nil {
		slog.Error("error while getting tests", "error", err)
		return settings
	}
	for _, machine := range machines {
		for _, test := range tests {
			metrics := getRustMetricFromTestName(test)
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
							Channel:     "rust-alerts",
							ProductLink: "rust",
						},
					},
				})
			}
		}
	}
	return settings
}

func getRustMetricFromTestName(test string) []string {
	if strings.Contains(test, "/indexing") {
		return []string{
			"rust_duration_from_start_to_work", "rust_duration_from_start_to_cargo_sync", "cargo_sync_execution_time", "rust_cargo_metadata_time", "rust_buildscript_evaluation_time",
			"rust_stdlib_fetching_time", "rust_def_maps_execution_time", "rust_macro_expansion_execution_time", "rust_class_instances_tree_size_mb", "indexingTimeWithoutPauses", "scanningTimeWithoutPauses",
		}
	}
	if strings.Contains(test, "/local-inspection") {
		return []string{"typingCodeAnalyzing#mean_value", "firstCodeAnalysis"}
	}
	if strings.Contains(test, "/global-inspection") {
		return []string{"globalInspections"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"completion#mean_value"}
	}
	if strings.Contains(test, "/findUsages") {
		return []string{"findUsages"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"typing#latency#mean_value", "typing#latency#max"}
	}
	return []string{}
}
