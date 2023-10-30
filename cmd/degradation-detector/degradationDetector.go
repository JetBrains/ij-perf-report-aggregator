package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "log"
  "math"
)

type Degradation struct {
  build        string
  timestamp    int64
  medianValues MedianValues
}

type MedianValues struct {
  previousValue float64
  newValue      float64
}

func inferDegradations(values []int, builds []string, timestamps []int64) []Degradation {

  changePoints, err := server.GetChangePointIndexes(values, 1)
  if err != nil {
    log.Fatalf("%v", err)
  }
  latestChangePoints := getChangePointsInLatestValues(changePoints, values, 20)

  segments, segmentChangePoints := getSegmentsBetweenChangePoints(latestChangePoints, values)
  degradations := make([]Degradation, 0)
  if len(segments) < 1 {
    fmt.Println("No significant change points were detected.")
    return degradations
  }
  previousMedian := server.CalculateMedian(segments[0])
  for i := 1; i < len(segments); i++ {
    currentMedian := server.CalculateMedian(segments[i])
    percentageChange := math.Abs((currentMedian - previousMedian) / previousMedian * 100)
    if math.Abs(percentageChange) > 10 {
      index := segmentChangePoints[i-1]
      build := builds[index]
      degradation := Degradation{
        build:        build,
        timestamp:    timestamps[index],
        medianValues: MedianValues{previousValue: previousMedian, newValue: currentMedian},
      }
      degradations = append(degradations, degradation)
    } else {
      fmt.Println("No significant change points detected.")
    }
    previousMedian = currentMedian
  }
  return degradations
}

func getChangePointsInLatestValues(changePoints []int, values []int, numberOfLastValuesToTake int) []int {
  filteredChanges := make([]int, 0)
  for _, change := range changePoints {
    if change >= len(values)-numberOfLastValuesToTake {
      filteredChanges = append(filteredChanges, change)
    }
  }
  return filteredChanges
}

func getSegmentsBetweenChangePoints(changePoints []int, values []int) ([][]int, []int) {
  var segments [][]int
  var segmentChangePoints []int
  prevChangePoint := 0
  for _, changePoint := range changePoints {
    segment := values[prevChangePoint:changePoint]
    segments = append(segments, segment)
    segmentChangePoints = append(segmentChangePoints, changePoint)
    prevChangePoint = changePoint
  }
  if prevChangePoint < len(values) {
    segment := values[prevChangePoint:]
    segments = append(segments, segment)
    segmentChangePoints = append(segmentChangePoints, len(values))
  }
  return segments, segmentChangePoints
}
