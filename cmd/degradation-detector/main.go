package main

import (
  detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/setting"
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
  for _, s := range generatePerformanceSettings(backendUrl, client) {
    analysisSettings = append(analysisSettings, s)
  }
  for _, s := range generateStartupSettings(backendUrl, client) {
    analysisSettings = append(analysisSettings, s)
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
  settings = append(settings, setting.GenerateStartupSettingsForIDEA(backendUrl, client)...)
  return settings
}

func generatePerformanceSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
  settings := make([]detector.PerformanceSettings, 0, 1000)
  settings = append(settings, setting.GenerateIdeaSettings(backendUrl, client)...)
  settings = append(settings, setting.GenerateWorkspaceSettings()...)
  settings = append(settings, setting.GenerateKotlinSettings()...)
  settings = append(settings, setting.GenerateMavenSettings()...)
  settings = append(settings, setting.GenerateGradleSettings()...)
  settings = append(settings, setting.GeneratePhpStormSettings(backendUrl, client)...)
  settings = append(settings, setting.GenerateUnitTestsSettings(backendUrl, client)...)
  return settings
}
