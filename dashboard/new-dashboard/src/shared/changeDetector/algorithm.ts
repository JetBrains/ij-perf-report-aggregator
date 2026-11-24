import { center, compare2, ComparisonVerdict, disparity, Measurement, Metric, ratio, shift, Sample, Threshold } from "pragmastat"

export enum ChangePointClassification {
  DEGRADATION = "Degradation",
  OPTIMIZATION = "Optimization",
  NO_CHANGE = "No Change",
  /** A real change on a metric with no good/bad direction: marked, but not colored as good or bad. */
  NEUTRAL = "Neutral",
}

/**
 * Which direction of change is good for a metric.
 *
 * `"lower"` (the default) treats a rise as a degradation; `"higher"` inverts that; `"stable"` treats
 * any change in either direction as a degradation, for metrics that should not move (for example the
 * file count of a fixed project). `"none"` still detects changes and marks them, but never colors them
 * good or bad — for metrics with no meaningful direction (for example a count that legitimately varies
 * run-to-run).
 */
export type BetterDirection = "lower" | "higher" | "stable" | "none"

/** Direction a metric value moved at a change point; drives the trend-arrow orientation. */
type ChangeDirection = "up" | "down"

/** A detected change point: how the value moved, and whether that move is good or bad. */
export interface DetectedChange {
  readonly classification: ChangePointClassification
  readonly direction: ChangeDirection
}

/**
 * When a change point counts as a change.
 *
 * `"default"` keeps the thresholds tuned for IDE metrics: a relative change of at least 5 % (2 % from 2000, 1 % from
 * 10000), an absolute shift of at least 10 and a disparity of at least 2. They suit integer milliseconds and counters,
 * but hide every change of small values such as frame times or frame rates.
 *
 * `"bounds"` is unit-free: the bounds of the Hodges–Lehmann ratio of the segments (after / before) at `misrate` must lie
 * wholly beyond a change of `minRelativeChange` in either direction, so that the change is confidently at least that
 * large (pragmastat's `compare2`). The values must be positive; others are left out.
 */
export type ChangeCriteria = { readonly kind: "default" } | { readonly kind: "bounds"; readonly misrate: number; readonly minRelativeChange: number }

export const DEFAULT_CHANGE_CRITERIA: ChangeCriteria = { kind: "default" }

export function detectChanges(
  seriesData: (string | number)[][],
  betterDirection: BetterDirection = "lower",
  criteria: ChangeCriteria = DEFAULT_CHANGE_CRITERIA
): Map<string, DetectedChange> {
  const dataset = seriesData[1] as number[] | undefined
  const changePointIndexes = getChangePointIndexes(dataset, 5)
  const classifications = classifyChangePoint(changePointIndexes, dataset, betterDirection, criteria)
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

/**
 * How a sample changed against an earlier one, both given as logarithms: the Hodges–Lehmann estimate of the ratio
 * after / before as a relative change (0.05 is 5 % more), with its bounds at the misrate.
 */
export interface RelativeChange {
  change: number
  /** The bounds of the change; null when the samples are too small for the misrate. */
  lower: number | null
  upper: number | null
  /** The bounds lie wholly beyond the minimum change in one direction: the change is confidently at least that large. */
  significant: boolean
}

/**
 * Compares two samples of logarithms (of positive values) with pragmastat's `compare2`: the shift of the logarithms is
 * the logarithm of the Hodges–Lehmann ratio (exactly pragmastat's `ratio`), and its bounds are tested against the
 * minimum change up and down, ±log(1 + minRelativeChange), so that a change and its inverse count alike.
 */
export function relativeChangeOfLogs(afterLogs: number[], beforeLogs: number[], misrate: number, minRelativeChange: number): RelativeChange {
  if (afterLogs.length === 0 || beforeLogs.length === 0) return { change: 0, lower: null, upper: null, significant: false }
  const minimum = Math.log1p(minRelativeChange)
  const after = Sample.of(afterLogs)
  const before = Sample.of(beforeLogs)
  try {
    const [up, down] = compare2(after, before, [new Threshold(Metric.Shift, new Measurement(minimum), misrate), new Threshold(Metric.Shift, new Measurement(-minimum), misrate)])
    return {
      change: Math.expm1(up.estimate.value),
      lower: Math.expm1(up.bounds.lower),
      upper: Math.expm1(up.bounds.upper),
      significant: up.verdict === ComparisonVerdict.Greater || down.verdict === ComparisonVerdict.Less,
    }
  } catch {
    // The samples are too small for the misrate: no pair of order statistics bounds the shift that tightly.
    return { change: Math.expm1(shift(after, before).value), lower: null, upper: null, significant: false }
  }
}

export const classifyChangePoint = (
  changePointIndexes: number[],
  dataset: number[] | undefined,
  betterDirection: BetterDirection = "lower",
  criteria: ChangeCriteria = DEFAULT_CHANGE_CRITERIA
) => {
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

    if (criteria.kind === "bounds") {
      const { change, significant } = relativeChangeOfLogs(positiveLogs(segmentAfter), positiveLogs(segmentBefore), criteria.misrate, criteria.minRelativeChange)
      const direction: ChangeDirection = change > 0 ? "up" : "down"
      classifications.push({ classification: significant ? judge(betterDirection, change) : ChangePointClassification.NO_CHANGE, direction })
      continue
    }

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
      classification = judge(betterDirection, shiftValue)
    }

    // Arrow points the way the value actually moved, independent of whether that move is good or bad.
    const direction: ChangeDirection = shiftValue > 0 ? "up" : "down"
    classifications.push({ classification, direction })
  }
  return classifications
}

/** The logarithms of the positive values; the others have none. */
function positiveLogs(values: number[]): number[] {
  return values.filter((value) => value > 0).map((value) => Math.log(value))
}

/** Whether a real change of the given sign (positive is a rise) is good or bad for the metric. */
function judge(betterDirection: BetterDirection, change: number): ChangePointClassification {
  // A real change, but the metric has no good/bad direction: mark it without judging it.
  if (betterDirection === "none") return ChangePointClassification.NEUTRAL
  // "stable" flags any change as a regression; "higher" inverts the default lower-is-better rule.
  const isRegression = betterDirection === "stable" ? true : betterDirection === "higher" ? change < 0 : change > 0
  return isRegression ? ChangePointClassification.DEGRADATION : ChangePointClassification.OPTIMIZATION
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
