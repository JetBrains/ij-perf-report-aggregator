export function median(arr: number[]): number {
  if (arr.length === 0) {
    throw new Error("Data array is empty")
  }
  const sortedArr = [...arr].sort((a, b) => a - b)
  const mid = Math.floor(sortedArr.length / 2)
  return sortedArr.length % 2 === 0 ? (sortedArr[mid - 1] + sortedArr[mid]) / 2 : sortedArr[mid]
}

export function calculatePercentile(arr: number[], percentile: number): number {
  if (arr.length === 0) {
    return 0
  }
  if (percentile < 0 || percentile > 100) {
    throw new Error("Percentile must be between 0 and 100")
  }

  const sortedArr = [...arr].sort((a, b) => a - b)
  const index = (percentile / 100) * (sortedArr.length - 1)
  const lowerIndex = Math.floor(index)
  const upperIndex = Math.ceil(index)

  // If the index is an integer, return the value directly
  if (lowerIndex === upperIndex) {
    return sortedArr[lowerIndex]
  }

  // Otherwise, interpolate between the two closest values
  const lowerValue = sortedArr[lowerIndex]
  const upperValue = sortedArr[upperIndex]
  const weight = index - lowerIndex

  return lowerValue + weight * (upperValue - lowerValue)
}

export function hodgesLehmannEstimator(segmentA: number[], segmentB: number[]): number {
  const pairwiseDifferences: number[] = []
  for (const valueA of segmentA) {
    for (const valueB of segmentB) {
      pairwiseDifferences.push(valueB - valueA)
    }
  }
  return median(pairwiseDifferences)
}

export function pooledShamos(x: number[], y: number[]): number {
  const n = x.length
  const m = y.length

  if (n < 2 || m < 2) {
    throw new Error("Both arrays must contain at least two elements")
  }

  const shamosX = shamosEstimator(x)
  const shamosY = shamosEstimator(y)

  return Math.sqrt(((n - 1) * shamosX * shamosX + (m - 1) * shamosY * shamosY) / (n + m - 2))
}

export function shamosEstimator(data: number[]): number {
  if (data.length < 2) {
    throw new Error("Data array must contain at least two elements")
  }
  const differences: number[] = []
  for (let i = 0; i < data.length; i++) {
    for (let j = i + 1; j < data.length; j++) {
      differences.push(Math.abs(data[i] - data[j]))
    }
  }
  return median(differences) * shamosBias(data.length)
}

export function shamosBias(n: number): number {
  const biasCoefficient = [
    Number.NaN,
    Number.NaN,
    0.18315,
    0.29894,
    0.1582782,
    0.1011748,
    0.1005038,
    0.0676993,
    0.0609574,
    0.054376,
    0.0476839,
    0.0426722,
    0.0385003,
    0.0353028,
    0.0323526,
    0.0299677,
    0.0280421,
    0.0262195,
    0.0247674,
    0.0232297,
    0.0220155,
    0.0208687,
    0.0199446,
    0.0189794,
    0.0182343,
    0.0174421,
    0.0166364,
    0.0160158,
    0.0153715,
    0.014894,
    0.0144027,
    0.0138855,
    0.013451,
    0.0130228,
    0.0127183,
    0.0122444,
    0.0118214,
    0.0115469,
    0.0113206,
    0.0109636,
    0.0106308,
    0.0104384,
    0.0100693,
    0.0098523,
    0.0096735,
    0.0094973,
    0.009221,
    0.0089781,
    0.0088083,
    0.0086574,
    0.0084772,
    0.008212,
    0.0081874,
    0.0079775,
    0.0078126,
    0.0076743,
    0.0075212,
    0.0074051,
    0.0072528,
    0.0071807,
    0.0070617,
    0.0069123,
    0.0067833,
    0.0066439,
    0.0065821,
    0.0064889,
    0.0063844,
    0.006293,
    0.006191,
    0.0061255,
    0.0060681,
    0.0058994,
    0.0058235,
    0.0057172,
    0.0056805,
    0.0056343,
    0.0055605,
    0.0055011,
    0.0053872,
    0.0053062,
    0.0052348,
    0.0052075,
    0.0051173,
    0.0050697,
    0.0049805,
    0.0048705,
    0.0048695,
    0.0048287,
    0.0047315,
    0.0046961,
    0.0046698,
    0.004601,
    0.0045544,
    0.0045191,
    0.0044245,
    0.0044074,
    0.0043579,
    0.0043536,
    0.0042874,
    0.004252,
    0.0041864,
  ]
  const asympt = 0.9538726
  return n <= 100 ? 1 / (asympt * (1 + biasCoefficient[n])) : 1 / (asympt * (1 + 0.414253297 / n - 0.442396799 / Math.pow(n, 2)))
}
