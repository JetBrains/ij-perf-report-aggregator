package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateKotlinBuildToolsSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	baseSettings := detector.PerformanceSettings{
		Db:    "perfintDev",
		Table: "kotlinBuildTools",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-hw-hetzner%",
		},
	}
	tests, err := detector.FetchAllTests(backendUrl, client, baseSettings)
	settings := make([]detector.PerformanceSettings, 0, 100)
	if err != nil {
		slog.Error("error while getting tests", "error", err)
		return settings
	}
	for _, test := range tests {
		metrics := []string{"ExternalSystemSyncProjectTask"}
		for _, metric := range metrics {
			settings = append(settings, detector.PerformanceSettings{
				Db:      baseSettings.Db,
				Table:   baseSettings.Table,
				Project: test,
				BaseSettings: detector.BaseSettings{
					Branch:  baseSettings.Branch,
					Machine: baseSettings.Machine,
					Metric:  metric,
					SlackSettings: detector.SlackSettings{
						Channel:     "kotlin-build-tools-build-notifications",
						ProductLink: "kotlinBuildTools",
					},
				},
			})
		}
	}
	return settings
}
