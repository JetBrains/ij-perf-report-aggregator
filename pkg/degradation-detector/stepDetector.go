package degradation_detector

import (
	"math"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
)

// StepSettings configure a step analysis. A run of a test that measures each scenario several times reports the median
// of its passes, so one value is already one stable measurement, and the analysis can judge it on the night it comes.
type StepSettings struct {
	// Relative is the smallest change that counts, in percent of the baseline. The default value is 20 %.
	Relative float64
	// Absolute is the smallest change that counts, in the unit of the metric. The default value is 10.
	Absolute float64
	// Noise is the number of robust standard deviations of the baseline that a change must exceed. A noisy metric
	// then needs a larger change than a stable one. The default value is 4.
	Noise float64
	// Baseline is the number of the values before a run that make its baseline. The default value is 10.
	Baseline int
}

const (
	// stepMinimumBaseline is the number of values that a baseline needs, so that its median and its spread mean
	// something. A new metric gets no report until it has them.
	stepMinimumBaseline = 5
	// stepLookback is the number of the latest runs whose steps the analysis reports. A step stays reported once: the
	// detector stores each degradation as an accident of its build and sends only the new ones. The lookback finds a
	// step also when the detector did not run on the night of the step.
	stepLookback = 3
	// stepDriftWindow is the number of the latest runs whose median the drift check compares to the baseline before
	// them. It finds a slow growth whose each step is too small for the step check.
	stepDriftWindow = 5
	// madToSigma turns the median absolute deviation into the standard deviation of a normal distribution.
	madToSigma = 1.4826
)

func (s StepSettings) withDefaults() StepSettings {
	if s.Relative == 0 {
		s.Relative = 20
	}
	if s.Absolute == 0 {
		s.Absolute = 10
	}
	if s.Noise == 0 {
		s.Noise = 4
	}
	if s.Baseline == 0 {
		s.Baseline = 10
	}
	return s
}

// StepKind tells which check of a step analysis found a change.
type StepKind int

const (
	// StepNight is a run whose value is out of the band of the baseline, while the run before it was in its own band.
	StepNight StepKind = iota
	// StepDrift is a window of runs whose median is out of the band of the baseline before them, while no run of the
	// window was a step.
	StepDrift
)

// StepChange is what a step analysis found: the check, the band and the number of the runs of the baseline.
type StepChange struct {
	kind     StepKind
	baseline int
	margin   float64
}

type band struct {
	center float64
	margin float64
	size   int
}

// bandOf returns the band of the baseline of the given values, or false when there are too few of them.
func (s StepSettings) bandOf(baseline []int) (band, bool) {
	if len(baseline) < stepMinimumBaseline {
		return band{}, false
	}
	center := statistic.Median(baseline)
	deviations := make([]int, len(baseline))
	for i, value := range baseline {
		deviations[i] = int(math.Round(math.Abs(float64(value) - center)))
	}
	spread := statistic.Median(deviations) * madToSigma
	margin := max(s.Relative/100*math.Abs(center), s.Absolute, s.Noise*spread)
	return band{center: center, margin: margin, size: len(baseline)}, true
}

// direction returns 1 when value is over the band, -1 when it is under it, and 0 when it is in it.
func (b band) direction(value float64) int {
	switch {
	case value > b.center+b.margin:
		return 1
	case value < b.center-b.margin:
		return -1
	}
	return 0
}

// stepAt tells whether the run at index is out of the band of the runs before it, and returns the band.
func (s StepSettings) stepAt(values []int, index int) (int, band) {
	if index < 1 {
		return 0, band{}
	}
	b, ok := s.bandOf(values[max(0, index-s.Baseline):index])
	if !ok {
		return 0, band{}
	}
	return b.direction(float64(values[index])), b
}

func wanted(direction int, reportType ReportType) bool {
	switch reportType {
	case DegradationEvent:
		return direction > 0
	case ImprovementEvent:
		return direction < 0
	case AllEvent:
		return direction != 0
	}
	return false
}

// newStepAt tells whether the run at index is a new step, and returns the band that it left. A run out of the band of
// its baseline is a new step, unless it continues an earlier step: the latest run of its baseline that was out of its
// own band in the same direction, with every run since then near the level of that run. Near is the margin of the band
// that the earlier run left, scaled to its level, so a noisy metric keeps its wide band at the new level, and at least
// the relative or the absolute change of the settings. Then the run is a new step only when it also left that level in the same
// direction. So a level after a step is not reported again, also not after one quiet night near the band edge, and a
// second step soon after the first is still reported. A run since the earlier step that is far from its level ends
// that step: it was one outlier night, and the run is judged against its baseline alone.
func (s StepSettings) newStepAt(values []int, index int) (int, band) {
	direction, b := s.stepAt(values, index)
	if direction == 0 {
		return 0, b
	}
	latest := -1
	var left band
	for earlier := index - 1; earlier >= max(1, index-s.Baseline); earlier-- {
		if before, earlierBand := s.stepAt(values, earlier); before == direction {
			latest, left = earlier, earlierBand
			break
		}
	}
	if latest < 0 {
		return direction, b
	}
	level := s.levelOf(float64(values[latest]), left)
	for _, value := range values[latest+1 : index] {
		if level.direction(float64(value)) != 0 {
			return direction, b
		}
	}
	if level.direction(float64(values[index])) == direction {
		return direction, level
	}
	return 0, b
}

// levelOf returns the band of the level of one run that left the band left: the margin of left scaled to the run, and at
// least the relative and the absolute change of the settings.
func (s StepSettings) levelOf(value float64, left band) band {
	margin := max(s.Relative/100*math.Abs(value), s.Absolute)
	if left.center != 0 {
		margin = max(margin, left.margin/math.Abs(left.center)*math.Abs(value))
	}
	return band{center: value, margin: margin, size: 1}
}

// detectSteps finds the changes of the latest runs. It reports a run as a step on the night it comes, see newStepAt. So
// the build of the report is the first run with the new value, and the run before it is the last run with the old one.
//
// When no run of the latest runs was a step, it compares the median of those runs to the baseline before them, and
// reports a drift at the run where that median left the band. A drift needs a window and a baseline with no step, so
// a step is never reported again as a drift.
func detectSteps(values []int, builds []string, timestamps []int64, s StepSettings, reportType ReportType) []Degradation {
	s = s.withDefaults()
	result := make([]Degradation, 0, 1)
	last := len(values) - 1
	if last < 1 || len(builds) != len(values) || len(timestamps) != len(values) {
		return result
	}

	for index := max(1, last-stepLookback+1); index <= last; index++ {
		direction, b := s.newStepAt(values, index)
		if !wanted(direction, reportType) {
			continue
		}
		result = append(result, Degradation{
			Build:         builds[index],
			PreviousBuild: builds[index-1],
			timestamp:     timestamps[index],
			medianValues:  CenterValues{previousValue: b.center, newValue: float64(values[index])},
			IsDegradation: direction > 0,
			step:          &StepChange{kind: StepNight, baseline: b.size, margin: b.margin},
		})
	}
	if len(result) > 0 {
		return result
	}

	if drift, ok := s.driftAt(values, last, reportType); ok {
		start := last - stepDriftWindow + 1
		result = append(result, Degradation{
			Build:         builds[last],
			PreviousBuild: builds[start-1],
			timestamp:     timestamps[last],
			medianValues:  CenterValues{previousValue: drift.center, newValue: statistic.Median(values[start:])},
			IsDegradation: statistic.Median(values[start:]) > drift.center,
			step:          &StepChange{kind: StepDrift, baseline: drift.size, margin: drift.margin},
		})
	}
	return result
}

// driftAt tells whether the median of the window that ends at index left the band of the baseline before the window at
// that run, with no step inside the window or its baseline.
func (s StepSettings) driftAt(values []int, index int, reportType ReportType) (band, bool) {
	isDrift := func(end int) (band, bool) {
		start := end - stepDriftWindow + 1
		if start < 1 {
			return band{}, false
		}
		for i := max(1, start-s.Baseline); i <= end; i++ {
			if direction, _ := s.stepAt(values, i); direction != 0 {
				return band{}, false
			}
		}
		b, ok := s.bandOf(values[max(0, start-s.Baseline):start])
		if !ok {
			return band{}, false
		}
		return b, wanted(b.direction(statistic.Median(values[start:end+1])), reportType)
	}
	b, ok := isDrift(index)
	if !ok {
		return band{}, false
	}
	if _, before := isDrift(index - 1); before {
		return band{}, false
	}
	return b, true
}
