package main

import (
	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/setting"
	_ "go.uber.org/automaxprocs"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"time"
)

func main() {
	backendUrl := getBackendUrl()
	client := createHttpClient()
	slog.Info("started")
	analysisSettings := make([]detector.Settings, 0, 25000)
	for _, s := range generatePerformanceSettings(backendUrl, client) {
		analysisSettings = append(analysisSettings, s)
	}
	for _, s := range generateStartupSettings(backendUrl, client) {
		analysisSettings = append(analysisSettings, s)
	}
	for _, s := range generateFleetStartupSettings() {
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
	return slices.Concat(
		setting.GenerateStartupSettingsForIDEA(backendUrl, client),
		setting.GenerateStartupSettingsForGoland(backendUrl, client),
	)
}

func generateFleetStartupSettings() []detector.FleetStartupSettings {
	return slices.Concat(
		setting.GenerateFleetStartupSettings(),
	)
}

func generatePerformanceSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	return slices.Concat(
		setting.GenerateIdeaSettings(backendUrl, client),
		setting.GenerateIdeaIndexingSettings(backendUrl, client),
		setting.GenerateWorkspaceSettings(),
		setting.GenerateKotlinSettings(),
		setting.GenerateMavenSettings(),
		setting.GenerateGradleSettings(),
		setting.GenerateVCSSettings(),
		setting.GeneratePhpStormSettings(backendUrl, client),
		setting.GenerateUnitTestsSettings(backendUrl, client),
		setting.GenerateGolandPerfSettings(backendUrl, client),
		setting.GenerateFleetPerformanceSettings(backendUrl, client),
	)
}
