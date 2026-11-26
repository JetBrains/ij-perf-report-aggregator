package degradation_detector

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server/meta"
)

// missingDataKey is what the backend stores a missing metric by. It accepts each key once.
type missingDataKey struct {
	buildType    string
	project      string
	metric       string
	missingSince int64
}

// PostMissingData stores the missing metrics and passes on the ones that the backend did not know yet.
//
// A series has one settings per receiver of its messages, and the backend accepts a missing metric once. So the
// settings of one series share one request: each of their channels gets the metric, and each channel gets it once.
func PostMissingData(client *http.Client, backendURL string, missingData <-chan MissingData) chan MissingData {
	url := backendURL + "/api/meta/missingData"
	insertionResults := make(chan MissingData, 100)
	go func() {
		defer close(insertionResults)
		type receiver struct {
			key     missingDataKey
			channel string
		}
		var wg sync.WaitGroup
		var mu sync.Mutex
		inserts := make(map[missingDataKey]func() bool)
		receivers := make(map[receiver]bool)
		for missingDatum := range missingData {
			wg.Go(func() {
				key := missingDataKey{
					buildType:    missingDatum.TCBuildType,
					project:      missingDatum.Settings.GetProject(),
					metric:       missingDatum.Settings.GetMetric(),
					missingSince: missingDatum.LastTimestamp,
				}
				to := receiver{key: key, channel: missingDatum.Settings.SlackChannel()}

				mu.Lock()
				insert, ok := inserts[key]
				if !ok {
					insert = sync.OnceValue(func() bool { return insertMissingData(client, url, key) })
					inserts[key] = insert
				}
				duplicate := receivers[to]
				receivers[to] = true
				mu.Unlock()

				if !duplicate && insert() {
					insertionResults <- missingDatum
				}
			})
		}
		wg.Wait()
	}()

	return insertionResults
}

// insertMissingData reports whether the backend stored the missing metric: false when it already has it, and when
// the request failed.
func insertMissingData(client *http.Client, url string, key missingDataKey) bool {
	insertParams := meta.MissingDataInsertParams{BuildType: key.buildType, Project: key.project, Metric: key.metric, MissingSince: key.missingSince}
	params, err := json.Marshal(insertParams)
	if err != nil {
		slog.Error("failed to marshal query", "error", err)
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(params))
	if err != nil {
		slog.Error("failed to create request", "error", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to send POST request", "error", err)
		return false
	}
	defer resp.Body.Close()

	// the missing metric already exists
	if resp.StatusCode == http.StatusConflict {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to post Details", "status", resp.Status)
		return false
	}
	return true
}
