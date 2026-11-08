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
// A single build has no spread of its own, so its disparity is judged against the LEGACY build-to-build
// spread of all `pairs`.
export function compareCell(
  single: boolean,
  run: number | null,
  pairs: readonly PairedRun[]
): { legacy: BaseStats; new: BranchStats; ratio: number; diffPercent: number; significant: boolean } | null {
  const values = pickRunValues(single, run, pairs)
  if (values == null) return null
  const legacy = computeBaseStats(values.legacy)
  const baseline = single ? { ...legacy, spread: computeBaseStats(pairs.map((pair) => pair.legacy)).spread } : legacy
  const newStats = computeBranchStats(baseline, values.new)
  const ratio = Number.isFinite(legacy.center) && legacy.center > 0 && Number.isFinite(newStats.center) ? newStats.center / legacy.center : Number.NaN
  return { legacy, new: newStats, ratio, diffPercent: newStats.diffPercent, significant: isSignificantBranch(newStats) }
}
