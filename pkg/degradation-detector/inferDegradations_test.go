package degradation_detector

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInferDegradations(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		degraded := PerformanceSettings{Metric: "degraded", MinimumSegmentLength: 3}
		stable := PerformanceSettings{Metric: "stable", MinimumSegmentLength: 3}
		data := make(chan QueryResultWithSettings, 2)
		data <- QueryResultWithSettings{
			queryResult: queryResultOf([]int{101, 99, 101, 99, 101, 201, 202, 201, 202, 201, 201, 201}),
			Settings:    degraded,
		}
		data <- QueryResultWithSettings{
			queryResult: queryResultOf([]int{101, 99, 101, 99, 101, 99, 101, 99, 101, 99, 101, 99}),
			Settings:    stable,
		}
		close(data)

		// ranging until close: if the output channel is never closed, synctest reports a deadlock instead of hanging
		var degradations []DegradationWithSettings
		for degradation := range InferDegradations(data) {
			degradations = append(degradations, degradation)
		}

		require.Len(t, degradations, 1)
		assert.Equal(t, degraded, degradations[0].Settings)
		assert.Equal(t, int64(5), degradations[0].Details.timestamp)
	})
}

func TestInferDegradationsIgnoresChangesBeforeTheSetTime(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		values := []int{101, 99, 101, 99, 101, 201, 202, 201, 202, 201, 201, 201}
		// The change starts at timestamp 5.
		for _, tc := range []struct {
			ignoreBefore int64
			wantChanges  int
		}{
			{ignoreBefore: 5, wantChanges: 1},
			{ignoreBefore: 6, wantChanges: 0},
		} {
			settings := PerformanceSettings{Metric: "degraded", MinimumSegmentLength: 3, IgnoreChangesBefore: time.UnixMilli(tc.ignoreBefore)}
			data := make(chan QueryResultWithSettings, 1)
			data <- QueryResultWithSettings{queryResult: queryResultOf(values), Settings: settings}
			close(data)

			var degradations []DegradationWithSettings
			for degradation := range InferDegradations(data) {
				degradations = append(degradations, degradation)
			}
			assert.Len(t, degradations, tc.wantChanges, "ignore before %d", tc.ignoreBefore)
		}
	})
}

func queryResultOf(values []int) queryResult {
	timestamps := make([]int64, len(values))
	for i := range timestamps {
		timestamps[i] = int64(i)
	}
	return queryResult{
		values:     values,
		builds:     make([]string, len(values)),
		timestamps: timestamps,
	}
}
