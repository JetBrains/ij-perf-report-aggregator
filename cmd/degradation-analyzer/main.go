package main

import (
	"log/slog"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/setting"
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
	metrics := detector.FetchMetricsFromClickhouse(analysisSettings, client, backendUrl)
	metricsForDegradation := make(chan detector.QueryResultWithSettings, 5)
	metricsForMissingMetrics := make(chan detector.QueryResultWithSettings, 5)
	util.Broadcast(metrics, metricsForDegradation, metricsForMissingMetrics)

	var wg sync.WaitGroup

	wg.Go(func() {
		degradations := detector.InferDegradations(metricsForDegradation)
		insertionResults := detector.PostDegradations(client, backendUrl, degradations)
		filteredResults := detector.FilterErrors(insertionResults)
		mergedResults := detector.MergeDegradations(filteredResults)
		detector.SendDegradationsToSlack(mergedResults, client)
	})

	wg.Go(func() {
		missingData := detector.InferMissingData(metricsForMissingMetrics)
		missingData = detector.PostMissingData(client, backendUrl, missingData)
		mergedMissingData := detector.MergeMissingData(missingData)
		detector.SendMissingDataMessages(mergedMissingData, client)
	})

	wg.Wait()
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
		setting.GenerateStartupSettingsForPhpStorm(backendUrl, client),
	)
}

func generateFleetStartupSettings() []detector.FleetStartupSettings {
	return slices.Concat(
		setting.GenerateFleetStartupSettings(),
	)
}

func generatePerformanceSettings(backendUrl string, client *http.Client) []detector.PerformanceSettings {
	return slices.Concat(
		setting.GenerateIdeaIndexingSettings(backendUrl, client),
		setting.GenerateWorkspaceSettings(),
		setting.GenerateKotlinSettings(),
		setting.GenerateKotlinIdeaSettings(backendUrl, client),
		setting.GenerateMavenSettings(),
		setting.GenerateGradleSettings(),
		setting.GenerateVCSSettings(),
		setting.GeneratePhpStormSettings(backendUrl, client),
		setting.GenerateClionSettings(backendUrl, client),
		setting.GenerateAllUnitTestsSettings(backendUrl, client),
		setting.GenerateGolandPerfSettings(backendUrl, client),
		setting.GenerateRustPerfSettings(backendUrl, client),
		setting.GenerateFleetPerformanceSettings(backendUrl, client),
		setting.GenerateRubyPerfSettings(backendUrl, client),
		setting.GenerateJavaSettings(backendUrl, client),
		setting.GenerateUltimateSettings(backendUrl, client),
		setting.GenerateAIASettings(),
		setting.GenerateAIATestTokenSettings(),
		setting.GenerateKotlinBuildToolsSettings(backendUrl, client),
		setting.GenerateKotlinMultiplatformToolingSettings(backendUrl, client),
		setting.GenerateUISettings(),
	)
}
