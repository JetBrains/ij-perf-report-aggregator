package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GeneratePhpStormSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	baseSettings := detector.PerformanceSettings{
		Db:    "perfintDev",
		Table: "phpstorm",
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
	modes := []string{"split", ""}
	for _, mode := range modes {
		for _, test := range tests {
			metrics := getMetricFromTestName(test)
			for _, metric := range metrics {
				settings = append(settings, detector.PerformanceSettings{
					Db:      baseSettings.Db,
					Table:   baseSettings.Table,
					Project: test,
					Mode:    mode,
					BaseSettings: detector.BaseSettings{
						Branch:  baseSettings.Branch,
						Machine: baseSettings.Machine,
						Metric:  metric,
						SlackSettings: detector.SlackSettings{
							Channel:     "phpstorm-performance-degradations",
							ProductLink: "phpstorm",
						},
					},
				})
			}
		}
	}
	return settings
}
