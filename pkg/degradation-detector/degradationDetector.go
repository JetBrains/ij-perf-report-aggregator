package degradation_detector

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
  "log/slog"
  "math"
)

type Degradation struct {
  Build         string
  timestamp     int64
  medianValues  MedianValues
  IsDegradation bool
}

type MedianValues struct {
  previousValue float64
  newValue      float64
}

type analysisSettings interface {
  GetDoNotReportImprovement() bool
  GetMinimumSegmentLength() int
  GetMedianDifferenceThreshold() float64
  GetEffectSizeThreshold() float64
}

func (v MedianValues) PercentageChange() float64 {
  return math.Abs((v.newValue - v.previousValue) / v.previousValue * 100)
}

func detectDegradations(values []int, builds []string, timestamps []int64, analysisSettings analysisSettings) []Degradation {
  degradations := make([]Degradation, 0)

  minimumSegmentLength := analysisSettings.GetMinimumSegmentLength()
  if minimumSegmentLength == 0 {
    minimumSegmentLength = 5
  }
  medianDifference := analysisSettings.GetMedianDifferenceThreshold()
  if medianDifference == 0 {
    medianDifference = 10
  }

  effectSizeThreshold := analysisSettings.GetEffectSizeThreshold()
  if effectSizeThreshold == 0 {
    effectSizeThreshold = 2
  }

  changePoints := statistic.GetChangePointIndexes(values, statistic.Min(minimumSegmentLength, 5))
  segments := getSegmentsBetweenChangePoints(changePoints, values)
  if len(segments) < 2 {
    slog.Debug("no significant change points were detected")
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
    currentMedian := statistic.Median(lastSegment)
    previousMedian := statistic.Median(segments[i])
    percentageChange := math.Abs((currentMedian - previousMedian) / previousMedian * 100)
    absoluteChange := math.Abs(currentMedian - previousMedian)

    if absoluteChange < 10 || percentageChange < medianDifference {
      break
    }

    es := statistic.EffectSize(lastSegment, segments[i])
    if es < effectSizeThreshold {
      break
    }

    isDegradation := false
    if currentMedian > previousMedian {
      isDegradation = true
    }
    if !isDegradation && analysisSettings.GetDoNotReportImprovement() {
      break
    }
    index := changePoints[len(segments)-2]
    degradations = append(degradations, Degradation{
      Build:         builds[index],
      timestamp:     timestamps[index],
      medianValues:  MedianValues{previousValue: previousMedian, newValue: currentMedian},
      IsDegradation: isDegradation,
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
