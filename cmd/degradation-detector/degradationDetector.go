package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "log"
  "math"
)

type Degradation struct {
  build            string
  timestamp        int64
  medianValues     MedianValues
  analysisSettings AnalysisSettings
  isDegradation    bool
}

type MedianValues struct {
  previousValue float64
  newValue      float64
}

func inferDegradations(values []int, builds []string, timestamps []int64, analysisSettings AnalysisSettings) []Degradation {
  numberOfLastValuesToTake := 40

  changePoints, err := server.GetChangePointIndexes(values, 1)
  if err != nil {
    log.Fatalf("%v", err)
  }

  segments := getSegmentsBetweenChangePoints(changePoints, values)
  degradations := make([]Degradation, 0)
  if len(segments) < 1 {
    fmt.Println("No significant change points were detected.")
    return degradations
  }
  previousMedian := server.CalculateMedian(segments[0])
  for i := 1; i < len(segments); i++ {
    currentMedian := server.CalculateMedian(segments[i])
    percentageChange := math.Abs((currentMedian - previousMedian) / previousMedian * 100)
    index := changePoints[i-1]
    isLatestChangePoint := index >= len(values)-numberOfLastValuesToTake
    if percentageChange > 10 && isLatestChangePoint {
      build := builds[index]
      isDegradation := false
      if currentMedian > previousMedian {
        isDegradation = true
      }
      degradation := Degradation{
        build:            build,
        timestamp:        timestamps[index],
        medianValues:     MedianValues{previousValue: previousMedian, newValue: currentMedian},
        analysisSettings: analysisSettings,
        isDegradation:    isDegradation,
      }
      degradations = append(degradations, degradation)
    }
    previousMedian = currentMedian
  }
  return degradations
}

func getSegmentsBetweenChangePoints(changePoints []int, values []int) [][]int {
  var segments [][]int
  prevChangePoint := 0
  for _, changePoint := range changePoints {
    segment := values[prevChangePoint:changePoint]
    segments = append(segments, segment)
    prevChangePoint = changePoint
  }
  if prevChangePoint < len(values) {
    segment := values[prevChangePoint:]
    segments = append(segments, segment)
  }
  return segments
}
