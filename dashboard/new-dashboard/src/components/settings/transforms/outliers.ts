import { rollingMad } from "../../../shared/changeDetector/statistic"

export function removeOutliers(data: (string | number)[][], windowSize: number = 5, threshold: number = 3): (string | number)[][] {
  if (data.length === 0) {
    return data
  }

  const values = data[1] as number[]

  const { medians, mads } = rollingMad(values, windowSize)
  const skippedIndex: number[] = []

  for (const [i, value] of values.entries()) {
    const madScore = Math.abs(value - medians[i]) / mads[i]
    if (madScore > threshold) {
      skippedIndex.push(i)
    }
  }

  return skippedIndex.length === 0 ? data : data.map((row) => row.filter((_, index) => !skippedIndex.includes(index)))
}
