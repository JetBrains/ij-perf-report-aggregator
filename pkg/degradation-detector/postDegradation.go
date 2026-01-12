package degradation_detector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server/meta"
)

type InsertionResults struct {
	Degradation DegradationWithSettings
	Error       error
}

func FilterErrors(insertionResults <-chan InsertionResults) <-chan DegradationWithSettings {
	ch := make(chan DegradationWithSettings, 100)
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

type accidentState struct {
	mu                   sync.Mutex
	posted               bool // true if POST attempt was made
	successfullyInserted bool // true if got 200 OK (newly inserted)
}

func PostDegradations(client *http.Client, backendURL string, degradations <-chan DegradationWithSettings) chan InsertionResults {
	url := backendURL + "/api/meta/accidents"
	insertionResults := make(chan InsertionResults, 100)
	go func() {
		defer close(insertionResults)
		var wg sync.WaitGroup
		accidentStates := sync.Map{} // map[string]*accidentState

		for degradation := range degradations {
			wg.Go(func() {
				d := degradation.Details
				if d.timestamp < time.Now().Add(-672*time.Hour).UnixMilli() { // Do not post degradations older than 28 days
					return
				}

				accidentKey := fmt.Sprintf("%s:%s", degradation.Settings.DBTestName(), d.Build)
				stateInterface, _ := accidentStates.LoadOrStore(accidentKey, &accidentState{})
				state, ok := stateInterface.(*accidentState)
				if !ok {
					insertionResults <- InsertionResults{Error: errors.New("unexpected type in accidentStates map")}
					return
				}

				state.mu.Lock()
				defer state.mu.Unlock()

				if state.posted {
					if state.successfullyInserted {
						insertionResults <- InsertionResults{DegradationWithSettings{d, degradation.Settings}, nil}
					}
					return
				}

				state.posted = true

				date := time.UnixMilli(d.timestamp).UTC().Format("2006-01-02")
				medianMessage := getMessageBasedOnMedianChange(d.medianValues)
				kind := "InferredRegression"
				if !d.IsDegradation {
					kind = "InferredImprovement"
				}
				insertParams := meta.AccidentInsertParams{Date: date, Test: degradation.Settings.DBTestName(), Kind: kind, Reason: medianMessage, BuildNumber: d.Build, UserName: "R2D2"}
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

				if resp.StatusCode == http.StatusOK {
					state.successfullyInserted = true
					insertionResults <- InsertionResults{DegradationWithSettings{d, degradation.Settings}, nil}
					return
				}

				// the accident already exists
				if resp.StatusCode == http.StatusConflict {
					return
				}

				insertionResults <- InsertionResults{Error: fmt.Errorf("failed to post Details: %v", resp.Status)}
			})
		}
		wg.Wait()
	}()

	return insertionResults
}
