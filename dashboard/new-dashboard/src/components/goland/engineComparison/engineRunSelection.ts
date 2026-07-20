// Pure run-selection helpers for the engine comparison. Kept out of the composable so the
// aggregated-vs-single reduction is unit-testable without a Vue render context.

import { SeriesRun } from "../../charts/compareQuery"

export const MS_PER_DAY = 86_400_000

// The UTC-day index of a timestamp (ms). Runs in the same nightly share a day, so this aligns the two
// engine variants (which are measured in the same build) onto one selectable snapshot.
export function dayKey(timestampMs: number): number {
  return Math.floor(timestampMs / MS_PER_DAY)
}

// The latest run of each day, keyed by day. Runs are assumed sorted ascending by timestamp, so a later
// same-day run overwrites an earlier one.
export function latestRunByDay(runs: readonly SeriesRun[]): Map<number, SeriesRun> {
  const byDay = new Map<number, SeriesRun>()
  for (const run of runs) byDay.set(dayKey(run.t), run)
  return byDay
}

// The [legacy, new] value arrays for one cell.
//   - aggregated (single=false): every run's value.
//   - single: the run of `selectedDay` (or the latest day both engines share when selectedDay is null).
// Returns null when a single-run day lacks one of the engines — the caller skips that cell.
export function pickRunValues(single: boolean, selectedDay: number | null, legacyRuns: readonly SeriesRun[], newRuns: readonly SeriesRun[]): { legacy: number[]; new: number[] } | null {
  if (!single) {
    return { legacy: legacyRuns.map((run) => run.v), new: newRuns.map((run) => run.v) }
  }
  const legacyByDay = latestRunByDay(legacyRuns)
  const newByDay = latestRunByDay(newRuns)
  let day = selectedDay
  if (day == null) {
    const commonDays = [...legacyByDay.keys()].filter((d) => newByDay.has(d))
    if (commonDays.length === 0) return null
    day = Math.max(...commonDays)
  }
  const legacyRun = legacyByDay.get(day)
  const newRun = newByDay.get(day)
  if (legacyRun == null || newRun == null) return null
  return { legacy: [legacyRun.v], new: [newRun.v] }
}
