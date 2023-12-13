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
  for i := 0; i < len(data); i++ {
    for j := i + 1; j < len(data); j++ {
      differences = append(differences, math.Abs(float64(data[i]-data[j])))
    }
  }
  return MedianF(differences), nil
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
