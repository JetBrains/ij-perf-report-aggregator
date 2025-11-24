package setting

import (
	"log/slog"
	"net/http"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateClionSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	baseSettings := detector.PerformanceSettings{
		Db:    "perfintDev",
		Table: "clion",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-performance-aws-%",
		},
	}
	branches := []string{"252", "253", "master"}
	tests, err := detector.FetchAllTests(backendUrl, client, baseSettings)
	settings := make([]detector.PerformanceSettings, 0, 100)
	if err != nil {
		slog.Error("error while getting tests", "error", err)
		return settings
	}
	modes := []string{"split", ""}
	for _, branch := range branches {
		for _, mode := range modes {
			for _, test := range tests {
				metrics := getClionMetricFromTestName(test)
				for _, metric := range metrics {
					settings = append(settings, detector.PerformanceSettings{
						Db:      baseSettings.Db,
						Table:   baseSettings.Table,
						Project: test,
						Mode:    mode,
						BaseSettings: detector.BaseSettings{
							Branch:  branch,
							Machine: baseSettings.Machine,
							Metric:  metric,
							SlackSettings: detector.SlackSettings{
								Channel:     "cidr-radler-perf-tests",
								ProductLink: "clion",
							},
						},
					})
				}
			}
		}
	}

	return settings
}

func getClionMetricFromTestName(test string) []string {
	if strings.Contains(test, "/gotoDeclaration") {
		return []string{
			"clionGotoDeclaration",
		}
	}
	if strings.Contains(test, "checkLocalTestConfig") {
		return []string{"waitFirstTestGutter"}
	}
	if strings.Contains(test, "/indexing") {
		return []string{"ocSymbolBuildingTimeMs", "backendIndexingTimeMs", "cidr.workspace.metrics#duration_in_write_action_ms", "ui.lagging#average", "rd.memory.allocatedManagedMemoryMb/afterIndexing"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"fus_time_to_show_90p"}
	}
	if strings.Contains(test, "/findUsages") {
		return []string{"%syncAction FindUsages"}
	}
	if strings.Contains(test, "/gotoDeclaration") {
		return []string{"clionGotoDeclaration"}
	}
	if strings.Contains(test, "/checkLocalTestConfig") {
		return []string{"waitFirstTestGutter"}
	}
	if strings.Contains(test, "/measureResolve") {
		return []string{"nova_total_memory_mb"}
	}
	return getMetricFromTestName(test)
}
