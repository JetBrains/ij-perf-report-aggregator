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
}

func (v CenterValues) PercentageChange() float64 {
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

		ratio, err := pragmastat.Ratio(lastSegment, segments[i])
		if err != nil {
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

		percentageChange := math.Abs((ratio - 1) * 100)
		absoluteChange := math.Abs(currentCenter - previousCenter)

		if absoluteChange < 10 || percentageChange < medianDifference {
			break
		}

		es := statistic.EffectSize(lastSegment, segments[i])
		if es < effectSizeThreshold {
			break
		}

		isDegradation := ratio > 1
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
