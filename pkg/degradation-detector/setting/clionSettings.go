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
	branches := []string{"253", "261", "262", "master"}
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
		return []string{"ocSymbolBuildingTimeMs", "backendIndexingTimeMs", "rd.memory.allocatedManagedMemoryMb/afterIndexing"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"fus_time_to_show_90p"}
	}
	if strings.Contains(test, "/findUsages") {
		return []string{"findUsagesInToolWindow", "findUsagesInToolWindow#number"}
	}
	if strings.Contains(test, "/gotoDeclaration") {
		return []string{"clionGotoDeclaration"}
	}
	if strings.Contains(test, "/go-to-all-with-warmup") {
		return []string{"searchEverywhere", "searchEverywhere_first_elements_added"}
	}
	if strings.Contains(test, "/checkLocalTestConfig") {
		return []string{"waitFirstTestGutter"}
	}
	if strings.Contains(test, "/measureResolve") {
		return []string{"nova_total_memory_mb"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"typing#latency#mean_value"}
	}
	return getMetricFromTestName(test)
}
