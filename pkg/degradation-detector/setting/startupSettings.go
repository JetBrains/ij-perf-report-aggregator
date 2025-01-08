package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateStartupSettingsForIDEA(backendUrl string, client *http.Client) []detector.StartupSettings {
	settings := make([]detector.StartupSettings, 0, 100)
	mainSettings := detector.StartupSettings{
		Product: "IU",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-hw-de-unit-%",
		},
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
	machines := []string{"intellij-windows-hw-de-%", "intellij-linux-hw-de-unit-%", "intellij-macos-perf-eqx-%"}
	metrics := []string{
		"appInit_d", "app initialization.end", "bootstrap_d",
		"classLoadingLoadedCount", "classLoadingPreparedCount", "editorRestoring",
		"codeAnalysisDaemon/fusExecutionTime", "runDaemon/executionTime", "startup/fusTotalDuration", "exitMetrics/application.exit",
	}
	for _, machine := range machines {
		for _, project := range projects {
			for _, metric := range metrics {
				settings = append(settings, detector.StartupSettings{
					Product: mainSettings.Product,
					Project: project,
					BaseSettings: detector.BaseSettings{
						Branch:  mainSettings.Branch,
						Machine: machine,
						Metric:  metric,
						AnalysisSettings: detector.AnalysisSettings{
							MinimumSegmentLength: 12,
						},
						SlackSettings: slackSettings,
					},
				})
			}
		}
	}
	return settings
}
