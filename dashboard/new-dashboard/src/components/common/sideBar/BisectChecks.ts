import { InfoData } from "./InfoSidebar"
import { nsToMs } from "../formatter"

// How many points on each side of the selected dot are inspected.
const WINDOW = 10

export interface BisectCheckWarning {
  title: string
  detail: string
}

interface Windows {
  before: number[]
  after: number[]
}

function isDegradation(direction: string): boolean {
  return direction === "DEGRADATION"
}

// windowsAround returns up to WINDOW finite values immediately before and after
// the selected dot (the dot itself is excluded). Returns null when the series
// context is missing or either side has no data to compare.
function windowsAround(data: InfoData): Windows | null {
  const values = data.seriesValues
  const index = data.pointIndex
  if (!values || index == null || index < 0) return null
  const before = values.slice(Math.max(0, index - WINDOW), index).filter((v) => Number.isFinite(v))
  const after = values.slice(index + 1, index + 1 + WINDOW).filter((v) => Number.isFinite(v))
  if (before.length === 0 || after.length === 0) return null
  return { before, after }
}

function toDisplayUnit(value: number, metricType: string | undefined): number {
  return metricType === "ns" ? nsToMs(value) : value
}

function plural(count: number): string {
  return count === 1 ? "" : "s"
}

// checkGraphStability warns when the before/after levels overlap, i.e. the change
// at the selected dot is not cleanly separated from the surrounding noise. For a
// degradation the highest value before the dot must stay below the lowest value
// after it; for an optimization the lowest before must stay above the highest after.
export function checkGraphStability(data: InfoData | null, direction: string): BisectCheckWarning | null {
  if (data == null) return null
  const w = windowsAround(data)
  if (w == null) return null

  const maxBefore = Math.max(...w.before)
  const minBefore = Math.min(...w.before)
  const maxAfter = Math.max(...w.after)
  const minAfter = Math.min(...w.after)

  const scope = `Looking at the ${w.before.length} point${plural(w.before.length)} before and ${w.after.length} point${plural(w.after.length)} after the selected one`

  if (isDegradation(direction)) {
    if (maxBefore < minAfter) return null
    return {
      title: "The metric may be too unstable to bisect",
      detail:
        `${scope}, the highest value before it (${maxBefore}) is not below the lowest value after it (${minAfter}). ` +
        `The values before and after the change overlap, so the regression is not clearly separated and the bisect may not converge. ` +
        `This is a heuristic and can be a false positive.`,
    }
  }

  if (minBefore > maxAfter) return null
  return {
    title: "The metric may be too unstable to bisect",
    detail:
      `${scope}, the lowest value before it (${minBefore}) is not above the highest value after it (${maxAfter}). ` +
      `The values before and after the change overlap, so the optimization is not clearly separated and the bisect may not converge. ` +
      `This is a heuristic and can be a false positive.`,
  }
}

// suggestTargetValue proposes a target value centred in the gap between the
// before and after levels: midway between the highest value before the dot and
// the lowest after it for a degradation, and midway between the lowest before and
// the highest after for an optimization. Returned in display units (milliseconds
// for durations); null when there isn't enough surrounding data.
export function suggestTargetValue(data: InfoData | null, direction: string): number | null {
  if (data == null) return null
  // Don't suggest a value when the before/after levels overlap: the midpoint
  // would fall inside the noise rather than in a real gap between the levels.
  if (checkGraphStability(data, direction) != null) return null
  const w = windowsAround(data)
  if (w == null) return null

  const before = w.before.map((v) => toDisplayUnit(v, data.metricType))
  const after = w.after.map((v) => toDisplayUnit(v, data.metricType))
  const midpoint = isDegradation(direction) ? (Math.max(...before) + Math.min(...after)) / 2 : (Math.min(...before) + Math.max(...after)) / 2
  return Math.round(midpoint)
}

// checkTargetValue warns when the entered target value does not sit between the
// before and after levels. For a degradation, values before the dot should stay
// at or below the target and values after it at or above it; inverse for an
// optimization. targetValue is in display units (milliseconds for durations).
export function checkTargetValue(data: InfoData | null, direction: string, targetValue: number): BisectCheckWarning | null {
  if (data == null || !Number.isFinite(targetValue)) return null
  const w = windowsAround(data)
  if (w == null) return null

  const before = w.before.map((v) => toDisplayUnit(v, data.metricType))
  const after = w.after.map((v) => toDisplayUnit(v, data.metricType))
  const maxBefore = Math.max(...before)
  const minBefore = Math.min(...before)
  const maxAfter = Math.max(...after)
  const minAfter = Math.min(...after)

  const issues: string[] = []
  if (isDegradation(direction)) {
    if (maxBefore > targetValue) {
      issues.push(`the highest value before it (${maxBefore}) is above the target (${targetValue})`)
    }
    if (minAfter < targetValue) {
      issues.push(`the lowest value after it (${minAfter}) is below the target (${targetValue})`)
    }
  } else {
    if (minBefore < targetValue) {
      issues.push(`the lowest value before it (${minBefore}) is below the target (${targetValue})`)
    }
    if (maxAfter > targetValue) {
      issues.push(`the highest value after it (${maxAfter}) is above the target (${targetValue})`)
    }
  }
  if (issues.length === 0) return null

  const expectation = isDegradation(direction) ? "below the target and the values after it above" : "above the target and the values after it below"
  const scope = `the ${w.before.length} point${plural(w.before.length)} before and ${w.after.length} point${plural(w.after.length)} after the selected one`
  return {
    title: "The target value may not separate before from after",
    detail:
      `For this ${direction.toLowerCase()}, across ${scope}, the values before it should be ${expectation} the target, but ${issues.join(" and ")}. ` +
      `This is a heuristic and can be a false positive; double-check the target value.`,
  }
}
