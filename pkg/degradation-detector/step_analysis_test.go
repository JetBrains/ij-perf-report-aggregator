package degradation_detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stepAnalysis() AnalysisSettings {
	return AnalysisSettings{ReportType: DegradationEvent, AnalysisKind: StepAnalysis, Step: StepSettings{Relative: 20, Absolute: 30}}
}

// open.M.coverRemovedMs of the Air performance tests, one value a night, 2026-11-03 to 2026-11-19 03:28.
var coverRemovedBefore = []int{2663, 2372, 2506, 2438, 2422, 2687, 2417, 2513, 2487, 2762, 2389, 2494, 2478, 2386, 2317, 2370, 2347, 2316, 2399}

func TestStep_ReportsTheFirstSlowNight(t *testing.T) {
	t.Parallel()
	// On 2026-11-19 22:19 the first click of the session row stopped to open the session, and the test clicked again
	// after 10 s.
	nights := append(append([]int{}, coverRemovedBefore...), 12930)
	degradations := detectOverNights(nights, stepAnalysis())
	require.Len(t, degradations, 1)
	d := degradations[0]
	assert.Equal(t, "b19", d.Build)
	assert.Equal(t, "b18", d.PreviousBuild)
	// The baseline is the median of the 10 nights before.
	assert.InDelta(t, 2387.5, d.medianValues.previousValue, 0.0001)
	assert.InDelta(t, 12930, d.medianValues.newValue, 0.0001)
	assert.Equal(t, "step", analysisName(d))
	assert.Contains(t, getReason(d), "Changed between builds b18 and b19")
}

func TestStep_ReportsTheStepOnTheNextRunsToo(t *testing.T) {
	t.Parallel()
	// The detector can miss the night of the step. The step stays the same accident, so a later run sends nothing new.
	nights := append(append([]int{}, coverRemovedBefore...), 12930, 12939)
	degradations := detectOverNights(nights, stepAnalysis())
	require.Len(t, degradations, 1)
	assert.Equal(t, "b19", degradations[0].Build)
}

func TestStep_DoesNotReportANewLevelAgain(t *testing.T) {
	t.Parallel()
	nights := append(append([]int{}, coverRemovedBefore...), 12930, 12939, 12678, 12588)
	assert.Empty(t, detectOverNights(nights, stepAnalysis()))
}

func TestStep_IgnoresTheNoiseOfAStableMetric(t *testing.T) {
	t.Parallel()
	for end := stepMinimumBaseline + 1; end <= len(coverRemovedBefore); end++ {
		assert.Empty(t, detectOverNights(coverRemovedBefore[:end], stepAnalysis()), "the first %d nights", end)
	}
}

func TestStep_NeedsMoreNoiseThanItsBandForANoisyMetric(t *testing.T) {
	t.Parallel()
	// A metric that moves by about 30 % from night to night: its band is wider than 20 %.
	noisy := []int{400, 520, 380, 610, 450, 390, 560, 470, 420, 600}
	assert.Empty(t, detectOverNights(append(append([]int{}, noisy...), 650), stepAnalysis()))
	assert.Len(t, detectOverNights(append(append([]int{}, noisy...), 1100), stepAnalysis()), 1)
}

func TestStep_IgnoresAChangeUnderTheAbsoluteMinimum(t *testing.T) {
	t.Parallel()
	// A metric of a few milliseconds doubles without a user noticing it.
	assert.Empty(t, detectOverNights([]int{10, 11, 10, 12, 11, 10, 11, 35}, stepAnalysis()))
}

func TestStep_NeedsABaseline(t *testing.T) {
	t.Parallel()
	assert.Empty(t, detectOverNights([]int{500, 510, 505, 498, 5000}, stepAnalysis()))
	assert.Len(t, detectOverNights([]int{500, 510, 505, 498, 502, 5000}, stepAnalysis()), 1)
}

func TestStep_ReportsNoImprovementForADegradationEvent(t *testing.T) {
	t.Parallel()
	nights := append(append([]int{}, coverRemovedBefore...), 900)
	assert.Empty(t, detectOverNights(nights, stepAnalysis()))

	all := stepAnalysis()
	all.ReportType = AllEvent
	degradations := detectOverNights(nights, all)
	require.Len(t, degradations, 1)
	assert.False(t, degradations[0].IsDegradation)
}

func TestStep_ReportsASlowDriftOnce(t *testing.T) {
	t.Parallel()
	// Each night is 3 % slower than the night before, so no night alone leaves the band of 20 %.
	nights := make([]int, 0, 20)
	nights = append(nights, 1000, 1010, 990, 1005, 995, 1000, 1008, 992, 1003, 997)
	value := 1000.0
	for range 10 {
		value *= 1.03
		nights = append(nights, int(value))
	}
	// The detector runs after each night and stores a report of one build once.
	reports := map[string]Degradation{}
	for end := 11; end <= len(nights); end++ {
		for _, d := range detectOverNights(nights[:end], stepAnalysis()) {
			reports[d.Build] = d
		}
	}
	require.Len(t, reports, 1)
	for _, d := range reports {
		assert.Equal(t, "drift", analysisName(d))
	}
}
