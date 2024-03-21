package setting

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

func GenerateFleetStartupSettings() []detector.FleetStartupSettings {
  settings := make([]detector.FleetStartupSettings, 0, 100)
  mainSettings := detector.FleetStartupSettings{
    Branch:  "master",
    Machine: "intellij-linux-hw-munit-%",
  }
  slackSettings := detector.SlackSettings{
    Channel:     "fleet-performance-hackathon-feb24",
    ProductLink: "fleet",
  }
  metrics := []string{"editor appeared", "time to edit", "terminal ready", "file tree rendered", "highlighting done", "window appeared"}

  for _, metric := range metrics {
    settings = append(settings, detector.FleetStartupSettings{
      Branch:  mainSettings.Branch,
      Machine: mainSettings.Machine,
      Metric:  metric + ".end",
      AnalysisSettings: detector.AnalysisSettings{
        MinimumSegmentLength:      7,
        MedianDifferenceThreshold: 5,
        EffectSizeThreshold:       0.5,
      },
      SlackSettings: slackSettings,
    })
  }
  return settings
}
