package analysis

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "log/slog"
  "net/http"
)

func GeneratePhpStormSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  baseSettings := detector.PerformanceSettings{
    Db:      "perfint",
    Table:   "phpstorm",
    Branch:  "master",
    Machine: "intellij-linux-hw-hetzner%",
  }
  tests, err := detector.FetchAllTests(backendUrl, client, baseSettings)
  settings := make([]detector.PerformanceSettings, 0, 100)
  if err != nil {
    slog.Error("error while getting tests", "error", err)
    return settings
  }
  for _, test := range tests {
    metrics := getMetricFromTestName(test)
    for _, metric := range metrics {
      settings = append(settings, detector.PerformanceSettings{
        Db:      baseSettings.Db,
        Table:   baseSettings.Table,
        Branch:  baseSettings.Branch,
        Machine: baseSettings.Machine,
        Project: test,
        Metric:  metric,
        SlackSettings: detector.SlackSettings{
          Channel:     "phpstorm-performance-degradations",
          ProductLink: "phpstorm",
        },
      })
    }

  }
  return settings
}
