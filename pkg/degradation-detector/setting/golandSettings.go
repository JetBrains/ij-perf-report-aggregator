package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateStartupSettingsForGoland(backendUrl string, client *http.Client) []detector.StartupSettings {
	settings := make([]detector.StartupSettings, 0, 100)
	mainSettings := detector.StartupSettings{
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-hw-de-unit-%",
		},
		Product: "GO",
	}
	slackSettings := detector.SlackSettings{
		Channel:     "goland-qa-duty",
		ProductLink: "goland",
	}
	projects, err := detector.FetchAllProjects(backendUrl, client, mainSettings)
	if err != nil {
		slog.Error("error while getting projects", "error", err)
		return settings
	}
	metrics := []string{"startup/fusTotalDuration", "progressMetric/Progress: Updating Go modules dependencies", "metrics.progressMetric/Progress: Updating Go modules dependencies#mean_value"}
	for _, project := range projects {
		for _, metric := range metrics {
			settings = append(settings, detector.StartupSettings{
				BaseSettings: detector.BaseSettings{
					Metric: metric,
					AnalysisSettings: detector.AnalysisSettings{
						MinimumSegmentLength: 12,
					},
					SlackSettings: slackSettings,
					Branch:        mainSettings.Branch,
					Machine:       mainSettings.Machine,
				},
				Product: mainSettings.Product,
				Project: project,
			})
		}
	}
	return settings
}

func GenerateGolandPerfSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	baseSettings := detector.PerformanceSettings{
		Db:    "perfint",
		Table: "goland",
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
		metrics := getMetricFromTestName(test)
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
						Channel:     "goland-qa-duty",
						ProductLink: "goland",
					},
				},
			})
		}

	}
	return settings
}
