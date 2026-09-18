import { median } from "../../../shared/changeDetector/statistic"

export function scaleToMedian(arr: number[], scale = getMedianScale(arr)): number[] {
  return scale === 1 ? arr : arr.map((value) => value * scale)
}

/** Shared by the values and their deviation band. */
export function getMedianScale(arr: number[]): number {
  if (arr.length === 0) {
    return 1
  }
  const medianValue = median(arr)
  return medianValue === 0 ? 1 : 50 / medianValue
}
