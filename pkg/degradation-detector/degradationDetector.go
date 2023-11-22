package degradation_detector

import (
  "log/slog"
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
  degradations := make([]Degradation, 0)

  minimumSegmentLength := analysisSettings.MinimumSegmentLength
  if minimumSegmentLength == 0 {
    minimumSegmentLength = 5
  }
  medianDifference := analysisSettings.MedianDifferenceThreshold
  if medianDifference == 0 {
    medianDifference = 10
  }

  changePoints := GetChangePointIndexes(values, 1)
  segments := getSegmentsBetweenChangePoints(changePoints, values)
  if len(segments) < 2 {
    slog.Info("no significant change points were detected")
    return degradations
  }
  lastSegment := segments[len(segments)-1]
  if len(lastSegment) < minimumSegmentLength {
    slog.Info("last segment is too short")
    return degradations
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

    if absoluteChange < 10 || percentageChange < medianDifference {
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
