package degradation_detector

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server/meta"
  "net/http"
  "sync"
  "time"
)

type InsertionResults struct {
  Degradation Degradation
  WasInserted bool
  Error       error
}

func PostDegradations(client *http.Client, backendURL string, degradations <-chan Degradation) chan InsertionResults {
  url := backendURL + "/api/meta/accidents"
  insertionResults := make(chan InsertionResults)
  go func() {
    defer close(insertionResults)
    var wg sync.WaitGroup
    for degradation := range degradations {
      wg.Add(1)
      func(degradation Degradation) {
        defer wg.Done()
        analysisSettings := degradation.analysisSettings
        date := time.UnixMilli(degradation.timestamp).UTC().Format("2006-01-02")
        medianMessage := getMessageBasedOnMedianChange(degradation.medianValues)
        kind := "InferredRegression"
        if !degradation.isDegradation {
          kind = "InferredImprovement"
        }
        insertParams := meta.AccidentInsertParams{Date: date, Test: analysisSettings.Test + "/" + analysisSettings.Metric, Kind: kind, Reason: medianMessage, BuildNumber: degradation.build}
        params, err := json.Marshal(insertParams)
        if err != nil {
          insertionResults <- InsertionResults{Error: fmt.Errorf("failed to marshal query: %w", err)}
          return
        }
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
        defer cancel()
        req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(params))
        if err != nil {
          insertionResults <- InsertionResults{Error: fmt.Errorf("failed to create request: %w", err)}
          return
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        if err != nil {
          insertionResults <- InsertionResults{Error: fmt.Errorf("failed to send POST request: %w", err)}
          return
        }
        defer resp.Body.Close()

        // the accident already exists
        if resp.StatusCode == http.StatusConflict {
          insertionResults <- InsertionResults{}
          return
        }

        if resp.StatusCode != http.StatusOK {
          insertionResults <- InsertionResults{Error: fmt.Errorf("failed to post Degradation: %v", resp.Status)}
          return
        }

        insertionResults <- InsertionResults{degradation, true, nil}
      }(degradation)
    }
    wg.Wait()
  }()

  return insertionResults
}
