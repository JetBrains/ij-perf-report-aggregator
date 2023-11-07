package degradation_detector

import (
  "fmt"
  "math"
)

type Degradation struct {
  build            string
  timestamp        int64
  medianValues     MedianValues
  analysisSettings Settings
  isDegradation    bool
}

type MedianValues struct {
  previousValue float64
  newValue      float64
}

func InferDegradations(values []int, builds []string, timestamps []int64, analysisSettings Settings) []Degradation {
  segmentSizeThreshold := 3
  degradations := make([]Degradation, 0)

  changePoints := GetChangePointIndexes(values, 1)
  segments := getSegmentsBetweenChangePoints(changePoints, values)
  if len(segments) < 2 {
    fmt.Println("No significant change points were detected.")
    return degradations
  }
  lastSegment := segments[len(segments)-1]
  if len(lastSegment) < segmentSizeThreshold {
    fmt.Println("Last segment is too short.")
    return degradations
  }

  minimumSegmentLength := analysisSettings.MinimumSegmentLength
  if minimumSegmentLength == 0 {
    minimumSegmentLength = 3
  }
  skippedSegments := 0
  for i := len(segments) - 2; i >= 0 && skippedSegments < 4; i-- {
    if len(segments[i]) < minimumSegmentLength {
      skippedSegments++
      continue
    }
    currentMedian := CalculateMedian(lastSegment)
    previousMedian := CalculateMedian(segments[i])
    percentageChange := math.Abs((currentMedian - previousMedian) / previousMedian * 100)
    absoluteChange := math.Abs(currentMedian - previousMedian)

    if absoluteChange < 10 || percentageChange < 10 {
      break
    }
    isDegradation := false
    if currentMedian > previousMedian {
      isDegradation = true
    }
    if !isDegradation && analysisSettings.DoNotReportImprovement {
      break
    }
    index := changePoints[len(segments)-2]
    degradations = append(degradations, Degradation{
      build:            builds[index],
      timestamp:        timestamps[index],
      medianValues:     MedianValues{previousValue: previousMedian, newValue: currentMedian},
      analysisSettings: analysisSettings,
      isDegradation:    isDegradation,
    })
    break
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
