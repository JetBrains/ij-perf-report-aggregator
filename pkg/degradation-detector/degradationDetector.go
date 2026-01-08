package degradation_detector

import (
	"log/slog"
	"math"

	"github.com/AndreyAkinshin/pragmastat/go/v3"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
)

type Degradation struct {
	Build         string
	timestamp     int64
	medianValues  CenterValues
	IsDegradation bool
}

type CenterValues struct {
	previousValue float64
	newValue      float64
}

type analysisSettings interface {
	GetReportType() ReportType
	GetMinimumSegmentLength() int
	GetMedianDifferenceThreshold() float64
	GetEffectSizeThreshold() float64
	GetDaysToCheckMissing() int

	GetAnalysisKind() AnalysisKind
	GetThresholdMode() ThresholdMode
	GetThresholdValue() float64
}

func (v CenterValues) PercentageChange() float64 {
	return math.Abs((v.newValue - v.previousValue) / v.previousValue * 100)
}

func detectDegradations(values []int, builds []string, timestamps []int64, analysisSettings analysisSettings) []Degradation {
	degradations := make([]Degradation, 0)

	if analysisSettings.GetAnalysisKind() == ThresholdAnalysis {
		return detectThresholdExceed(values, builds, timestamps, analysisSettings)
	}

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

	changePoints := statistic.GetChangePointIndexes(values, min(5, len(values)/2))
	segments := GetSegmentsBetweenChangePoints(changePoints, values)
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

		currentCenter, err := pragmastat.Center(lastSegment)
		if err != nil {
			skippedSegments++
			continue
		}
		previousCenter, err := pragmastat.Center(segments[i])
		if err != nil {
			skippedSegments++
			continue
		}

		ratio := currentCenter / previousCenter

		percentageChange := math.Abs((ratio - 1) * 100)
		absoluteChange := math.Abs(currentCenter - previousCenter)

		if absoluteChange < 10 || percentageChange < medianDifference {
			break
		}

		es := statistic.EffectSize(lastSegment, segments[i])
		if es < effectSizeThreshold {
			break
		}

		isDegradation := currentCenter > previousCenter
		reportType := analysisSettings.GetReportType()

		if !isDegradation && reportType == DegradationEvent {
			break
		}
		if isDegradation && reportType == ImprovementEvent {
			break
		}
		index := changePoints[len(segments)-2]
		degradations = append(degradations, Degradation{
			Build:         builds[index],
			timestamp:     timestamps[index],
			medianValues:  CenterValues{previousValue: previousCenter, newValue: currentCenter},
			IsDegradation: isDegradation,
		})
		break
	}

	return degradations
}

// detectThresholdExceed emits a degradation when the latest value crosses the configured threshold
// according to the selected ThresholdMode. For GreaterThan mode, strictly greater (>) is used; for
// LessThan mode, strictly less (<) is used.
func detectThresholdExceed(values []int, builds []string, timestamps []int64, s analysisSettings) []Degradation {
	result := make([]Degradation, 0, 1)
	if len(values) == 0 || len(builds) == 0 || len(timestamps) == 0 {
		return result
	}
	lastIdx := len(values) - 1
	last := float64(values[lastIdx])
	threshold := s.GetThresholdValue()

	meets := false
	switch s.GetThresholdMode() {
	case ThresholdGreaterThan:
		meets = last > threshold
	case ThresholdLessThan:
		meets = last < threshold
	}
	if !meets {
		return result
	}
	// Treat exceeding threshold as degradation event; consumer can filter by ReportType if needed
	var previous float64
	if lastIdx > 0 {
		previous = float64(values[lastIdx-1])
	} else {
		previous = threshold
	}
	result = append(result, Degradation{
		Build:         builds[lastIdx],
		timestamp:     timestamps[lastIdx],
		medianValues:  CenterValues{previousValue: previous, newValue: last},
		IsDegradation: true,
	})
	return result
}

func GetSegmentsBetweenChangePoints(changePoints []int, values []int) [][]int {
	segments := make([][]int, 0, len(changePoints)+1)
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
