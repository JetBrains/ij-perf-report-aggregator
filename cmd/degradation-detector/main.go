package main

import (
  "context"
  "fmt"
  "log"
  "math"
)

const URL = "http://localhost:9044"

type AnalysisSettings struct {
  test   string
  metric string
  db     string
  table  string
}

func main() {
  analysisSettings := AnalysisSettings{
    test:   "intellij_sources/localInspection/java_file",
    metric: "firstCodeAnalysis",
    db:     "perfint",
    table:  "idea",
  }

  ctx := context.Background()

  timestamps, values, builds, err := getDataFromClickhouse(ctx, analysisSettings)
  if err != nil {
    log.Printf("%v", err)
  }

  degradations := inferDegradations(values, builds, timestamps)

  insertionResults := postDegradation(ctx, analysisSettings, degradations)

  for _, result := range insertionResults {
    if result.error != nil {
      log.Printf("%v", result.error)
      continue
    }
    if !result.wasInserted {
      continue
    }
    err = sendSlackMessage(ctx, result.degradation, analysisSettings)
    if err != nil {
      log.Printf("%v", err)
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
