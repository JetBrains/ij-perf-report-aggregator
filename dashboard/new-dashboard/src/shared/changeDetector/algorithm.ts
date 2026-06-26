import { center, disparity, ratio, shift, Sample } from "pragmastat"

export enum ChangePointClassification {
  DEGRADATION = "Degradation",
  OPTIMIZATION = "Optimization",
  NO_CHANGE = "No Change",
}

/**
 * Which direction of change is good for a metric.
 *
 * `"lower"` (the default) treats a rise as a degradation; `"higher"` inverts that; `"stable"` treats
 * any change in either direction as a degradation, for metrics that should not move (for example a
 * completion-result count).
 */
export type BetterDirection = "lower" | "higher" | "stable"

/** Direction a metric value moved at a change point; drives the trend-arrow orientation. */
export type ChangeDirection = "up" | "down"

/** A detected change point: how the value moved, and whether that move is good or bad. */
export interface DetectedChange {
  readonly classification: ChangePointClassification
  readonly direction: ChangeDirection
}

export function detectChanges(seriesData: (string | number)[][], betterDirection: BetterDirection = "lower"): Map<string, DetectedChange> {
  const dataset = seriesData[1] as number[] | undefined
  const changePointIndexes = getChangePointIndexes(dataset, 5)
  const classifications = classifyChangePoint(changePointIndexes, dataset, betterDirection)
  const resultMap = new Map<string, DetectedChange>()

  for (const [index, value] of changePointIndexes.entries()) {
    const extractedValues = extractValuesFromMatrix(seriesData, value)
    resultMap.set(JSON.stringify(extractedValues), classifications[index])
  }
  return resultMap
}

const whichMin = (values: number[]): number => {
  return values.indexOf(Math.min(...values))
}

const getSegmentCost = (partialSums: number[][], tau1: number, tau2: number, k: number, n: number): number => {
  let sum = 0
  for (let i = 0; i < k; i++) {
    const actualSum = partialSums[i][tau2] - partialSums[i][tau1]

    if (actualSum !== 0 && actualSum !== (tau2 - tau1) * 2) {
      const fit = (actualSum * 0.5) / (tau2 - tau1)
      const lnp = (tau2 - tau1) * (fit * Math.log(fit) + (1 - fit) * Math.log1p(-fit))
      sum += lnp
    }
  }
  return ((2 * -Math.log(2 * n - 1)) / k) * sum
}

export const classifyChangePoint = (changePointIndexes: number[], dataset: number[] | undefined, betterDirection: BetterDirection = "lower") => {
  if (dataset == undefined) return []
  const classifications: DetectedChange[] = []

  for (let i = 0; i < changePointIndexes.length; i++) {
    // If it's the first change point, take data from the beginning, otherwise from the previous change point.
    const startBefore = i === 0 ? 0 : changePointIndexes[i - 1]
    const endBefore = changePointIndexes[i]

    const startAfter = changePointIndexes[i]
    // If it's the last change point, take data till the end, otherwise till the next change point.
    const endAfter = i === changePointIndexes.length - 1 ? dataset.length : changePointIndexes[i + 1]

    const segmentBefore = dataset.slice(startBefore, endBefore)
    const segmentAfter = dataset.slice(startAfter, endAfter)

    // ratio() requires strictly positive values; clamp zeros to 1
    const positiveSegmentBefore = segmentBefore.map((v) => Math.max(1, v))
    const positiveSegmentAfter = segmentAfter.map((v) => Math.max(1, v))

    const sampleBefore = Sample.of(segmentBefore)
    const sampleAfter = Sample.of(segmentAfter)
    const positiveSampleBefore = Sample.of(positiveSegmentBefore)
    const positiveSampleAfter = Sample.of(positiveSegmentAfter)

    const centerBefore = center(sampleBefore).value
    const shiftValue = shift(sampleAfter, sampleBefore).value
    const ratioValue = ratio(positiveSampleAfter, positiveSampleBefore).value
    const percentageDifference = Math.abs((ratioValue - 1) * 100)
    const absoluteChange = Math.abs(shiftValue)
    let effectSize: number
    try {
      effectSize = Math.abs(disparity(sampleAfter, sampleBefore).value)
    } catch {
      // disparity throws when a segment has zero spread (all identical values);
      // if there is a shift despite zero spread, it's a clear change
      effectSize = absoluteChange > 0 ? 100 : 0
    }
    let classification

    if (
      (centerBefore < 2000 && percentageDifference < 5) ||
      (centerBefore >= 2000 && centerBefore < 10000 && percentageDifference < 2) ||
      (centerBefore >= 10000 && percentageDifference < 1)
    ) {
      classification = ChangePointClassification.NO_CHANGE
    } else if (absoluteChange < 10) {
      classification = ChangePointClassification.NO_CHANGE
    } else if (effectSize < 2) {
      classification = ChangePointClassification.NO_CHANGE
    } else {
      // "stable" flags any change as a regression; "higher" inverts the default lower-is-better rule.
      const isRegression = betterDirection === "stable" ? true : betterDirection === "higher" ? shiftValue < 0 : shiftValue > 0
      classification = isRegression ? ChangePointClassification.DEGRADATION : ChangePointClassification.OPTIMIZATION
    }

    // Arrow points the way the value actually moved, independent of whether that move is good or bad.
    const direction: ChangeDirection = shiftValue > 0 ? "up" : "down"
    classifications.push({ classification, direction })
  }
  return classifications
}

const getPartialSums = (data: number[], k: number): number[][] => {
  const n = data.length
  const partialSums: number[][] = Array.from({ length: k }, () => Array.from({ length: n + 1 }, () => 0))
  const sortedData = data.toSorted((a, b) => a - b)

  for (let i = 0; i < k; i++) {
    const z = -1 + (2 * i + 1) / k
    const p = 1 / (1 + (2 * n - 1) ** -z)
    const t = sortedData[Math.floor((n - 1) * p)]

    for (let tau = 1; tau <= n; tau++) {
      partialSums[i][tau] = partialSums[i][tau - 1]
      if (data[tau - 1] < t) partialSums[i][tau] += 2
      if (data[tau - 1] === t) partialSums[i][tau] += 1
    }
  }
  return partialSums
}
const extractValuesFromMatrix = (matrix: (string | number)[][], index: number): (string | number)[] => {
  return matrix.map((row) => row[index])
}

export function getChangePointIndexes(data: number[] | undefined, minDistance: number = 1): number[] {
  if (data == undefined) return []
  const n = data.length

  if (n <= 2) return []
  minDistance = Math.min(minDistance, n / 3)

  const penalty = 3 * Math.log(n)
  const k = Math.min(n, Math.ceil(4 * Math.log(n)))

  const partialSums = getPartialSums(data, k)
  const cost = (tau1: number, tau2: number): number => getSegmentCost(partialSums, tau1, tau2, k, n)

  const bestCost: number[] = Array.from({ length: n + 1 }, () => 0)
  bestCost[0] = -penalty
  for (let currentTau = minDistance; currentTau < 2 * minDistance; currentTau++) {
    bestCost[currentTau] = cost(0, currentTau)
  }

  const previousChangePointIndex: number[] = Array.from({ length: n + 1 }, () => 0)
  let previousTaus: number[] = [0, minDistance]

  for (let currentTau = 2 * minDistance; currentTau < n + 1; currentTau++) {
    const costForPreviousTau = previousTaus.map((previousTau) => bestCost[previousTau] + cost(previousTau, currentTau) + penalty)

    const bestPreviousTauIndex = whichMin(costForPreviousTau)
    bestCost[currentTau] = costForPreviousTau[bestPreviousTauIndex]
    previousChangePointIndex[currentTau] = previousTaus[bestPreviousTauIndex]

    const currentBestCost = bestCost[currentTau]
    previousTaus = previousTaus.filter((_, i) => costForPreviousTau[i] < currentBestCost + penalty)
    previousTaus.push(currentTau - (minDistance - 1))
  }

  const changePointIndexes: number[] = []
  let currentIndex = previousChangePointIndex[n]
  while (currentIndex !== 0) {
    changePointIndexes.push(currentIndex)
    currentIndex = previousChangePointIndex[currentIndex]
  }
  return changePointIndexes.toReversed()
}
