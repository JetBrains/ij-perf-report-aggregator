import { computed, onUnmounted, Ref, shallowRef, triggerRef, watch } from "vue"
import { DimensionConfigurator, selectedToArray } from "../../../configurators/DimensionConfigurator"
import { PredefinedMeasureConfigurator } from "../../../configurators/MeasureConfigurator"
import { injectOrError } from "../../../shared/injectionKeys"
import { branchConfiguratorKey, dashboardConfiguratorsKey, serverConfiguratorKey } from "../../../shared/keys"
import { DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../common/dataQuery"
import { DataQueryExecutor, DataQueryResult } from "../../common/DataQueryExecutor"
import { formatMeasureValue, resolveMeasureUnit } from "../../common/formatter"
import { indexSeriesRuns, seriesKey, SeriesRun } from "../../charts/compareQuery"
import { BaseStats, BranchStats, computeBaseStats, computeBranchStats, isSignificantBranch } from "../../charts/compareStats"
import { bucketOf, buildMeasure, HighlightingProject, phaseOf, projects, Quantity } from "./highlightingMetrics"
import { computeEngineAggregates, EngineAggregates, EngineComparePoint } from "./engineCompareAggregates"
import { dayKey, pickRunValues } from "./engineRunSelection"

// How the comparison reduces the runs in the window: "aggregated" = robust center over every run
// (SPEC-style verdict), "single" = one run's before/after (see EngineComparisonMode below).
export type EngineComparisonMode = "aggregated" | "single"

// A selectable run in single-run mode. `day` is the UTC-day key used to align the two engines; `label`
// is the human date. LEGACY and NEW of a project share a build, so a day picks a coherent snapshot.
export interface RunDay {
  day: number
  label: string
}

// One comparison cell: LEGACY vs NEW for one project base and one metric type, carrying both the raw
// robust stats and the derived ratio/diff the views render. Assignable to EngineComparePoint so the
// pure aggregator can consume rows directly.
export interface EngineCompareRow extends EngineComparePoint {
  // Concrete measure name (with the quantity suffix), for the drill-down link.
  measure: string
  legacy: BaseStats
  new: BranchStats
}

// The LEGACY and NEW centers of a row, formatted with the row's resolved measure unit. Shared by the
// heatmap and the ranked-bars tooltip so their "LEGACY → NEW" text stays identical.
export function formatEngineCell(row: EngineCompareRow): { before: string; after: string } {
  const unit = resolveMeasureUnit(row.measure)
  return { before: formatMeasureValue(row.legacy.center, unit), after: formatMeasureValue(row.new.center, unit) }
}

interface EngineComparisonOptions {
  // The metric types the user selected in the toolbar (subset of METRIC_TYPES, in that order). Read-only
  // so a getter-backed ref (toRef(() => prop)) is accepted alongside a plain writable ref.
  selectedMetricTypes: Readonly<Ref<string[]>>
  // The quantity to plot (duration or a JVM sub-metric).
  quantity: Readonly<Ref<Quantity>>
  // Aggregated vs single-run view.
  mode: Readonly<Ref<EngineComparisonMode>>
  // In single-run mode, the run day to show; null means the latest run present.
  selectedRunDay: Readonly<Ref<number | null>>
}

interface EngineComparison {
  rows: Ref<EngineCompareRow[]>
  aggregates: Ref<EngineAggregates>
  runDays: Ref<RunDay[]>
  loading: Ref<boolean>
}

// One live query per project base, mirroring today's per-project GroupProjectsChart ownership. Measure
// names embed the project's scenario file, so the bases cannot share a single executor. The last result
// is cached so a mode/run switch recomputes without a refetch (the query is unchanged).
interface BaseEntry {
  project: HighlightingProject
  legacyProject: string
  newProject: string
  measures: Ref<string[]>
  lastData: DataQueryResult | null
  lastConfiguration: DataQueryExecutorConfiguration | null
  unsubscribe: () => void
}

// Drives the LEGACY-vs-NEW comparison: builds its own executors (independent of the lazy drill-down
// accordion's mount state), indexes each result by (branch, project variant, measure), and derives one
// row per (base, metric type). The engine variants are fixed (both are always queried); the Engine
// toolbar selector gates only the drill-down charts and the comparison's own render, not this query.
export function useEngineComparison(opts: EngineComparisonOptions): EngineComparison {
  const serverConfigurator = injectOrError(serverConfiguratorKey)
  const dashboardConfigurators = injectOrError(dashboardConfiguratorsKey) as DataQueryConfigurator[]
  const branchConfigurator = injectOrError(branchConfiguratorKey)

  const rowsByBase = shallowRef(new Map<string, EngineCompareRow[]>())
  const loadingByBase = shallowRef(new Map<string, boolean>())
  // day key -> representative timestamp (ms) for labeling; union across bases and engines.
  const runTimeByDay = shallowRef(new Map<number, number>())
  const entries = new Map<string, BaseEntry>()

  const selectedBranches = (): string[] => {
    const list = selectedToArray(branchConfigurator?.selected.value)
    return list.length > 0 ? list : [""]
  }

  const measuresFor = (project: HighlightingProject): string[] => opts.selectedMetricTypes.value.map((type) => buildMeasure(project, type, opts.quantity.value.suffix))

  function setLoading(base: string, loading: boolean): void {
    if (loadingByBase.value.get(base) === loading) return
    loadingByBase.value.set(base, loading)
    triggerRef(loadingByBase)
  }

  // Record every run's day so the picker can list available runs, keeping the newest timestamp per day.
  function noteRunDays(runs: readonly SeriesRun[]): void {
    let changed = false
    for (const run of runs) {
      const day = dayKey(run.t)
      const existing = runTimeByDay.value.get(day)
      if (existing == null || run.t > existing) {
        runTimeByDay.value.set(day, run.t)
        changed = true
      }
    }
    if (changed) triggerRef(runTimeByDay)
  }

  function recompute(project: HighlightingProject, data: DataQueryResult, configuration: DataQueryExecutorConfiguration): void {
    const entry = entries.get(project.base)
    if (entry == null) return

    const types = opts.selectedMetricTypes.value
    const branches = selectedBranches()
    const measures = measuresFor(project)
    const indexed = indexSeriesRuns(data, configuration, branches, [entry.legacyProject, entry.newProject], measures)

    const rows: EngineCompareRow[] = []
    // types[i] and measures[i] are parallel (measuresFor maps the former to the latter), so reuse the
    // already-built measure name instead of recomputing it.
    for (let i = 0; i < types.length; i++) {
      const type = types[i]
      const measure = measures[i]
      // Merge across selected branches; the dashboard is single-branch by design, so this is a no-op
      // for the common case and a defensive concat if the user widens the branch selection.
      const legacyRuns = branches.flatMap((branch) => indexed.get(seriesKey(branch, entry.legacyProject, measure)) ?? [])
      const newRuns = branches.flatMap((branch) => indexed.get(seriesKey(branch, entry.newProject, measure)) ?? [])
      noteRunDays(legacyRuns)
      noteRunDays(newRuns)

      const values = pickRunValues(opts.mode.value === "single", opts.selectedRunDay.value, legacyRuns, newRuns)
      if (values == null || (values.legacy.length === 0 && values.new.length === 0)) continue

      const legacy = computeBaseStats(values.legacy)
      const newStats = computeBranchStats(legacy, values.new)
      const ratio = Number.isFinite(legacy.center) && legacy.center > 0 && Number.isFinite(newStats.center) ? newStats.center / legacy.center : Number.NaN
      rows.push({
        base: project.base,
        title: project.title,
        metricType: type,
        phase: phaseOf(type),
        bucket: bucketOf(type),
        measure,
        legacy,
        new: newStats,
        ratio,
        diffPercent: newStats.diffPercent,
        significant: isSignificantBranch(newStats),
      })
    }

    rowsByBase.value.set(project.base, rows)
    triggerRef(rowsByBase)
  }

  function recomputeFromCache(): void {
    for (const entry of entries.values()) {
      if (entry.lastData != null && entry.lastConfiguration != null) {
        recompute(entry.project, entry.lastData, entry.lastConfiguration)
      }
    }
  }

  for (const project of projects) {
    const legacyProject = `${project.base}/highlighting`
    const newProject = `${project.base}-types2/highlighting`
    const scenario = new DimensionConfigurator("project", true)
    scenario.selected.value = [legacyProject, newProject]
    const measures = shallowRef<string[]>(measuresFor(project))
    // "auto" value unit — declared units (sub-metrics) and the stored type (duration) resolve formatting.
    const measureConfigurator = new PredefinedMeasureConfigurator(measures, shallowRef(false), "line", "auto", {}, null, "item")
    const executor = new DataQueryExecutor([...dashboardConfigurators, scenario, serverConfigurator, measureConfigurator])
    const entry: BaseEntry = { project, legacyProject, newProject, measures, lastData: null, lastConfiguration: null, unsubscribe: () => {} }
    setLoading(project.base, true)
    entry.unsubscribe = executor.subscribe((data, configuration, loading) => {
      setLoading(project.base, loading)
      if (loading || data == null) return
      entry.lastData = data
      entry.lastConfiguration = configuration
      recompute(project, data, configuration)
    })
    entries.set(project.base, entry)
  }

  // Measures change (metric-type or quantity selection) triggers a fresh query per base. Empty metric
  // selection makes each measure configurator produce no query, so clear rows explicitly.
  watch(
    [opts.selectedMetricTypes, opts.quantity],
    () => {
      if (opts.selectedMetricTypes.value.length === 0) {
        if (rowsByBase.value.size > 0) {
          rowsByBase.value.clear()
          triggerRef(rowsByBase)
        }
        for (const base of entries.keys()) setLoading(base, false)
        return
      }
      for (const entry of entries.values()) {
        entry.measures.value = measuresFor(entry.project)
        setLoading(entry.project.base, true)
      }
    },
    { flush: "post" }
  )

  // Mode / run-day change reuses the cached query result — same data, different reduction.
  watch([opts.mode, opts.selectedRunDay], () => {
    recomputeFromCache()
  })

  onUnmounted(() => {
    for (const entry of entries.values()) entry.unsubscribe()
    entries.clear()
  })

  // Flatten in the canonical project order regardless of which executor resolved first.
  const rows = computed<EngineCompareRow[]>(() => projects.flatMap((project) => rowsByBase.value.get(project.base) ?? []))
  const aggregates = computed<EngineAggregates>(() => computeEngineAggregates(rows.value))
  const runDays = computed<RunDay[]>(() =>
    [...runTimeByDay.value.entries()]
      .toSorted(([a], [b]) => b - a)
      .map(([day, timestamp]) => ({ day, label: new Date(timestamp).toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" }) }))
  )
  const loading = computed<boolean>(() => [...loadingByBase.value.values()].some(Boolean))

  return { rows, aggregates, runDays, loading }
}
