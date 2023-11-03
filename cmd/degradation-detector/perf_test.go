package main

import (
  "context"
  "fmt"
  "log"
  "os"
  "sync"
  "testing"
)

func TestChangeDetector(_ *testing.T) {
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
  degradationsChan := make(chan []Degradation)
  var wg sync.WaitGroup

  // Create a semaphore with a capacity of 8.
  semaphore := make(chan struct{}, 8)

  for _, analysisSetting := range analysisSettings {
    wg.Add(1)
    go func(as AnalysisSettings) {
      defer wg.Done()
      // Acquire a slot in the semaphore before proceeding.
      semaphore <- struct{}{}
      log.Printf("Processing %v", as)
      timestamps, values, builds, err := getDataFromClickhouse(ctx, backendUrl, as)
      if err != nil {
        log.Printf("%v", err)
        degradationsChan <- nil // or handle the error differently
      } else {
        degradationsChan <- inferDegradations(values, builds, timestamps, as)
      }
      // Release the slot when finished.
      <-semaphore
    }(analysisSetting)
  }

  // Start a goroutine to close the channel once all goroutines have finished.
  go func() {
    wg.Wait()
    close(degradationsChan)
  }()

  // Collect results from the channel.
  degradations := make([]Degradation, 0, 1000)
  for ds := range degradationsChan {
    if ds != nil {
      degradations = append(degradations, ds...)
    }
  }
  fmt.Println(degradations)
}
