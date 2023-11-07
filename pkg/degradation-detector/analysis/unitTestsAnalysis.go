package analysis

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "log"
)

func GenerateUnitTestsSettings(backendUrl string) []detector.Settings {
  mainSettings := detector.Settings{
    Db:          "perfUnitTests",
    Table:       "report",
    Channel:     "ij-perf-unit-tests",
    Branch:      "master",
    Machine:     "%",
    Metric:      "attempt.average.ms",
    ProductLink: "perfUnit",
  }
  tests, err := detector.GetAllTests(backendUrl, mainSettings)
  if err != nil {
    log.Printf("Can't get tests: %v", err)
    return nil
  }
  settings := make([]detector.Settings, 0, 1000)
  for _, test := range tests {
    settings = append(settings, detector.Settings{
      Db:          "perfUnitTests",
      Table:       "report",
      Channel:     "ij-perf-report-aggregator",
      Branch:      "master",
      Machine:     "%",
      Test:        test,
      Metric:      "attempt.average.ms",
      ProductLink: "perfUnit",
    })
  }
  return settings
}
