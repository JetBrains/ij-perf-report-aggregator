package degradation_detector

import (
	"log/slog"
	"math"
	"time"

	"github.com/AndreyAkinshin/pragmastat/go/v4"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
)

type Degradation struct {
	Build string
	// PreviousBuild is the run before Build, when the analysis knows it: the last run with the old value. A step
	// analysis sets it, so the change came with the commits between the two runs.
	PreviousBuild string
	timestamp     int64
	medianValues  CenterValues
	IsDegradation bool
	// threshold is set when a threshold analysis found the degradation: newValue of medianValues is over it.
	threshold *Threshold
	// step is set when a step analysis found the degradation.
	step *StepChange
}

// Threshold is the limit that a threshold analysis compared, and the number of the latest values of its median.
type Threshold struct {
	value  float64
	window int
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
	GetIgnoreChangesBefore() time.Time

	GetAnalysisKind() AnalysisKind
	GetThresholdMode() ThresholdMode
	GetThresholdValue() float64
	GetThresholdWindowing() (window int, onCrossing bool)
	GetStepSettings() StepSettings
}

func (v CenterValues) PercentageChange() float64 {
	return math.Abs((v.newValue - v.previousValue) / v.previousValue * 100)
}

func detectDegradations(values []int, builds []string, timestamps []int64, analysisSettings analysisSettings) []Degradation {
	degradations := make([]Degradation, 0)

	switch analysisSettings.GetAnalysisKind() {
	case ThresholdAnalysis:
		return detectThresholdExceed(values, builds, timestamps, analysisSettings)
	case StepAnalysis:
		return detectSteps(values, builds, timestamps, analysisSettings.GetStepSettings(), analysisSettings.GetReportType())
	case ChangePointAnalysis:
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

// detectThresholdExceed emits a degradation when the median of the latest values of the window of
// GetThresholdWindowing crosses the configured threshold according to the selected ThresholdMode. For GreaterThan
// mode, strictly greater (>) is used; for LessThan mode, strictly less (<) is used. The previous value of the
// degradation is the median of the same number of values one run earlier. When the analysis reports only the
// crossing, a run whose previous median already crossed the threshold emits nothing.
func detectThresholdExceed(values []int, builds []string, timestamps []int64, s analysisSettings) []Degradation {
	result := make([]Degradation, 0, 1)
	window, onCrossing := s.GetThresholdWindowing()
	if len(values) < window || len(builds) == 0 || len(timestamps) == 0 {
		return result
	}
	lastIdx := len(values) - 1
	threshold := s.GetThresholdValue()
	exceeds := func(value float64) bool {
		switch s.GetThresholdMode() {
		case ThresholdGreaterThan:
			return value > threshold
		case ThresholdLessThan:
			return value < threshold
		}
		return false
	}

	current := statistic.Median(values[len(values)-window:])
	if !exceeds(current) {
		return result
	}
	// Treat exceeding threshold as degradation event; consumer can filter by ReportType if needed
	previous := threshold
	if len(values) > window {
		previous = statistic.Median(values[len(values)-window-1 : lastIdx])
		if onCrossing && exceeds(previous) {
			return result
		}
	}
	result = append(result, Degradation{
		Build:         builds[lastIdx],
		timestamp:     timestamps[lastIdx],
		medianValues:  CenterValues{previousValue: previous, newValue: current},
		IsDegradation: true,
		threshold:     &Threshold{value: threshold, window: window},
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
