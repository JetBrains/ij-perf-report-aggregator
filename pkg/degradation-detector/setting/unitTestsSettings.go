package setting

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "log/slog"
  "net/http"
)

func GenerateUnitTestsSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  settings := make([]detector.PerformanceSettings, 0, 1000)
  mainSettings := detector.PerformanceSettings{
    Db:      "perfUnitTests",
    Table:   "report",
    Branch:  "master",
    Machine: "intellij-linux-%-hetzner-%",
    Metric:  "attempt.mean.ms",
  }
  slackSettings := detector.SlackSettings{
    Channel:     "ij-perf-unit-tests",
    ProductLink: "perfUnit",
  }
  tests, err := detector.FetchAllTests(backendUrl, client, mainSettings)
  if err != nil {
    slog.Error("error while getting tests", "error", err)
    return settings
  }
  for _, test := range tests {
    settings = append(settings, detector.PerformanceSettings{
      Project:       test,
      Db:            mainSettings.Db,
      Table:         mainSettings.Table,
      Branch:        mainSettings.Branch,
      Machine:       mainSettings.Machine,
      Metric:        mainSettings.Metric,
      SlackSettings: slackSettings,
      AnalysisSettings: detector.AnalysisSettings{
        DoNotReportImprovement:    true,
        MinimumSegmentLength:      30,
        MedianDifferenceThreshold: 10,
        EffectSizeThreshold:       2,
      },
    })
  }
  return settings
}
