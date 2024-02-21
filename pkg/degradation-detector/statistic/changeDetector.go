package statistic

import (
  "math"
  "slices"
)

func GetChangePointIndexes(data []int, minDistance int) []int {
  n := len(data)
  if n <= 2 {
    return []int{}
  }
  if minDistance < 1 || minDistance > n {
    panic("minDistance should be in range from 1 to len(data)")
  }

  penalty := 3 * math.Log(float64(n))
  k := int(math.Min(float64(n), math.Ceil(4*math.Log(float64(n)))))
  partialSums := getPartialSums(data, k)

  cost := func(tau1, tau2 int) float64 {
    return getSegmentCost(partialSums, tau1, tau2, k, n)
  }

  bestCost := make([]float64, n+1)
  bestCost[0] = -penalty
  for currentTau := minDistance; currentTau < 2*minDistance; currentTau++ {
    bestCost[currentTau] = cost(0, currentTau)
  }

  previousChangePointIndex := make([]int, n+1)
  previousTaus := []int{0, minDistance}
  var costForPreviousTau []float64

  for currentTau := 2 * minDistance; currentTau < n+1; currentTau++ {
    costForPreviousTau = costForPreviousTau[:0]
    for _, previousTau := range previousTaus {
      costForPreviousTau = append(costForPreviousTau, bestCost[previousTau]+cost(previousTau, currentTau)+penalty)
    }
    bestPreviousTauIndex := whichMin(costForPreviousTau)
    bestCost[currentTau] = costForPreviousTau[bestPreviousTauIndex]
    previousChangePointIndex[currentTau] = previousTaus[bestPreviousTauIndex]
    var newPreviousTaus = make([]int, 0, len(previousTaus))
    for i, cost2 := range costForPreviousTau {
      if cost2 < bestCost[currentTau]+penalty {
        newPreviousTaus = append(newPreviousTaus, previousTaus[i])
      }
    }
    previousTaus = newPreviousTaus
    previousTaus = append(previousTaus, currentTau-(minDistance-1))
  }

  var changePointIndexes []int
  currentIndex := previousChangePointIndex[n]
  for currentIndex != 0 {
    changePointIndexes = append(changePointIndexes, currentIndex)
    currentIndex = previousChangePointIndex[currentIndex]
  }
  slices.Reverse(changePointIndexes)
  return changePointIndexes
}

func getPartialSums(data []int, k int) [][]int {
  n := len(data)
  partialSums := make([][]int, k)
  for i := range partialSums {
    partialSums[i] = make([]int, n+1)
  }

  sortedData := make([]int, n)
  copy(sortedData, data)
  slices.Sort(sortedData)

  for i := range k {
    z := -1 + (2*float64(i)+1.0)/float64(k)
    p := 1.0 / (1 + math.Pow(2*float64(n)-1, -z))
    t := sortedData[int(math.Trunc(float64(n-1)*p))]

    for tau := 1; tau <= n; tau++ {
      partialSums[i][tau] = partialSums[i][tau-1]
      if data[tau-1] < t {
        partialSums[i][tau] += 2 // Using doubled value (2) instead of original 1.0
      }
      if data[tau-1] == t {
        partialSums[i][tau]++
      }
    }
  }

  return partialSums
}

func getSegmentCost(partialSums [][]int, tau1, tau2, k, n int) float64 {
  var sum float64
  tauDiff := float64(tau2 - tau1)
  tauDiffDoubled := (tau2 - tau1) * 2
  for i := range k {
    // actualSum is (count(data[j] < t) * 2 + count(data[j] == t) * 1) for j=tau1..tau2-1
    actualSum := partialSums[i][tau2] - partialSums[i][tau1]

    if actualSum != 0 && actualSum != tauDiffDoubled {
      fit := float64(actualSum) / float64(tauDiffDoubled)
      lnp := tauDiff * (fit*math.Log(fit) + (1-fit)*math.Log1p(-fit))
      sum += lnp
    }
  }
  c := -math.Log(float64(2*n - 1))  // Constant from Lemma 3.1 in [Haynes2017]
  return 2.0 * c / float64(k) * sum // See Section 3.1 "Discrete approximation" in [Haynes2017]
}

func whichMin(values []float64) int {
  if len(values) == 0 {
    panic("slice should contain elements")
  }

  minValue := values[0]
  minIndex := 0
  for i := 1; i < len(values); i++ {
    if values[i] < minValue {
      minValue = values[i]
      minIndex = i
    }
  }

  return minIndex
}
