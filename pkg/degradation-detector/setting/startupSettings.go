package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateStartupSettingsForIDEA(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	settings := make([]detector.PerformanceSettings, 0, 100)
	mainSettings := detector.PerformanceSettings{
		Db:      "perfintDev",
		Table:   "idea",
		Branch:  "master",
		Machine: "intellij-linux-hw-de-unit-%",
	}
	// Notifications are off for now: the IntelliJ Ultimate team does not want startup alerts in
	// ij-u-team-performance-issues-check. With no channel, degradations are still detected and stored,
	// but neither degradation nor missing-data messages are sent.
	slackSettings := detector.SlackSettings{
		Channel:     "",
		ProductLink: "intellij",
	}
	projects, err := fetchStartupProjects(backendUrl, client, mainSettings)
	if err != nil {
		slog.Error("error while getting projects", "error", err)
		return settings
	}
	machines := []string{"intellij-windows-hw-de-%", "intellij-linux-hw-de-unit-%", "intellij-macos-perf-eqx-%"}
	metrics := []string{
		"appInit_d", "app initialization.end", "bootstrap_d",
		"classLoadingLoadedCount", "classLoadingPreparedCount", "editorRestoring",
		"fus_daemon_finished_full_duration_since_started_ms", "runDaemon/executionTime", "fus_startup_totalDuration", "exitMetrics/application.exit",
		"fus_reopen_startup_code_loaded_and_visible_in_editor",
		"fus_reopen_startup_first_ui_shown",
		"fus_reopen_startup_frame_became_interactive",
		"fus_reopen_startup_frame_became_visible",
	}
	for _, machine := range machines {
		for _, project := range projects {
			for _, metric := range metrics {
				settings = append(settings, detector.PerformanceSettings{
					Db:                   mainSettings.Db,
					Table:                mainSettings.Table,
					Project:              project,
					Branch:               mainSettings.Branch,
					Machine:              machine,
					Metric:               metric,
					MinimumSegmentLength: 12,
					ReportType:           detector.DegradationEvent,
					SlackSettings:        slackSettings,
				})
			}
		}
	}
	return settings
}
