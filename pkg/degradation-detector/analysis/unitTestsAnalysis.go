package analysis

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "log/slog"
  "net/http"
)

func GenerateUnitTestsSettings(backendUrl string, client *http.Client) []detector.Settings {
  mainSettings := detector.Settings{
    Db:          "perfUnitTests",
    Table:       "report",
    Channel:     "ij-perf-unit-tests",
    Branch:      "master",
    Machine:     "%",
    Metric:      "attempt.mean.ms",
    ProductLink: "perfUnit",
  }
  tests, err := detector.GetAllTests(backendUrl, client, mainSettings)
  if err != nil {
    slog.Error("error while getting tests", "error", err)
    return nil
  }
  settings := make([]detector.Settings, 0, 1000)
  for _, test := range tests {
    settings = append(settings, detector.Settings{
      Test:                      test,
      Db:                        mainSettings.Db,
      Table:                     mainSettings.Table,
      Channel:                   mainSettings.Channel,
      Branch:                    mainSettings.Branch,
      Machine:                   mainSettings.Machine,
      Metric:                    mainSettings.Metric,
      ProductLink:               mainSettings.ProductLink,
      DoNotReportImprovement:    true,
      MinimumSegmentLength:      20,
      MedianDifferenceThreshold: 20,
    })
  }
  return settings
}
