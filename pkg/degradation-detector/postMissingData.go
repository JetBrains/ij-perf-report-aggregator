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

func PostMissingData(client *http.Client, backendURL string, missingData <-chan MissingData) chan MissingData {
	url := backendURL + "/api/meta/missingData"
	insertionResults := make(chan MissingData, 100)
	go func() {
		defer close(insertionResults)
		var wg sync.WaitGroup
		for missingDatum := range missingData {
			wg.Go(func() {
				insertParams := meta.MissingDataInsertParams{BuildType: missingDatum.TCBuildType, Project: missingDatum.Settings.GetProject(), Metric: missingDatum.Settings.GetMetric(), MissingSince: missingDatum.LastTimestamp}
				params, err := json.Marshal(insertParams)
				if err != nil {
					slog.Error("failed to marshal query", "error", err)
					return
				}
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				defer cancel()
				req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(params))
				if err != nil {
					slog.Error("failed to create request", "error", err)
					return
				}
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					slog.Error("failed to send POST request", "error", err)
					return
				}
				defer resp.Body.Close()

				// the accident already exists
				if resp.StatusCode == http.StatusConflict {
					return
				}

				if resp.StatusCode != http.StatusOK {
					slog.Error("failed to post Details", "status", resp.Status)
					return
				}

				insertionResults <- missingDatum
			})
		}
		wg.Wait()
	}()

	return insertionResults
}
