package statistic

import (
	"math"
	"slices"

	"github.com/AndreyAkinshin/pragmastat/go/v3"
)

func MedianF(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	var sortedNums []float64
	if slices.IsSorted(nums) {
		sortedNums = nums
	} else {
		sortedNums = make([]float64, len(nums))
		copy(sortedNums, nums)
		slices.Sort(sortedNums)
	}

	middle := len(sortedNums) / 2
	if len(sortedNums)%2 == 0 {
		return (sortedNums[middle-1] + sortedNums[middle]) / 2
	}

	return sortedNums[middle]
}

func Median(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}

	var sortedNums []int
	if slices.IsSorted(nums) {
		sortedNums = nums
	} else {
		sortedNums = make([]int, len(nums))
		copy(sortedNums, nums)
		slices.Sort(sortedNums)
	}

	middle := len(sortedNums) / 2
	if len(sortedNums)%2 == 0 {
		return float64(sortedNums[middle-1]+sortedNums[middle]) / 2
	}

	return float64(sortedNums[middle])
}

func EffectSize(segmentA, segmentB []int) float64 {
	shift, err := pragmastat.Shift(segmentB, segmentA)
	if err != nil {
		return 100 // this will make sure that the change point won't be reported since slices are too small
	}

	avgSpread, err := pragmastat.AvgSpread(segmentA, segmentB)
	if err != nil {
		return 100
	}

	if avgSpread == 0 {
		return 100
	}

	return math.Abs(shift / avgSpread)
}
