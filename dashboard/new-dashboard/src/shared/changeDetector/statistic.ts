export function median(arr: number[]): number {
  if (arr.length === 0) {
    throw new Error("Data array is empty")
  }
  const sortedArr = [...arr].sort((a, b) => a - b)
  const mid = Math.floor(sortedArr.length / 2)
  return sortedArr.length % 2 === 0 ? (sortedArr[mid - 1] + sortedArr[mid]) / 2 : sortedArr[mid]
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
  // Generate all unique pairs and calculate their absolute differences
  for (let i = 0; i < data.length; i++) {
    for (let j = i + 1; j < data.length; j++) {
      differences.push(Math.abs(data[i] - data[j]))
    }
  }
  return median(differences)
}
