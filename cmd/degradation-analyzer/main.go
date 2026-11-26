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
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/alerts"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/setting"
)

func main() {
	backendUrl := getBackendUrl()
	client := createHttpClient()
	slog.Info("started")
	analysisSettings := collectSettings(loadAlertsConfig(client), backendUrl, client)
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
		forSlack := make(chan detector.DegradationWithSettings, 100)
		forHooks := make(chan detector.DegradationWithSettings, 100)
		util.Broadcast(mergedResults, forSlack, forHooks)
		var send sync.WaitGroup
		send.Go(func() { detector.SendDegradationsToSlack(forSlack, client) })
		send.Go(func() { detector.SendDegradationsToHooks(forHooks, client, setting.Hooks()) })
		send.Wait()
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

// loadAlertsConfig reads who gets which notification from the latest release of ide-infrastructure-alerts.
// GitHub being down must not cancel the run, so on a failure the run uses the config that is built into the image.
func loadAlertsConfig(client *http.Client) *alerts.Config {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	config, err := alerts.Load(ctx, client)
	if config == nil {
		slog.Error("failed to load the alerts config", "error", err)
		reportHealth(client, ":rotating_light: degradation-analyzer did not run: it failed to load the alerts config: "+err.Error())
		os.Exit(1)
	}
	if err != nil {
		slog.Error("failed to load the latest alerts config, using the one built into the image", "error", err, "version", config.Version)
		reportHealth(client, ":warning: degradation-analyzer failed to load the latest settings from ide-infrastructure-alerts: "+err.Error()+"\n"+
			"This run uses the settings of release `"+config.Version+"` that are built into the image, so later changes of the settings have no effect.")
	}
	slog.Info("alerts config loaded", "version", config.Version, "commit", config.Commit)
	return config
}

func collectSettings(config *alerts.Config, backendUrl string, client *http.Client) []detector.Settings {
	analysisSettings := make([]detector.Settings, 0, 25000)
	empty := make([]string, 0)
	backend := setting.NewBackend(backendUrl, client)
	if err := setting.CheckUnitTestOwners(backend); err != nil {
		slog.Error("failed to load the owners of the unit tests", "error", err)
		reportHealth(client, ":warning: degradation-analyzer failed to load the owners of the unit tests: "+err.Error()+"\n"+
			"In this run the messages about a unit test can go to the fallback channel instead of the channel of its owner.")
	}
	for _, g := range setting.Generators(config, backend) {
		settings := g.Generate()
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
	slog.Error("no settings generated, the corresponding subscribers get no degradation or missing-data alerts", "generators", names)
	reportHealth(client, ":rotating_light: degradation-analyzer generated no settings for: `"+strings.Join(names, "`, `")+"`\n"+
		"Those subscribers get no degradation and no missing-data alerts until this is fixed. "+
		"Check the CronJob logs and the subscriber's ij-perf.yaml in ide-infrastructure-alerts for a config or backend error.")
}

func reportHealth(client *http.Client, text string) {
	channel := os.Getenv("HEALTH_SLACK_CHANNEL")
	if channel == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := detector.SendSlackMessage(ctx, client, detector.SlackMessage{Channel: channel, Text: text})
	if err != nil {
		slog.Error("failed to send health message to slack", "error", err, "channel", channel)
	}
}
