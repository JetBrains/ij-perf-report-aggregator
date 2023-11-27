package analysis

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "log/slog"
  "net/http"
)

func GenerateStartupSettingsForIDEA(backendUrl string, client *http.Client) []detector.StartupSettings {
  settings := make([]detector.StartupSettings, 0, 100)
  mainSettings := detector.StartupSettings{
    Channel:     "ij-perf-startup-tests",
    Branch:      "master",
    Machine:     "intellij-linux-hw-munit-%",
    Product:     "IU",
    ProductLink: "intellij",
  }
  projects, err := detector.FetchAllProjects(backendUrl, client, mainSettings)
  if err != nil {
    slog.Error("error while getting projects", "error", err)
    return settings
  }
  metrics := []string{"appInit_d", "app initialization.end", "bootstrap_d",
    "classLoadingLoadedCount", "classLoadingPreparedCount", "editorRestoring",
    "codeAnalysisDaemon/fusExecutionTime", "runDaemon/executionTime"}
  for _, project := range projects {
    for _, metric := range metrics {
      settings = append(settings, detector.StartupSettings{
        Channel:     mainSettings.Channel,
        Branch:      mainSettings.Branch,
        Machine:     mainSettings.Machine,
        Product:     mainSettings.Product,
        ProductLink: mainSettings.ProductLink,
        Project:     project,
        Metric:      metric,
        CommonAnalysisSettings: detector.CommonAnalysisSettings{
          DoNotReportImprovement: true,
        },
      })
    }
  }
  return settings
}
