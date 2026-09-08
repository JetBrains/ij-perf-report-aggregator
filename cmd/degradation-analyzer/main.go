package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
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
	analysisSettings := collectSettings(backendUrl, client)
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

func collectSettings(backendUrl string, client *http.Client) []detector.Settings {
	analysisSettings := make([]detector.Settings, 0, 25000)
	empty := make([]string, 0)
	for _, g := range setting.Generators() {
		settings := g.Generate(backendUrl, client)
		slog.Info("settings generated", "generator", g.Name, "count", len(settings))
		if len(settings) == 0 {
			empty = append(empty, g.Name)
		}
		analysisSettings = append(analysisSettings, settings...)
	}
	if len(empty) > 0 {
		reportEmptyGenerators(client, empty)
	}
	return analysisSettings
}

func reportEmptyGenerators(client *http.Client, names []string) {
	slog.Error("no settings generated, the corresponding products get no degradation or missing-data alerts", "generators", names)
	channel := os.Getenv("HEALTH_SLACK_CHANNEL")
	if channel == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := detector.SendSlackMessage(ctx, client, detector.SlackMessage{
		Channel: channel,
		Text: ":rotating_light: degradation-analyzer generated no settings for: `" + strings.Join(names, "`, `") + "`\n" +
			"Those products get no degradation and no missing-data alerts until this is fixed. " +
			"Check the CronJob logs or the settings generator for a config or backend error.",
	})
	if err != nil {
		slog.Error("failed to send health message to slack", "error", err, "channel", channel)
	}
}
