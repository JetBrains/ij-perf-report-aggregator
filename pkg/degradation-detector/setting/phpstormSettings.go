package setting

import (
	"log/slog"
	"net/http"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateStartupSettingsForPhpStorm(backendUrl string, client *http.Client) []detector.StartupSettings {
	settings := make([]detector.StartupSettings, 0, 100)
	mainSettings := detector.StartupSettings{
		Db:      "perfintDev",
		Table:   "phpstorm",
		Product: "PS",
		BaseSettings: detector.BaseSettings{
			Branch:  "master",
			Machine: "intellij-linux-hw-de-unit-%",
		},
	}
	slackSettings := detector.SlackSettings{
		Channel:     "phpstorm-performance-degradations",
		ProductLink: "phpstorm",
	}
	projects, err := detector.FetchAllProjects(backendUrl, client, mainSettings)
	if err != nil {
		slog.Error("error while getting projects", "error", err)
		return settings
	}
	machines := []string{"intellij-windows-hw-de-%", "intellij-linux-hw-de-unit-%", "intellij-macos-perf-eqx-%"}
	metrics := []string{
		"fus_reopen_startup_code_loaded_and_visible_in_editor",
		"fus_reopen_startup_first_ui_shown",
		"fus_reopen_startup_frame_became_interactive",
		"fus_reopen_startup_frame_became_visible",
		"fus_startup_totalDuration",
		"fus_daemon_finished_full_duration_since_started_ms",
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
