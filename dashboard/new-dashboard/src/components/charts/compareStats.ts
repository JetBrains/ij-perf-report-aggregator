// Robust per-row statistics for the "Compare with base" table.
// Location, shift, and spread all come from pragmastat's pairwise (Hodges-Lehmann / Shamos)
// family — internally consistent and more efficient than positional median + MAD on clean data.

import { Sample, center, shift, spread } from "pragmastat"

function sortedCopy(values: number[]): number[] {
  return values.toSorted((a, b) => a - b)
}

export function percentile(sorted: number[], p: number): number {
  const n = sorted.length
  if (n === 0) return Number.NaN
  if (n === 1) return sorted[0]
  const rank = (p / 100) * (n - 1)
  const lo = Math.floor(rank)
  const hi = Math.ceil(rank)
  if (lo === hi) return sorted[lo]
  return sorted[lo] + (rank - lo) * (sorted[hi] - sorted[lo])
}

function safeCenter(values: number[]): number {
  if (values.length === 0) return Number.NaN
  if (values.length === 1) return values[0]
  try {
    return center(Sample.of(values)).value
  } catch {
    return Number.NaN
  }
}

function safeShift(branchValues: number[], baseValues: number[]): number {
  if (branchValues.length === 0 || baseValues.length === 0) return Number.NaN
  try {
    return shift(Sample.of(branchValues), Sample.of(baseValues)).value
  } catch {
    return Number.NaN
  }
}

function safeSpread(values: number[]): number {
  // Pragmastat throws AssumptionError on tie-dominant samples (> 50 % identical). We catch and
  // return 0; computeBaseStats then falls back to a p10–p90-based estimate so the disparity
  // column doesn't go to ±∞ on quantized perf data (e.g. integer counters where most builds
  // report the exact same number).
  if (values.length < 2) return 0
  try {
    return spread(Sample.of(values)).value
  } catch {
    return 0
  }
}

// p90 − p10 ≈ 2.5631 · σ for normal data, and Shamos ≈ 0.954 · σ — so (p90 − p10) / 2.687
// is the p10–p90-range estimate on the Shamos scale. Used only when Shamos itself collapses.
const P10_P90_TO_SHAMOS = 2.687

export interface BaseStats {
  count: number
  center: number // Hodges-Lehmann pseudomedian
  p10: number
  p90: number
  spread: number // Shamos spread; 0 when the sample is tie-dominant
  // Kept on the struct so computeBranchStats can compute pragmastat's `shift` (median of pairwise
  // branch[i] − base[j]) without the caller threading raw values through a second argument.
  values: readonly number[]
}

export interface BranchStats {
  count: number
  center: number // Hodges-Lehmann pseudomedian
  diffPercent: number // (branch.center − base.center) / base.center * 100
  disparity: number // shift(branch, base) / base.spread; ±Infinity if spread=0 and centers differ
}

export function computeBaseStats(values: number[]): BaseStats {
  const sorted = sortedCopy(values)
  const p10 = percentile(sorted, 10)
  const p90 = percentile(sorted, 90)
  let dispersion = safeSpread(values)
  if (dispersion === 0) {
    const range = p90 - p10
    if (range > 0) dispersion = range / P10_P90_TO_SHAMOS
  }
  return {
    count: values.length,
    center: safeCenter(values),
    p10,
    p90,
    spread: dispersion,
    values,
  }
}

export function computeBranchStats(base: BaseStats, values: number[]): BranchStats {
  const branchCenter = safeCenter(values)

  let diffPercent = Number.NaN
  if (Number.isFinite(base.center) && base.center !== 0 && Number.isFinite(branchCenter)) {
    diffPercent = ((branchCenter - base.center) / base.center) * 100
  }

  let disparity = Number.NaN
  if (Number.isFinite(branchCenter) && Number.isFinite(base.center)) {
    if (base.spread > 0) {
      disparity = safeShift(values, base.values as number[]) / base.spread
    } else if (branchCenter === base.center) {
      disparity = 0
    } else {
      disparity = branchCenter > base.center ? Number.POSITIVE_INFINITY : Number.NEGATIVE_INFINITY
    }
  }

  return { count: values.length, center: branchCenter, diffPercent, disparity }
}

export const DISPARITY_SIGNIFICANT_THRESHOLD = 2
