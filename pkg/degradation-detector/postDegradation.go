package degradation_detector

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server/meta"
  "net/http"
  "time"
)

type InsertionResults struct {
  Degradation Degradation
  WasInserted bool
  Error       error
}

func PostDegradations(ctx context.Context, client *http.Client, backendURL string, degradations []Degradation) []InsertionResults {
  url := backendURL + "/api/meta/accidents"
  insertionResults := make([]InsertionResults, len(degradations))
  for i, degradation := range degradations {
    func(degradation Degradation, i int) {
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
        insertionResults[i] = InsertionResults{Error: fmt.Errorf("failed to marshal query: %w", err)}
        return
      }

      req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(params))
      if err != nil {
        insertionResults[i] = InsertionResults{Error: fmt.Errorf("failed to create request: %w", err)}
        return
      }
      req.Header.Set("Content-Type", "application/json")

      resp, err := client.Do(req)
      if err != nil {
        insertionResults[i] = InsertionResults{Error: fmt.Errorf("failed to send POST request: %w", err)}
        return
      }
      defer resp.Body.Close()

      if resp.StatusCode != http.StatusOK {
        insertionResults[i] = InsertionResults{Error: fmt.Errorf("failed to post Degradation: %v", resp.Status)}
        return
      }

      // the accident already exists
      if resp.StatusCode == http.StatusConflict {
        insertionResults[i] = InsertionResults{}
        return
      }
      insertionResults[i] = InsertionResults{degradation, true, nil}
    }(degradation, i)
  }
  return insertionResults
}
