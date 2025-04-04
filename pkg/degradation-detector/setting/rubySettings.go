package setting

import (
	"log/slog"
	"net/http"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateRubyPerfSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	machines := []string{"intellij-linux-performance-aws-%"}
	baseSettings := detector.PerformanceSettings{
		Db:    "perfint",
		Table: "ruby",
		BaseSettings: detector.BaseSettings{
			Branch: "master",
		},
	}
	modes := []string{"", "split"}
	tests, err := detector.FetchAllTests(backendUrl, client, baseSettings)
	settings := make([]detector.PerformanceSettings, 0, 100)
	if err != nil {
		slog.Error("error while getting tests", "error", err)
		return settings
	}
	for _, mode := range modes {
		for _, machine := range machines {
			for _, test := range tests {
				metrics := getRubyMetricFromTestName(test)
				for _, metric := range metrics {
					medianThreshold := 10.0
					if metric == "gcPause" || metric == "freedMemoryByGC" {
						medianThreshold = 20
					}
					settings = append(settings, detector.PerformanceSettings{
						Db:      baseSettings.Db,
						Table:   baseSettings.Table,
						Project: test,
						Mode:    mode,
						BaseSettings: detector.BaseSettings{
							Branch:  baseSettings.Branch,
							Machine: machine,
							Metric:  metric,
							AnalysisSettings: detector.AnalysisSettings{
								MinimumSegmentLength:      14,
								MedianDifferenceThreshold: medianThreshold,
								ReportType:                detector.DegradationEvent,
							},
							SlackSettings: detector.SlackSettings{
								Channel:     "rubymine-performance-alerts",
								ProductLink: "rubymine",
							},
						},
					})
				}
			}
		}
	}
	return settings
}

func getRubyMetricFromTestName(test string) []string {
	if strings.Contains(test, "/indexing") {
		return []string{
			"indexingTimeWithoutPauses",
			"numberOfIndexedFilesWritingIndexValue", "indexSize", "scanningTimeWithoutPauses",
			"gcPause", "freedMemoryByGC",
		}
	}
	if strings.Contains(test, "inspections-test/") {
		return []string{"globalInspections", "gcPause", "freedMemoryByGC"}
	}
	if strings.Contains(test, "/typing") {
		return []string{"firstCodeAnalysis", "test#average_awt_delay", "typing"}
	}
	if strings.Contains(test, "/completion") {
		return []string{"firstCodeAnalysis", "completion", "completion#firstElementShown#mean_value"}
	}
	if strings.Contains(test, "/getSymbolMembers") {
		return []string{"getSymbolMembers", "getSymbolMembers#number"}
	}
	// very likely find usages
	return []string{"findUsages", "findUsages#number", "gcPause", "freedMemoryByGC"}
}
