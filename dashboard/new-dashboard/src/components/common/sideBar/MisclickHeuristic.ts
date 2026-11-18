import { AccidentKind } from "../../../configurators/accidents/AccidentsConfigurator"
import { InfoData } from "./InfoSidebar"

const NO_CHANGE_THRESHOLD = 0.005
const NEIGHBOR_DOMINANCE_RATIO = 3
const LOOKAHEAD_WINDOW = 5
const LOOKAHEAD_LARGE_CHANGE_THRESHOLD = 0.07
// How many points before the selected one describe the usual noise, and how many
// of them are needed before the noise is judged at all.
const NOISE_WINDOW = 10
const MIN_NOISE_POINTS = 5

export interface MisclickWarning {
  title: string
  detail: string
}

export function detectPossibleMisclick(data: InfoData | null, kind: string): MisclickWarning | null {
  if (data == null) return null
  if (kind != AccidentKind.Regression && kind != AccidentKind.Improvement) return null

  const value = data.series[0]?.rawValue
  if (value == null || !Number.isFinite(value)) return null

  const prev = data.previousValue
  if (prev == null || !Number.isFinite(prev)) return null

  const relDeltaPrev = relativeChange(value, prev)
  const lookahead = findDominantLookaheadChange(data, value, kind)

  if (relDeltaPrev < NO_CHANGE_THRESHOLD) {
    if (lookahead != null) {
      return {
        title: "The selected point doesn't look like the actual event",
        detail:
          `The selected point is within ${percent(relDeltaPrev)} of the previous point, ` +
          `but a ${percent(lookahead.relChange)} change matching this ${kind.toLowerCase()} occurs ${lookahead.offset} point${lookahead.offset == 1 ? "" : "s"} later. ` +
          `Did you mean to click that point?`,
      }
    }
    return {
      title: "The selected point looks unchanged",
      detail: `The selected point is within ${percent(relDeltaPrev)} of the previous point, so there is no meaningful change to attribute to ${withArticle(kind.toLowerCase())}.`,
    }
  }

  if (!directionMatchesKind(value - prev, kind)) {
    const actualKind = kind == AccidentKind.Regression ? "improvement" : "regression"
    const selectedKind = kind.toLowerCase()
    return {
      title: `This change looks like ${withArticle(actualKind)}, not ${withArticle(selectedKind)}`,
      detail:
        `The selected point has a ${value > prev ? "higher" : "lower"} value than the previous point ` +
        `(${percent(relDeltaPrev)} change), which would be ${withArticle(actualKind)}. ` +
        `Please verify the correct point and event type are selected.`,
    }
  }

  // A change above the fixed threshold can still be ordinary noise: on a metric that
  // swings by tens of percent, a point inside the range of the points before it
  // doesn't stand out, whatever its delta to the immediate predecessor.
  const noise = findNoiseRangeBefore(data)
  if (noise != null && value >= noise.min && value <= noise.max) {
    const noiseDetail =
      `The selected point changed by ${percent(relDeltaPrev)} from the previous one, ` +
      `but the ${noise.count} points before it already vary by ${percent(relativeChange(noise.max, noise.min))} and its value lies within their range`
    if (lookahead != null) {
      return {
        title: "The selected point doesn't look like the actual event",
        detail:
          `${noiseDetail}. ` +
          `A ${percent(lookahead.relChange)} change matching this ${kind.toLowerCase()} occurs ${lookahead.offset} point${lookahead.offset == 1 ? "" : "s"} later. ` +
          `Did you mean to click that point?`,
      }
    }
    return {
      title: "The selected point doesn't stand out from the noise",
      detail: `${noiseDetail}, so there is no clear change to attribute to ${withArticle(kind.toLowerCase())}.`,
    }
  }

  if (lookahead != null && lookahead.relChange >= relDeltaPrev * NEIGHBOR_DOMINANCE_RATIO) {
    return {
      title: "A later point has a much larger change",
      detail:
        `A ${percent(lookahead.relChange)} change matching this ${kind.toLowerCase()} occurs ${lookahead.offset} point${lookahead.offset == 1 ? "" : "s"} later, ` +
        `which is much larger than the ${percent(relDeltaPrev)} change at the selected point. ` +
        `Did you mean to click that point?`,
    }
  }

  return null
}

interface LookaheadHit {
  offset: number
  relChange: number
}

function findDominantLookaheadChange(data: InfoData, value: number, kind: string): LookaheadHit | null {
  const values = data.seriesValues
  const index = data.pointIndex
  if (!values || index == null || index < 0) {
    // Fall back to the immediate next value (from Delta) when series context is unavailable
    const next = data.nextValue
    if (next == null || !Number.isFinite(next)) return null
    const rel = relativeChange(next, value)
    if (rel >= LOOKAHEAD_LARGE_CHANGE_THRESHOLD && directionMatchesKind(next - value, kind)) {
      return { offset: 1, relChange: rel }
    }
    return null
  }

  let best: LookaheadHit | null = null
  const end = Math.min(values.length - 1, index + LOOKAHEAD_WINDOW)
  for (let i = index + 1; i <= end; i++) {
    const futureValue = values[i]
    if (!Number.isFinite(futureValue)) continue
    const rel = relativeChange(futureValue, value)
    if (rel < LOOKAHEAD_LARGE_CHANGE_THRESHOLD) continue
    if (!directionMatchesKind(futureValue - value, kind)) continue
    if (best == null || rel > best.relChange) {
      best = { offset: i - index, relChange: rel }
    }
  }
  return best
}

interface NoiseRange {
  min: number
  max: number
  count: number
}

// findNoiseRangeBefore returns the range of up to NOISE_WINDOW finite values right
// before the selected point, or null when the series context is missing or too
// short to tell noise from a change.
function findNoiseRangeBefore(data: InfoData): NoiseRange | null {
  const values = data.seriesValues
  const index = data.pointIndex
  if (!values || index == null || index < 0) return null
  const before = values.slice(Math.max(0, index - NOISE_WINDOW), index).filter((v) => Number.isFinite(v))
  if (before.length < MIN_NOISE_POINTS) return null
  return { min: Math.min(...before), max: Math.max(...before), count: before.length }
}

function relativeChange(a: number, b: number): number {
  const denom = Math.max(Math.abs(a), Math.abs(b))
  if (denom === 0) return 0
  return Math.abs(a - b) / denom
}

function directionMatchesKind(signedDelta: number, kind: string): boolean {
  if (signedDelta === 0) return false
  const worsened = signedDelta > 0
  return kind == AccidentKind.Regression ? worsened : !worsened
}

function percent(ratio: number): string {
  return `${(ratio * 100).toFixed(1)}%`
}

function withArticle(word: string): string {
  return `${/^[aeiou]/i.test(word) ? "an" : "a"} ${word}`
}
