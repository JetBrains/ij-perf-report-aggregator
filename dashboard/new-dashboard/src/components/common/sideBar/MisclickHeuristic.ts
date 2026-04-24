import { AccidentKind } from "../../../configurators/accidents/AccidentsConfigurator"
import { typeIsCounter } from "../formatter"
import { InfoData } from "./InfoSidebar"

const NO_CHANGE_THRESHOLD = 0.005
const NEIGHBOR_DOMINANCE_RATIO = 3
const LOOKAHEAD_WINDOW = 5
const LOOKAHEAD_LARGE_CHANGE_THRESHOLD = 0.07

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
  const counterMetric = typeIsCounter(data.metricType ?? "d")

  if (prev == null || !Number.isFinite(prev)) return null

  const relDeltaPrev = relativeChange(value, prev)
  const lookahead = findDominantLookaheadChange(data, value, kind, counterMetric)

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

  if (!directionMatchesKind(value - prev, kind, counterMetric)) {
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

function findDominantLookaheadChange(data: InfoData, value: number, kind: string, counterMetric: boolean): LookaheadHit | null {
  const values = data.seriesValues
  const index = data.pointIndex
  if (!values || index == null || index < 0) {
    // Fall back to the immediate next value (from Delta) when series context is unavailable
    const next = data.nextValue
    if (next == null || !Number.isFinite(next)) return null
    const rel = relativeChange(next, value)
    if (rel >= LOOKAHEAD_LARGE_CHANGE_THRESHOLD && directionMatchesKind(next - value, kind, counterMetric)) {
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
    if (!directionMatchesKind(futureValue - value, kind, counterMetric)) continue
    if (best == null || rel > best.relChange) {
      best = { offset: i - index, relChange: rel }
    }
  }
  return best
}

function relativeChange(a: number, b: number): number {
  const denom = Math.max(Math.abs(a), Math.abs(b))
  if (denom === 0) return 0
  return Math.abs(a - b) / denom
}

function directionMatchesKind(signedDelta: number, kind: string, counterMetric: boolean): boolean {
  if (signedDelta === 0) return false
  const worsened = counterMetric ? signedDelta < 0 : signedDelta > 0
  return kind == AccidentKind.Regression ? worsened : !worsened
}

function percent(ratio: number): string {
  return `${(ratio * 100).toFixed(1)}%`
}

function withArticle(word: string): string {
  return `${/^[aeiou]/i.test(word) ? "an" : "a"} ${word}`
}
