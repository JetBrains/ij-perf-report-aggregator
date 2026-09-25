package degradation_detector

import (
	"testing"
	"testing/synctest"

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
