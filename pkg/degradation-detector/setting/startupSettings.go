package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateStartupSettingsForIDEA(backendUrl string, client *http.Client) []detector.StartupSettings {
	settings := make([]detector.StartupSettings, 0, 100)
	mainSettings := detector.StartupSettings{
		Branch:  "master",
		Machine: "intellij-linux-hw-de-unit-%",
		Product: "IU",
	}
	slackSettings := detector.SlackSettings{
		Channel:     "ij-startup-idea-reports",
		ProductLink: "intellij",
	}
	projects, err := detector.FetchAllProjects(backendUrl, client, mainSettings)
	if err != nil {
		slog.Error("error while getting projects", "error", err)
		return settings
	}
	metrics := []string{
		"appInit_d", "app initialization.end", "bootstrap_d",
		"classLoadingLoadedCount", "classLoadingPreparedCount", "editorRestoring",
		"codeAnalysisDaemon/fusExecutionTime", "runDaemon/executionTime", "startup/fusTotalDuration", "exitMetrics/application.exit",
	}
	for _, project := range projects {
		for _, metric := range metrics {
			settings = append(settings, detector.StartupSettings{
				Branch:  mainSettings.Branch,
				Machine: mainSettings.Machine,
				Product: mainSettings.Product,
				Project: project,
				Metric:  metric,
				AnalysisSettings: detector.AnalysisSettings{
					MinimumSegmentLength: 12,
				},
				SlackSettings: slackSettings,
			})
		}
	}
	return settings
}
