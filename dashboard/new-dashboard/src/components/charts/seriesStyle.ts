import { SERIES_NAME_SEPARATOR } from "../common/DataQueryExecutor"
import { chartTheme } from "./chartTheme"

// The ECharts default palette, which the light theme uses.
const LIGHT_PALETTE = ["#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc"]
const LINE_TYPES = ["solid", "dashed", "dotted"] as const
// The soft cap on the series one chart can tell apart; beyond it, split the chart instead.
const MAX_STYLED_SERIES = 6

export interface SeriesStyle {
  readonly color: string
  readonly lineType: (typeof LINE_TYPES)[number]
}

/**
 * Splits a raw series name (producer titles joined by the separator) into the measure it plots and the rest,
 * which identifies the project. Returns null when the name does not end with one of `measures`.
 */
export function splitSeriesName(rawSeriesName: string, measures: readonly string[]): { group: string; measure: string } | null {
  let best: string | null = null
  for (const measure of measures) {
    const matches = rawSeriesName === measure || rawSeriesName.endsWith(SERIES_NAME_SEPARATOR + measure)
    if (matches && (best == null || measure.length > best.length)) best = measure
  }
  if (best == null) return null
  return { group: rawSeriesName.slice(0, Math.max(0, rawSeriesName.length - best.length - SERIES_NAME_SEPARATOR.length)), measure: best }
}

/**
 * Styles for a chart that plots two or three measures for several projects: one color per project and one line
 * type per measure, so that a project keeps its color across measures. Returns null for any other chart, which
 * keeps the default palette: with one project the measures are told apart by color, and past six series line
 * types stop telling projects apart, so such a chart is better split per project.
 */
export function resolveSeriesStyles(series: readonly ({ group: string; measure: string } | null)[], measures: readonly string[], darkMode: boolean): (SeriesStyle | null)[] | null {
  const groups = [...new Set(series.filter((it) => it != null).map((it) => it.group))]
  const plotted = measures.filter((measure) => series.some((it) => it?.measure === measure))
  if (groups.length < 2 || plotted.length < 2 || plotted.length > LINE_TYPES.length || groups.length * plotted.length > MAX_STYLED_SERIES) return null
  const palette = darkMode ? chartTheme.color : LIGHT_PALETTE
  return series.map((it) =>
    it == null
      ? null
      : {
          color: palette[groups.indexOf(it.group)],
          lineType: LINE_TYPES[plotted.indexOf(it.measure)],
        }
  )
}
