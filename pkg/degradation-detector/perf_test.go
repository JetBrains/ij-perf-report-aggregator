package degradation_detector

import (
  "context"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/analysis"
  "log"
  "os"
  "sync"
  "testing"
)

func TestDegradationDetector(_ *testing.T) {
  backendUrl := os.Getenv("BACKEND_URL")
  if len(backendUrl) == 0 {
    backendUrl = "https://ij-perf-api.labs.jb.gg"
    log.Printf("BACKEND_URL is not set, using default value: %s", backendUrl)
  }

  analysisSettings := make([]analysis.Settings, 0, 1000)
  analysisSettings = append(analysisSettings, analysis.GenerateIdeaSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateWorkspaceSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateKotlinSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateMavenSettings()...)
  analysisSettings = append(analysisSettings, analysis.GenerateGradleSettings()...)
  analysisSettings = append(analysisSettings, analysis.GeneratePhpStormSettings()...)

  ctx := context.Background()
  degradationsChan := make(chan []Degradation)
  var wg sync.WaitGroup

  // Create a semaphore with a capacity of 8.
  semaphore := make(chan struct{}, 8)

  for _, analysisSetting := range analysisSettings {
    wg.Add(1)
    go func(as analysis.Settings) {
      defer wg.Done()
      // Acquire a slot in the semaphore before proceeding.
      semaphore <- struct{}{}
      log.Printf("Processing %v", as)
      timestamps, values, builds, err := GetDataFromClickhouse(ctx, backendUrl, as)
      if err != nil {
        log.Printf("%v", err)
        degradationsChan <- nil // or handle the Error differently
      } else {
        degradationsChan <- InferDegradations(values, builds, timestamps, as)
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
