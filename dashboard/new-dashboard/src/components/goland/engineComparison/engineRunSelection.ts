// Pure run-selection helpers for the engine comparison. Kept out of the composable so the
// aggregated-vs-single reduction is unit-testable without a Vue render context.

import { SeriesRun } from "../../charts/compareQuery"
import { BaseStats, BranchStats, computeBaseStats, computeBranchStats, isSignificantBranch } from "../../charts/compareStats"

// One build measured with both engines.
export interface PairedRun {
  t: number
  legacy: number
  new: number
}

// The builds present in both series, ascending by timestamp. One build publishes both engine variants
// with the same generated_time, so an exact timestamp match pairs them; a build missing one engine is dropped.
export function pairRuns(legacyRuns: readonly SeriesRun[], newRuns: readonly SeriesRun[]): PairedRun[] {
  const newByTime = new Map<number, number>()
  for (const run of newRuns) newByTime.set(run.t, run.v)
  const pairs: PairedRun[] = []
  for (const run of legacyRuns) {
    const value = newByTime.get(run.t)
    if (value != null) pairs.push({ t: run.t, legacy: run.v, new: value })
  }
  return pairs.toSorted((a, b) => a.t - b.t)
}

// The [legacy, new] value arrays for one cell, from paired builds only.
//   - aggregated (single=false): every paired build.
//   - single: the paired build whose timestamp is `run`.
// Returns null when no paired build qualifies — the caller skips that cell.
export function pickRunValues(single: boolean, run: number | null, pairs: readonly PairedRun[]): { legacy: number[]; new: number[] } | null {
  const selected = single ? pairs.filter((pair) => pair.t === run) : pairs
  if (selected.length === 0) return null
  return { legacy: selected.map((pair) => pair.legacy), new: selected.map((pair) => pair.new) }
}

// The LEGACY-vs-NEW statistics of one cell, or null when no paired build qualifies (see pickRunValues).
// The cell is judged by its per-build NEW ÷ LEGACY ratios: both engines run in one build, so the drift
// the build shares (machine, commit) cancels out. The ratio is the center of the selected builds' ratios,
// and its effect size is measured against the spread of the log ratios over all `pairs`, which also
// gives a single build a noise estimate.
export function compareCell(
  single: boolean,
  run: number | null,
  pairs: readonly PairedRun[]
): { legacy: BaseStats; new: BranchStats; ratio: number; diffPercent: number; significant: boolean } | null {
  const values = pickRunValues(single, run, pairs)
  if (values == null) return null
  const legacy = computeBaseStats(values.legacy)
  const newStats = computeBranchStats(legacy, values.new)
  const selected = values.legacy.flatMap((legacyValue, i) => logRatio(legacyValue, values.new[i]))
  if (selected.length === 0) return { legacy, new: newStats, ratio: Number.NaN, diffPercent: Number.NaN, significant: false }
  const center = computeBaseStats(selected).center
  const noise = computeBaseStats(pairs.flatMap((pair) => logRatio(pair.legacy, pair.new))).spread
  const ratio = Math.exp(center)
  const diffPercent = (ratio - 1) * 100
  const disparity = noise > 0 ? center / noise : center === 0 ? 0 : Math.sign(center) * Number.POSITIVE_INFINITY
  const judged = { ...newStats, diffPercent, disparity }
  return { legacy, new: judged, ratio, diffPercent, significant: isSignificantBranch(judged) }
}

function logRatio(legacyValue: number, newValue: number): number[] {
  return legacyValue > 0 && newValue > 0 ? [Math.log(newValue / legacyValue)] : []
}
