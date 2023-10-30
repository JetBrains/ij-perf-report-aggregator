package main

import (
  "context"
  "fmt"
  "log"
  "math"
  "os"
)

type AnalysisSettings struct {
  test    string
  metric  string
  db      string
  table   string
  channel string
}

func main() {
  backendUrl := os.Getenv("BACKEND_URL")
  if len(backendUrl) == 0 {
    log.Printf("BACKEND_URL is not set, using default value: %s", backendUrl)
    backendUrl = "http://localhost:9044"
  }

  analysisSettings := generateIdeaAnalysisSettings()
  for _, analysisSetting := range analysisSettings {
    ctx := context.Background()

    timestamps, values, builds, err := getDataFromClickhouse(ctx, backendUrl, analysisSetting)
    if err != nil {
      log.Printf("%v", err)
    }

    degradations := inferDegradations(values, builds, timestamps)

    insertionResults := postDegradation(ctx, backendUrl, analysisSetting, degradations)

    for _, result := range insertionResults {
      if result.error != nil {
        log.Printf("%v", result.error)
        continue
      }
      if !result.wasInserted {
        continue
      }
      err = sendSlackMessage(ctx, result.degradation, analysisSetting)
      if err != nil {
        log.Printf("%v", err)
      }
    }
  }
}

func getMessageBasedOnMedianChange(medianValues MedianValues) string {
  percentageChange := math.Abs((medianValues.newValue - medianValues.previousValue) / medianValues.previousValue * 100)
  medianMessage := fmt.Sprintf("Median changed by: %.2f%%. Median was %.2f and now it is %.2f.", percentageChange, medianValues.previousValue, medianValues.newValue)
  if medianValues.newValue > medianValues.previousValue {
    return "Degradation detected. " + medianMessage
  }
  return "Improvement detected. " + medianMessage
}
