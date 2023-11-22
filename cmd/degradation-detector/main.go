package main

import (
  "context"
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/analysis"
  "github.com/alitto/pond"
  _ "go.uber.org/automaxprocs"
  "log/slog"
  "net/http"
  "os"
  "sync"
  "time"
)

func main() {
  backendUrl := getBackendUrl()
  client := createHttpClient()
  slog.Info("started")
  analysisSettings := generateAnalysisSettings(backendUrl, client)
  degradations := getDegradations(analysisSettings, client, backendUrl)
  insertionResults := detector.PostDegradations(client, backendUrl, degradations)
  sendDegradationsToSlack(insertionResults, client)
  slog.Info("finished")
}

func getBackendUrl() string {
  backendUrl := os.Getenv("BACKEND_URL")
  if backendUrl == "" {
    backendUrl = "https://ij-perf-api.labs.jb.gg" // Default URL
    slog.Info("BACKEND_URL is not set, using default value: %s", "url", backendUrl)
  }
  return backendUrl
}

func createHttpClient() *http.Client {
  return &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
      MaxIdleConns:        20,
      MaxIdleConnsPerHost: 10,
    },
  }
}

func generateAnalysisSettings(backendUrl string, client *http.Client) <-chan detector.Settings {
  settings := make([]detector.Settings, 0, 1000)
  settings = append(settings, analysis.GenerateIdeaSettings()...)
  settings = append(settings, analysis.GenerateWorkspaceSettings()...)
  settings = append(settings, analysis.GenerateKotlinSettings()...)
  settings = append(settings, analysis.GenerateMavenSettings()...)
  settings = append(settings, analysis.GenerateGradleSettings()...)
  settings = append(settings, analysis.GeneratePhpStormSettings()...)
  settings = append(settings, analysis.GenerateUnitTestsSettings(backendUrl, client)...)
  settingsChan := make(chan detector.Settings)
  go func() {
    for _, setting := range settings {
      settingsChan <- setting
    }
    close(settingsChan)
  }()
  return settingsChan
}

func getDegradations(analysisSettings <-chan detector.Settings, client *http.Client, backendUrl string) <-chan detector.Degradation {
  degradationChan := make(chan detector.Degradation, 5)
  go func() {
    defer close(degradationChan)
    var wg sync.WaitGroup
    pool := pond.New(5, 1000)
    for analysisSetting := range analysisSettings {
      wg.Add(1)
      analysisSetting := analysisSetting
      pool.Submit(func() {
        defer wg.Done()
        slog.Info("processing", "settings", analysisSetting)
        ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
        defer cancel()
        timestamps, values, builds, err := detector.GetDataFromClickhouse(ctx, client, backendUrl, analysisSetting)
        if err != nil {
          slog.Error("error while getting data from clickhouse", "error", err, "settings", analysisSetting)
          return
        }
        for _, degradation := range detector.InferDegradations(values, builds, timestamps, analysisSetting) {
          degradationChan <- degradation
        }
      })
    }
    wg.Wait()
  }()

  return degradationChan
}

func sendDegradationsToSlack(insertionResults <-chan detector.InsertionResults, client *http.Client) {
  var wg sync.WaitGroup
  for result := range insertionResults {
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
      ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
      err := detector.SendSlackMessage(ctx, client, result.Degradation)
      if err != nil {
        slog.Error("error while sending slack message", "error", err)
      }
      cancel()
    }(result)
  }
  wg.Wait()
}
