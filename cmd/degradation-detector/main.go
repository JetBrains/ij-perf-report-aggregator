package main

import (
  "context"
  "fmt"
  "log"
  "math"
  "os"
  "sync"
)

type AnalysisSettings struct {
  test        string
  metric      string
  db          string
  table       string
  branch      string
  machine     string
  channel     string
  productLink string
}

func main() {
  backendUrl := os.Getenv("BACKEND_URL")
  if len(backendUrl) == 0 {
    backendUrl = "https://ij-perf-api.labs.jb.gg"
    log.Printf("BACKEND_URL is not set, using default value: %s", backendUrl)
  }

  analysisSettings := make([]AnalysisSettings, 0, 1000)
  analysisSettings = append(analysisSettings, generateIdeaAnalysisSettings()...)
  analysisSettings = append(analysisSettings, generateWorkspaceAnalysisSettings()...)
  analysisSettings = append(analysisSettings, generateKotlinAnalysisSettings()...)
  analysisSettings = append(analysisSettings, generateMavenAnalysisSettings()...)
  analysisSettings = append(analysisSettings, generateGradleAnalysisSettings()...)
  analysisSettings = append(analysisSettings, generatePhpStormAnalysisSettings()...)

  ctx := context.Background()
  degradations := make([]Degradation, 0, 1000)
  for _, analysisSetting := range analysisSettings {
    log.Printf("Processing %v", analysisSetting)
    timestamps, values, builds, err := getDataFromClickhouse(ctx, backendUrl, analysisSetting)
    if err != nil {
      log.Printf("%v", err)
    }

    degradations = append(degradations, inferDegradations(values, builds, timestamps, analysisSetting)...)
  }

  insertionResults := postDegradation(ctx, backendUrl, degradations)

  var wg sync.WaitGroup
  for _, result := range insertionResults {
    if result.error != nil {
      log.Printf("%v", result.error)
      continue
    }
    if !result.wasInserted {
      continue
    }
    wg.Add(1)
    go func(result InsertionResults) {
      defer wg.Done()
      err := sendSlackMessage(ctx, result.degradation)
      if err != nil {
        log.Printf("%v", err)
      }
    }(result)
  }
  wg.Wait()
}

func getMessageBasedOnMedianChange(medianValues MedianValues) string {
  percentageChange := math.Abs((medianValues.newValue - medianValues.previousValue) / medianValues.previousValue * 100)
  medianMessage := fmt.Sprintf("Median changed by: %.2f%%. Median was %.2f and now it is %.2f.", percentageChange, medianValues.previousValue, medianValues.newValue)
  if medianValues.newValue > medianValues.previousValue {
    return "Degradation detected. " + medianMessage
  }
  return "Improvement detected. " + medianMessage
}