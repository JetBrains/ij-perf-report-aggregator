package main

import (
  "context"
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/analysis"
  "log/slog"
  "net/http"
  "os"
  "sync"
  "time"
)

func main() {
  backendUrl := os.Getenv("BACKEND_URL")
  if len(backendUrl) == 0 {
    backendUrl = "https://ij-perf-api.labs.jb.gg" //http://localhost:9044
    slog.Info("BACKEND_URL is not set, using default value: %s", "url", backendUrl)
  }

  client := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
      MaxIdleConns:        20,
      MaxIdleConnsPerHost: 10,
    },
  }

  analysisSettings := make([]detector.Settings, 0, 1000)
  analysisSettings = append(analysisSettings, analysis.GenerateIdeaSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateWorkspaceSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateKotlinSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateMavenSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateGradleSettings()...)
  analysisSettings = append(analysisSettings, analysis.GeneratePhpStormSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateUnitTestsSettings(backendUrl, client)...)

  ctx := context.Background()
  degradations := make([]detector.Degradation, 0, 1000)
  for _, analysisSetting := range analysisSettings {
    slog.Info("processing", "settings", analysisSetting)
    timestamps, values, builds, err := detector.GetDataFromClickhouse(ctx, client, backendUrl, analysisSetting)
    if err != nil {
      slog.Error("error while getting data from clickhouse", "error", err)
    }

    degradations = append(degradations, detector.InferDegradations(values, builds, timestamps, analysisSetting)...)
  }

  insertionResults := detector.PostDegradations(ctx, client, backendUrl, degradations)

  var wg sync.WaitGroup
  for _, result := range insertionResults {
    if result.Error != nil {
      slog.Error("error while inserting degradation", "error", result.Error, "degradation", result.Degradation)
      continue
    }
    if !result.WasInserted {
      continue
    }
    wg.Add(1)
    go func(result detector.InsertionResults) {
      defer wg.Done()
      err := detector.SendSlackMessage(ctx, client, result.Degradation)
      if err != nil {
        slog.Error("error while sending slack message", "error", err)
      }
    }(result)
  }
  wg.Wait()
}
