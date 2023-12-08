package degradation_detector

import (
  "errors"
  "math"
  "slices"
)

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
  return medianF(differences), nil
}

func hodgesLehmannEstimator(segmentA, segmentB []int) float64 {
  var pairwiseDifferences []int
  for _, valueA := range segmentA {
    for _, valueB := range segmentB {
      pairwiseDifferences = append(pairwiseDifferences, valueB-valueA)
    }
  }
  return medianI(pairwiseDifferences)
}

func medianF(numbers []float64) float64 {
  input := make([]float64, len(numbers))
  copy(input, numbers)
  slices.Sort(input)
  middle := len(input) / 2
  if len(input)%2 == 0 {
    return (input[middle-1] + input[middle]) / 2
  }
  return input[middle]
}

func medianI(numbers []int) float64 {
  input := make([]int, len(numbers))
  copy(input, numbers)
  slices.Sort(input)
  middle := len(input) / 2
  if len(input)%2 == 0 {
    return float64(input[middle-1]+input[middle]) / 2
  }
  return float64(input[middle])
}
