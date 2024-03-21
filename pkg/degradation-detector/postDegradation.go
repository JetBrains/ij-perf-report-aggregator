package degradation_detector

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server/meta"
  "log/slog"
  "net/http"
  "sync"
  "time"
)

type InsertionResults struct {
  Degradation DegradationWithSettings
  Error       error
}

func FilterErrors(insertionResults <-chan InsertionResults) <-chan DegradationWithSettings {
  ch := make(chan DegradationWithSettings)
  go func() {
    for result := range insertionResults {
      if result.Error != nil {
        slog.Error("error while inserting degradation", "error", result.Error, "degradation", result.Degradation)
        continue
      }
      ch <- result.Degradation
    }
    close(ch)
  }()
  return ch
}

type accidentWriter interface {
  DBTestName() string
}

func (s PerformanceSettings) DBTestName() string {
  return s.Project + "/" + s.Metric
}

func (s StartupSettings) DBTestName() string {
  return s.Product + "/" + s.Project + "/" + s.Metric
}

func (s FleetStartupSettings) DBTestName() string {
  return "fleet" + "/" + s.Metric
}

func PostDegradations(client *http.Client, backendURL string, degradations <-chan DegradationWithSettings) chan InsertionResults {
  url := backendURL + "/api/meta/accidents"
  insertionResults := make(chan InsertionResults)
  go func() {
    defer close(insertionResults)
    var wg sync.WaitGroup
    for degradation := range degradations {
      wg.Add(1)
      func() {
        defer wg.Done()
        d := degradation.Details
        date := time.UnixMilli(d.timestamp).UTC().Format("2006-01-02")
        medianMessage := getMessageBasedOnMedianChange(d.medianValues)
        kind := "InferredRegression"
        if !d.IsDegradation {
          kind = "InferredImprovement"
        }
        insertParams := meta.AccidentInsertParams{Date: date, Test: degradation.Settings.DBTestName(), Kind: kind, Reason: medianMessage, BuildNumber: d.Build}
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
          return
        }

        if resp.StatusCode != http.StatusOK {
          insertionResults <- InsertionResults{Error: fmt.Errorf("failed to post Details: %v", resp.Status)}
          return
        }

        insertionResults <- InsertionResults{DegradationWithSettings{d, degradation.Settings}, nil}
      }()
    }
    wg.Wait()
  }()

  return insertionResults
}
