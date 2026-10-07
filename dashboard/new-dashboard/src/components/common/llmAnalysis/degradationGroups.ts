import { LlmAnalysisDegradationMatch } from "./LlmAnalysisClient"

// one chart (project + metric), possibly degraded at several builds, e.g. on different machines;
// date, rangeSize and matchedCommits are of the strongest build (the smallest range), so they describe one build
interface DegradedChart {
  project: string
  metric: string
  date: string
  rangeSize: number
  matchedCommits: string[]
  // sorted by range size, smallest first
  builds: LlmAnalysisDegradationMatch[]
}

export interface DegradationGroups {
  // set when all charts matched the same commits, then they are shown once instead of on every row
  commonCommits: string[] | null
  groups: { metric: string; charts: DegradedChart[] }[]
  chartCount: number
  // smallest range size over all charts
  strongestRangeSize: number | null
}

export function groupDegradations(degradations: LlmAnalysisDegradationMatch[]): DegradationGroups {
  const charts = new Map<string, DegradedChart>()
  for (const d of degradations.toSorted((a, b) => a.rangeSize - b.rangeSize)) {
    const key = `${d.project}\u0000${d.metric}`
    const chart = charts.get(key)
    if (chart == null) {
      charts.set(key, { project: d.project, metric: d.metric, date: d.date, rangeSize: d.rangeSize, matchedCommits: d.matchedCommits, builds: [d] })
    } else {
      chart.builds.push(d)
    }
  }

  const byMetric = new Map<string, DegradedChart[]>()
  for (const chart of charts.values()) {
    byMetric.set(chart.metric, [...(byMetric.get(chart.metric) ?? []), chart])
  }
  // the metric with more degraded charts first
  const groups = [...byMetric.entries()].map(([metric, metricCharts]) => ({ metric, charts: metricCharts })).toSorted((a, b) => b.charts.length - a.charts.length)

  const all = [...charts.values()]
  const commitsKey = (chart: DegradedChart) => chart.matchedCommits.toSorted().join(",")
  const commonCommits = all.length > 0 && all.every((c) => commitsKey(c) === commitsKey(all[0])) ? all[0].matchedCommits : null
  return {
    commonCommits,
    groups,
    chartCount: all.length,
    strongestRangeSize: all.length === 0 ? null : Math.min(...all.map((c) => c.rangeSize)),
  }
}
