package setting

import (
	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateFleetStartupSettings() []detector.FleetStartupSettings {
	settings := make([]detector.FleetStartupSettings, 0, 100)
	machines := []string{"intellij-linux-hw-munit-%", "intellij-windows-hw-munit-%", "intellij-macos-perf-eqx-%"}
	mainSettings := detector.FleetStartupSettings{
		Branch: "master",
	}
	slackSettings := detector.SlackSettings{
		Channel:     "fleet-performance-tests-notifications",
		ProductLink: "fleet",
	}
	metrics := []string{"editor appeared", "time to edit", "terminal ready", "file tree rendered", "highlighting done", "window appeared"}

	for _, machine := range machines {
		for _, metric := range metrics {
			settings = append(settings, detector.FleetStartupSettings{
				Branch:  mainSettings.Branch,
				Machine: machine,
				Metric:  metric + ".end",
				AnalysisSettings: detector.AnalysisSettings{
					MinimumSegmentLength:      7,
					MedianDifferenceThreshold: 5,
					EffectSizeThreshold:       0.5,
					ReportType:                detector.DegradationEvent,
				},
				SlackSettings: slackSettings,
			})
		}
	}
	return settings
}
