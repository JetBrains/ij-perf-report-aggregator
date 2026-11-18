package degradation_detector

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The nights of a metric of the Air performance tests, one run a night, and the detection over them.
func detectOverNights(values []int, settings AnalysisSettings) []Degradation {
	builds := make([]string, len(values))
	times := make([]int64, len(values))
	for i := range values {
		builds[i] = "b" + strconv.Itoa(i)
		times[i] = int64(i)
	}
	return detectDegradations(values, builds, times, settings)
}

func windowedLimit(limit float64) AnalysisSettings {
	return AnalysisSettings{
		ReportType:          DegradationEvent,
		AnalysisKind:        ThresholdAnalysis,
		ThresholdMode:       ThresholdGreaterThan,
		ThresholdValue:      limit,
		ThresholdWindow:     3,
		ThresholdOnCrossing: true,
	}
}

func TestThresholdWindow_OneSlowNightDoesNotCrossTheLimit(t *testing.T) {
	t.Parallel()
	// manyAgents.edtMaxMs on 2026-11-18: 51 ms the night before, then one night of 136 ms.
	assert.Empty(t, detectOverNights([]int{48, 51, 136}, windowedLimit(100)))
}

func TestThresholdWindow_ReportsTheNightOfTheCrossing(t *testing.T) {
	t.Parallel()
	degradations := detectOverNights([]int{45, 48, 44, 130, 128}, windowedLimit(100))
	require.Len(t, degradations, 1)
	assert.Equal(t, "b4", degradations[0].Build)
	assert.InDelta(t, 48, degradations[0].medianValues.previousValue, 0.0001)
	assert.InDelta(t, 128, degradations[0].medianValues.newValue, 0.0001)
}

func TestThresholdWindow_DoesNotReportEachNightOverTheLimit(t *testing.T) {
	t.Parallel()
	// answer.M.streamEndToSettledMs stays at 365-449 ms, over a limit of 300 ms, so each night reported it before.
	nights := []int{396, 449, 420, 365}
	assert.Empty(t, detectOverNights(nights, windowedLimit(300)))

	each := windowedLimit(300)
	each.ThresholdOnCrossing = false
	assert.Len(t, detectOverNights(nights, each), 1, "without the crossing rule, the last night is reported again")
}

func TestThresholdWindow_NeedsAWindowOfValues(t *testing.T) {
	t.Parallel()
	assert.Empty(t, detectOverNights([]int{500, 500}, windowedLimit(100)))
}

func TestThresholdWindow_TheReasonNamesTheLimitAndTheMedian(t *testing.T) {
	t.Parallel()
	degradations := detectOverNights([]int{45, 48, 44, 130, 128}, windowedLimit(100))
	require.Len(t, degradations, 1)
	assert.Equal(t, "Threshold exceeded. The median of the latest 3 values is 128.00, and the limit is 100.00. One run earlier it was 48.00.", getReason(degradations[0]))

	single := detectOverNights([]int{90, 96}, AnalysisSettings{AnalysisKind: ThresholdAnalysis, ThresholdValue: 95})
	require.Len(t, single, 1)
	assert.Equal(t, "Threshold exceeded. The latest value is 96.00, and the limit is 95.00. One run earlier it was 90.00.", getReason(single[0]))
}

func TestTheReasonOfAStepOfTheMedianStaysTheSame(t *testing.T) {
	t.Parallel()
	d := Degradation{medianValues: CenterValues{previousValue: 10, newValue: 20}}
	assert.Equal(t, "Degradation detected. Median changed by: 100.00%. Median was 10.00 and now it is 20.00.", getReason(d))
}
