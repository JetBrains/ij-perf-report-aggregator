package main

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/analysis"
  _ "go.uber.org/automaxprocs"
  "log/slog"
  "net/http"
  "os"
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
  degradations := detector.GetDegradations(analysisSettings, client, backendUrl)
  insertionResults := detector.PostDegradations(client, backendUrl, degradations)
  filteredResults := detector.FilterErrors(insertionResults)
  mergedResults := detector.MergeDegradations(filteredResults)
  detector.SendDegradationsToSlack(mergedResults, client)
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
