package main

import (
  "context"
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/analysis"
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
  analysisSettings := generateAnalysisSettings(backendUrl, client)
  degradations := getDegradations(analysisSettings, client, backendUrl)
  insertionResults := writeDegradations(client, backendUrl, degradations)
  postDegradations(insertionResults, client)
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

func generateAnalysisSettings(backendUrl string, client *http.Client) []detector.Settings {
  settings := make([]detector.Settings, 0, 1000)
  settings = append(settings, analysis.GenerateIdeaSettings()...)
  settings = append(settings, analysis.GenerateWorkspaceSettings()...)
  settings = append(settings, analysis.GenerateKotlinSettings()...)
  settings = append(settings, analysis.GenerateMavenSettings()...)
  settings = append(settings, analysis.GenerateGradleSettings()...)
  settings = append(settings, analysis.GeneratePhpStormSettings()...)
  settings = append(settings, analysis.GenerateUnitTestsSettings(backendUrl, client)...)
  return settings
}

func getDegradations(analysisSettings []detector.Settings, client *http.Client, backendUrl string) []detector.Degradation {
  degradationChan := make(chan []detector.Degradation)
  var wgAnalysis sync.WaitGroup
  for _, analysisSetting := range analysisSettings {
    wgAnalysis.Add(1)
    go func(analysisSetting detector.Settings) {
      defer wgAnalysis.Done()
      slog.Info("processing", "settings", analysisSetting)
      ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
      defer cancel()
      timestamps, values, builds, err := detector.GetDataFromClickhouse(ctx, client, backendUrl, analysisSetting)
      if err != nil {
        slog.Error("error while getting data from clickhouse", "error", err, "settings", analysisSetting)
        return
      }
      degradationChan <- detector.InferDegradations(values, builds, timestamps, analysisSetting)
    }(analysisSetting)
  }

  go func() {
    wgAnalysis.Wait()
    close(degradationChan)
  }()

  degradations := make([]detector.Degradation, 0, 1000)
  for d := range degradationChan {
    degradations = append(degradations, d...)
  }
  return degradations
}

func postDegradations(insertionResults []detector.InsertionResults, client *http.Client) {
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

func writeDegradations(client *http.Client, backendUrl string, degradations []detector.Degradation) []detector.InsertionResults {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
  defer cancel()
  return detector.PostDegradations(ctx, client, backendUrl, degradations)
}
