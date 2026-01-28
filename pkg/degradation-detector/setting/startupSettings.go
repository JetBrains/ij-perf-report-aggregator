package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateStartupSettingsForIDEA(backendUrl string, client *http.Client) []detector.StartupSettings {
	settings := make([]detector.StartupSettings, 0, 100)
	mainSettings := detector.StartupSettings{
		Db:      "perfintDev",
		Table:   "idea",
		Product: "IU",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-hw-de-unit-%",
		},
	}
	slackSettings := detector.SlackSettings{
		Channel:     "ij-u-team-performance-issues-check",
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
		"fus_daemon_finished_full_duration_since_started_ms", "runDaemon/executionTime", "fus_startup_totalDuration", "exitMetrics/application.exit", "fus_reopen_startup_code_loaded_and_visible_in_editor",
	}
	for _, machine := range machines {
		for _, project := range projects {
			for _, metric := range metrics {
				settings = append(settings, detector.StartupSettings{
					Db:      mainSettings.Db,
					Table:   mainSettings.Table,
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
