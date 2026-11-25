import { ValueUnit } from "../common/chart"
import { MeasureUnit, resolveMeasureUnit } from "../common/formatter"
import { BetterDirection } from "../../shared/changeDetector/algorithm"
import { getMetricDescription } from "../../shared/metricsDescription"

export interface MetricInfoEntry {
  labels: string[]
  description?: string
  unit: MeasureUnit
  betterDirection: BetterDirection
  url?: string
}

/**
 * The entries of the info popup of a chart: one per distinct description, unit, direction and link of its measures,
 * so that a chart of several measures names each of them. `description` and `betterDirection` override those of
 * every measure.
 */
export function metricInfoEntries(
  measures: readonly string[],
  options: { description?: string; betterDirection?: BetterDirection; valueUnit: ValueUnit; label: (measure: string) => string }
): MetricInfoEntry[] {
  const entries = new Map<string, MetricInfoEntry>()
  for (const measure of measures) {
    const info = getMetricDescription(measure)
    const description = options.description ?? info?.description
    const unit = resolveMeasureUnit(measure, { valueUnit: options.valueUnit, scaling: false })
    const betterDirection = options.betterDirection ?? info?.betterDirection ?? "lower"
    const url = info?.url
    const key = JSON.stringify([description, unit, betterDirection, url])
    const entry = entries.get(key)
    if (entry == undefined) {
      entries.set(key, { labels: [options.label(measure)], description, unit, betterDirection, url })
    } else {
      entry.labels.push(options.label(measure))
    }
  }
  return [...entries.values()]
}
