import { DataQueryExecutorConfiguration } from "../common/dataQuery"
import { DataQueryResult, SERIES_NAME_SEPARATOR } from "../common/DataQueryExecutor"

// Series-name parts assembled by DataQueryExecutor follow the order the producers were registered,
// but a producer whose `getSeriesName(i)` is empty (e.g. when its dimension has size 1) contributes
// nothing. That makes positional parsing fragile — set membership against the known dimension
// values is reliable and order-agnostic.
export interface IndexedSeries {
  // values keyed by `${branch}::${project}::${metric}`
  byKey: Map<string, number[]>
  unresolvedNames: string[]
}

interface Dimensions {
  branches: readonly string[]
  projects: readonly string[]
  measures: readonly string[]
  branchesSet: ReadonlySet<string>
  projectsSet: ReadonlySet<string>
  measuresSet: ReadonlySet<string>
}

function makeDimensions(branches: readonly string[], projects: readonly string[], measures: readonly string[]): Dimensions {
  return {
    branches,
    projects,
    measures,
    branchesSet: new Set(branches),
    projectsSet: new Set(projects),
    measuresSet: new Set(measures),
  }
}

export function indexSeries(
  data: DataQueryResult,
  configuration: DataQueryExecutorConfiguration,
  branches: readonly string[],
  projects: readonly string[],
  measures: readonly string[]
): IndexedSeries {
  const byKey = new Map<string, number[]>()
  const unresolvedNames: string[] = []
  const dims = makeDimensions(branches, projects, measures)
  for (let i = 0; i < data.length; i++) {
    const seriesData = data[i]
    if (seriesData == null || seriesData[1] == null) continue
    const values = (seriesData[1] as unknown[]).filter((v): v is number => typeof v === "number" && Number.isFinite(v))
    if (values.length === 0) continue
    const seriesName = configuration.seriesNames[i] ?? ""
    const measureName = configuration.measureNames[i] ?? ""
    const resolved = resolveSeriesKey(seriesName, measureName, dims)
    if (resolved.branch === "" || resolved.project === "" || resolved.metric === "") continue
    if (resolved.fellBack) {
      unresolvedNames.push(seriesName.length > 0 ? seriesName : measureName)
    }
    const key = seriesKey(resolved.branch, resolved.project, resolved.metric)
    const existing = byKey.get(key)
    if (existing == null) {
      byKey.set(key, values)
    } else {
      existing.push(...values)
    }
  }
  return { byKey, unresolvedNames }
}

export function seriesKey(branch: string, project: string, metric: string): string {
  return `${branch}::${project}::${metric}`
}

// One measured build within a series: its generated_time (ms) and value.
export interface SeriesRun {
  t: number
  v: number
}

// Like indexSeries, but retains the per-build (timestamp, value) pairs instead of merging values into
// one bag — the single-run comparison needs to pick a specific run. Each key's runs are sorted by
// timestamp ascending. Series-name resolution reuses resolveSeriesKey, so it matches indexSeries.
export function indexSeriesRuns(
  data: DataQueryResult,
  configuration: DataQueryExecutorConfiguration,
  branches: readonly string[],
  projects: readonly string[],
  measures: readonly string[]
): Map<string, SeriesRun[]> {
  const byKey = new Map<string, SeriesRun[]>()
  const dims = makeDimensions(branches, projects, measures)
  for (let i = 0; i < data.length; i++) {
    const seriesData = data[i]
    // Column 0 is the timestamp array, column 1 the value array (MeasureConfigurator forces this order).
    if (seriesData == null || seriesData[0] == null || seriesData[1] == null) continue
    const timestamps = seriesData[0] as unknown[]
    const values = seriesData[1] as unknown[]
    const seriesName = configuration.seriesNames[i] ?? ""
    const measureName = configuration.measureNames[i] ?? ""
    const resolved = resolveSeriesKey(seriesName, measureName, dims)
    if (resolved.branch === "" || resolved.project === "" || resolved.metric === "") continue

    const runs: SeriesRun[] = []
    for (let k = 0; k < values.length; k++) {
      const t = timestamps[k]
      const v = values[k]
      if (typeof t === "number" && Number.isFinite(t) && typeof v === "number" && Number.isFinite(v)) {
        runs.push({ t, v })
      }
    }
    if (runs.length === 0) continue
    runs.sort((a, b) => a.t - b.t)

    const key = seriesKey(resolved.branch, resolved.project, resolved.metric)
    const existing = byKey.get(key)
    if (existing == null) {
      byKey.set(key, runs)
    } else {
      existing.push(...runs)
      existing.sort((a, b) => a.t - b.t)
    }
  }
  return byKey
}

interface ResolvedSeriesKey {
  branch: string
  project: string
  metric: string
  fellBack: boolean
}

export function resolveSeriesKey(seriesName: string, measureName: string, dims: Dimensions): ResolvedSeriesKey {
  // A dimension whose configurator has size 1 omits its part from the series name, so seed
  // those positions directly. Only the >1-size dimensions need to be matched against the parts.
  const singleBranch = dims.branches.length === 1
  const singleProject = dims.projects.length === 1
  const singleMeasure = dims.measures.length === 1

  let branch = singleBranch ? dims.branches[0] : ""
  let project = singleProject ? dims.projects[0] : ""
  let metric = singleMeasure ? dims.measures[0] : ""

  if (singleBranch && singleProject && singleMeasure) {
    return { branch, project, metric, fellBack: false }
  }

  const parts = (seriesName.length > 0 ? seriesName : measureName).split(SERIES_NAME_SEPARATOR)
  const unclaimed: string[] = []
  for (const part of parts) {
    if (branch === "" && dims.branchesSet.has(part)) branch = part
    else if (project === "" && dims.projectsSet.has(part)) project = part
    else if (metric === "" && dims.measuresSet.has(part)) metric = part
    else unclaimed.push(part)
  }

  // Positional fallback: producers register in the same order as the configurator list
  // (branch, project, measure under the current pipeline), so unclaimed parts line up
  // with whatever dimensions still need a value. This preserves the prior single-pair
  // behavior — if a measure name happens not to match `measuresSet`, the last part of
  // the series name is used.
  let fellBack = false
  let unclaimedIdx = 0
  if (branch === "") {
    branch = unclaimed[unclaimedIdx++] ?? dims.branches[0] ?? ""
    fellBack = true
  }
  if (project === "") {
    project = unclaimed[unclaimedIdx++] ?? dims.projects[0] ?? ""
    fellBack = true
  }
  if (metric === "") {
    metric = unclaimed[unclaimedIdx++] ?? dims.measures.at(-1) ?? ""
    fellBack = true
  }

  return { branch, project, metric, fellBack }
}
