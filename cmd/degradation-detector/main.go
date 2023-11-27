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
  analysisSettings := make([]detector.Settings, 0, 2000)
  for _, setting := range generatePerformanceSettings(backendUrl, client) {
    analysisSettings = append(analysisSettings, setting)
  }
  for _, setting := range generateStartupSettings(backendUrl, client) {
    analysisSettings = append(analysisSettings, setting)
  }
  degradations := getDegradations(analysisSettings, client, backendUrl)
  insertionResults := detector.PostDegradations(client, backendUrl, degradations)
  filteredResults := filterErrors(insertionResults)
  mergedResults := detector.MergeDegradations(filteredResults)
  sendDegradationsToSlack(mergedResults, client)
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

func generateStartupSettings(backendUrl string, client *http.Client) []detector.StartupSettings {
  settings := make([]detector.StartupSettings, 0, 1000)
  settings = append(settings, analysis.GenerateStartupSettingsForIDEA(backendUrl, client)...)
  return settings
}

func generatePerformanceSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  settings := make([]detector.PerformanceSettings, 0, 1000)
  settings = append(settings, analysis.GenerateIdeaSettings()...)
  settings = append(settings, analysis.GenerateWorkspaceSettings()...)
  settings = append(settings, analysis.GenerateKotlinSettings()...)
  settings = append(settings, analysis.GenerateMavenSettings()...)
  settings = append(settings, analysis.GenerateGradleSettings()...)
  settings = append(settings, analysis.GeneratePhpStormSettings()...)
  settings = append(settings, analysis.GenerateUnitTestsSettings(backendUrl, client)...)
  return settings
}

func getDegradations(settings []detector.Settings, client *http.Client, backendUrl string) <-chan detector.DegradationWithContext {
  degradationChan := make(chan detector.DegradationWithContext, 5)
  go func() {
    defer close(degradationChan)
    var wg sync.WaitGroup
    pool := pond.New(5, 1000)
    for _, setting := range settings {
      wg.Add(1)
      setting := setting
      pool.Submit(func() {
        defer wg.Done()
        slog.Info("fetching from clickhouse", "settings", setting)
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        timestamps, values, builds, err := detector.GetDataFromClickhouse(ctx, client, backendUrl, setting)
        if err != nil {
          slog.Error("error while getting data from clickhouse", "error", err, "settings", setting)
          return
        }
        for _, degradation := range detector.InferDegradations(values, builds, timestamps, setting) {
          degradationChan <- detector.DegradationWithContext{Details: degradation, Settings: setting}
        }
      })
    }
    wg.Wait()
  }()

  return degradationChan
}

func filterErrors(insertionResults <-chan detector.InsertionResults) <-chan detector.DegradationWithContext {
  ch := make(chan detector.DegradationWithContext)
  go func() {
    for result := range insertionResults {
      if result.Error != nil {
        slog.Error("error while inserting degradation", "error", result.Error, "degradation", result.Degradation)
        continue
      }
      ch <- result.Degradation
    }
    close(ch)
  }()
  return ch
}

func sendDegradationsToSlack(insertionResults <-chan detector.MultipleDegradationWithContext, client *http.Client) {
  var wg sync.WaitGroup
  for result := range insertionResults {
    wg.Add(1)
    go func(result detector.MultipleDegradationWithContext) {
      defer wg.Done()
      ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
      err := detector.SendSlackMessage(ctx, client, result)
      if err != nil {
        slog.Error("error while sending slack message", "error", err)
      }
      cancel()
    }(result)
  }
  wg.Wait()
}
