package outlier_detection

import (
	"math"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
)

// RollingMAD calculates rolling Median Absolute Deviation for outlier detection
func rollingMAD(data []int, windowSize int) ([]float64, []float64) {
	medians := make([]float64, len(data))
	mads := make([]float64, len(data))

	for i := range data {
		start := max(0, i-windowSize/2)
		end := min(len(data), i+windowSize/2+1)
		window := data[start:end]

		// Calculate median of window
		med := statistic.Median(window)

		// Calculate MAD (Median Absolute Deviation)
		absDevs := make([]int, len(window))
		for j, val := range window {
			absDevs[j] = int(math.Abs(float64(val) - med))
		}
		mad := statistic.Median(absDevs)

		medians[i] = med
		mads[i] = mad
	}

	return medians, mads
}

// RemoveOutliers removes outliers from data using MAD-based detection
func RemoveOutliers(data []int, windowSize int, threshold float64) []int {
	if len(data) == 0 {
		return data
	}

	medians, mads := rollingMAD(data, windowSize)
	filtered := make([]int, 0, len(data))

	for i, value := range data {
		madScore := math.Abs(float64(value)-medians[i]) / mads[i]
		if madScore <= threshold {
			filtered = append(filtered, value)
		}
	}

	return filtered
}
