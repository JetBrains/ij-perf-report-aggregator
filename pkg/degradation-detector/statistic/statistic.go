package statistic

import (
	"errors"
	"math"
	"slices"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
	hle := hodgesLehmannEstimator(segmentA, segmentB)
	pS, err := pooledShamos(segmentA, segmentB)
	if err != nil {
		return 100 // this will make sure that the change point won't be reported since slices are too small
	}
	return math.Abs(hle / pS)
}

func pooledShamos(x, y []int) (float64, error) {
	n := len(x)
	m := len(y)

	if n < 2 || m < 2 {
		return 0, errors.New("both slices must contain at least two elements")
	}

	shamosX, err := shamosEstimator(x)
	if err != nil {
		return 0, err
	}
	shamosY, err := shamosEstimator(y)
	if err != nil {
		return 0, err
	}

	return math.Sqrt(((float64(n-1) * shamosX * shamosX) + (float64(m-1) * shamosY * shamosY)) / float64(n+m-2)), nil
}

func shamosEstimator(data []int) (float64, error) {
	if len(data) < 2 {
		return 0, errors.New("data slice must contain at least two elements")
	}
	var differences []float64
	for i := range data {
		for j := i + 1; j < len(data); j++ {
			differences = append(differences, math.Abs(float64(data[i]-data[j])))
		}
	}
	return MedianF(differences) * shamosBias(len(data)), nil
}

func hodgesLehmannEstimator(segmentA, segmentB []int) float64 {
	var pairwiseDifferences []int
	for _, valueA := range segmentA {
		for _, valueB := range segmentB {
			pairwiseDifferences = append(pairwiseDifferences, valueB-valueA)
		}
	}
	return Median(pairwiseDifferences)
}

func shamosBias(n int) float64 {
	biasCoefficient := []float64{
		math.NaN(), math.NaN(), 0.1831500, 0.2989400, 0.1582782, 0.1011748, 0.1005038, 0.0676993,
		0.0609574, 0.0543760, 0.0476839, 0.0426722, 0.0385003, 0.0353028, 0.0323526,
		0.0299677, 0.0280421, 0.0262195, 0.0247674, 0.0232297, 0.0220155, 0.0208687,
		0.0199446, 0.0189794, 0.0182343, 0.0174421, 0.0166364, 0.0160158, 0.0153715,
		0.0148940, 0.0144027, 0.0138855, 0.0134510, 0.0130228, 0.0127183, 0.0122444,
		0.0118214, 0.0115469, 0.0113206, 0.0109636, 0.0106308, 0.0104384, 0.0100693,
		0.0098523, 0.0096735, 0.0094973, 0.0092210, 0.0089781, 0.0088083, 0.0086574,
		0.0084772, 0.0082120, 0.0081874, 0.0079775, 0.0078126, 0.0076743, 0.0075212,
		0.0074051, 0.0072528, 0.0071807, 0.0070617, 0.0069123, 0.0067833, 0.0066439,
		0.0065821, 0.0064889, 0.0063844, 0.0062930, 0.0061910, 0.0061255, 0.0060681,
		0.0058994, 0.0058235, 0.0057172, 0.0056805, 0.0056343, 0.0055605, 0.0055011,
		0.0053872, 0.0053062, 0.0052348, 0.0052075, 0.0051173, 0.0050697, 0.0049805,
		0.0048705, 0.0048695, 0.0048287, 0.0047315, 0.0046961, 0.0046698, 0.0046010,
		0.0045544, 0.0045191, 0.0044245, 0.0044074, 0.0043579, 0.0043536, 0.0042874,
		0.0042520, 0.0041864,
	}
	asympt := 0.9538726
	if n <= 100 {
		return 1 / (asympt * (1 + biasCoefficient[n]))
	}
	nf := float64(n)
	return 1 / (asympt * (1 + 0.414253297/nf - 0.442396799/math.Exp2(nf)))
}
