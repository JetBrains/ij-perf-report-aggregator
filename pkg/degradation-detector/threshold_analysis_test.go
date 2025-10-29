package degradation_detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThresholdGreaterThan_TriggersWhenStrictlyGreater(t *testing.T) {
	t.Parallel()
	values := []int{90, 96}
	builds := []string{"b0", "b1"}
	times := []int64{0, 123}

	settings := AnalysisSettings{
		AnalysisKind:   ThresholdAnalysis,
		ThresholdMode:  ThresholdGreaterThan,
		ThresholdValue: 95,
	}

	degradations := detectDegradations(values, builds, times, settings)
	assert.Len(t, degradations, 1)
	deg := degradations[0]
	assert.Equal(t, "b1", deg.Build)
	assert.Equal(t, int64(123), deg.timestamp)
	assert.True(t, deg.IsDegradation)
	assert.InEpsilon(t, 90.0, deg.medianValues.previousValue, 0.0001)
	assert.InEpsilon(t, 96.0, deg.medianValues.newValue, 0.0001)
}

func TestThresholdGreaterThan_DoesNotTriggerOnEqual(t *testing.T) {
	t.Parallel()
	values := []int{95}
	builds := []string{"b0"}
	times := []int64{0}

	settings := AnalysisSettings{
		AnalysisKind:   ThresholdAnalysis,
		ThresholdMode:  ThresholdGreaterThan,
		ThresholdValue: 95,
	}

	degradations := detectDegradations(values, builds, times, settings)
	assert.Empty(t, degradations)
}

func TestThresholdLessThan_TriggersWhenStrictlyLess(t *testing.T) {
	t.Parallel()
	values := []int{100, 80}
	builds := []string{"b0", "b1"}
	times := []int64{0, 77}

	settings := AnalysisSettings{
		AnalysisKind:   ThresholdAnalysis,
		ThresholdMode:  ThresholdLessThan,
		ThresholdValue: 85,
	}

	degradations := detectDegradations(values, builds, times, settings)
	assert.Len(t, degradations, 1)
	deg := degradations[0]
	assert.Equal(t, "b1", deg.Build)
	assert.Equal(t, int64(77), deg.timestamp)
	assert.True(t, deg.IsDegradation)
	assert.InEpsilon(t, 100.0, deg.medianValues.previousValue, 0.0001)
	assert.InEpsilon(t, 80.0, deg.medianValues.newValue, 0.0001)
}

func TestThresholdLessThan_DoesNotTriggerOnEqual(t *testing.T) {
	t.Parallel()
	values := []int{85}
	builds := []string{"b0"}
	times := []int64{0}

	settings := AnalysisSettings{
		AnalysisKind:   ThresholdAnalysis,
		ThresholdMode:  ThresholdLessThan,
		ThresholdValue: 85,
	}

	degradations := detectDegradations(values, builds, times, settings)
	assert.Empty(t, degradations)
}

func TestThreshold_NoDataReturnsEmpty(t *testing.T) {
	t.Parallel()
	var values []int
	var builds []string
	var times []int64

	settings := AnalysisSettings{
		AnalysisKind:   ThresholdAnalysis,
		ThresholdMode:  ThresholdGreaterThan,
		ThresholdValue: 95,
	}

	degradations := detectDegradations(values, builds, times, settings)
	assert.Empty(t, degradations)
}
