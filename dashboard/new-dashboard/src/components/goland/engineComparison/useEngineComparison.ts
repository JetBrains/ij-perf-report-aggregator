import { computed, inject, onUnmounted, Ref, shallowRef, triggerRef, watch } from "vue"
import { DimensionConfigurator, selectedToArray } from "../../../configurators/DimensionConfigurator"
import { PredefinedMeasureConfigurator } from "../../../configurators/MeasureConfigurator"
import { defaultModeName, TestModeConfigurator } from "../../../configurators/TestModeConfigurator"
import { injectOrError } from "../../../shared/injectionKeys"
import { branchConfiguratorKey, dashboardConfiguratorsKey, serverConfiguratorKey, testModeConfiguratorKey } from "../../../shared/keys"
import { DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../common/dataQuery"
import { DataQueryExecutor, DataQueryResult } from "../../common/DataQueryExecutor"
import { formatMeasureValue, resolveMeasureUnit } from "../../common/formatter"
import { indexSeriesRuns, seriesKey } from "../../charts/compareQuery"
import { BaseStats, BranchStats } from "../../charts/compareStats"
import { bucketOf, buildMeasure, HighlightingProject, phaseOf, projects, Quantity } from "./highlightingMetrics"
import { computeEngineAggregates, EngineAggregates, EngineComparePoint } from "./engineCompareAggregates"
import { compareCell, dayKey, pairRuns } from "./engineRunSelection"

// How the comparison reduces the runs in the window: "aggregated" = robust center over every run
// (SPEC-style verdict), "single" = one run's before/after (see EngineComparisonMode below).
export type EngineComparisonMode = "aggregated" | "single"

// A selectable run in single-run mode. `day` is the UTC-day key that groups the paired builds; `label`
// is the human date.
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
  // In single-run mode, the run day to show; null means the newest day with a paired build in any project.
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
  const dashboardModeConfigurator = inject(testModeConfiguratorKey, null)

  // The two engines are compared within one test mode: the default mode when it is selected, else the first selected
  // mode. Without the filter, the runs of every mode join the series of the engine project they publish under.
  const modeConfigurator = new TestModeConfigurator(false)
  watch(
    () => selectedToArray(dashboardModeConfigurator?.selected.value),
    (modes) => {
      modeConfigurator.selected.value = modes.length === 0 || modes.includes(defaultModeName) ? defaultModeName : modes[0]
    },
    { immediate: true }
  )
  const queryConfigurators = [...dashboardConfigurators.filter((configurator) => configurator !== dashboardModeConfigurator), modeConfigurator]

  const rowsByBase = shallowRef(new Map<string, EngineCompareRow[]>())
  const loadingByBase = shallowRef(new Map<string, boolean>())
  // Per base: day key -> newest paired build timestamp (ms) of that day. Replaced on each result, so a
  // narrower time range or another branch drops the days it no longer has.
  const runTimeByDayByBase = shallowRef(new Map<string, Map<number, number>>())
  const runTimeByDay = computed(() => {
    const union = new Map<number, number>()
    for (const byDay of runTimeByDayByBase.value.values()) {
      for (const [day, t] of byDay) union.set(day, Math.max(t, union.get(day) ?? t))
    }
    return union
  })
  // The day a null run selection resolves to: one day for every cell, so a "latest" summary is one snapshot.
  const newestDay = computed<number | null>(() => (runTimeByDay.value.size === 0 ? null : Math.max(...runTimeByDay.value.keys())))
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

  function recompute(project: HighlightingProject, data: DataQueryResult, configuration: DataQueryExecutorConfiguration): void {
    const entry = entries.get(project.base)
    if (entry == null) return

    const types = opts.selectedMetricTypes.value
    const branches = selectedBranches()
    const measures = measuresFor(project)
    const indexed = indexSeriesRuns(data, configuration, branches, [entry.legacyProject, entry.newProject], measures)

    // types[i] and measures[i] are parallel (measuresFor maps the former to the latter), so reuse the
    // already-built measure name instead of recomputing it. Days are noted for all cells before any cell
    // resolves the newest day.
    const baseDays = new Map<number, number>()
    const cells = types.map((type, i) => {
      const measure = measures[i]
      const runsOf = (variant: string) => branches.flatMap((branch) => indexed.get(seriesKey(branch, variant, measure)) ?? [])
      const pairs = pairRuns(runsOf(entry.legacyProject), runsOf(entry.newProject))
      for (const pair of pairs) baseDays.set(dayKey(pair.t), Math.max(pair.t, baseDays.get(dayKey(pair.t)) ?? 0))
      return { type, measure, pairs }
    })
    runTimeByDayByBase.value.set(project.base, baseDays)
    triggerRef(runTimeByDayByBase)

    const single = opts.mode.value === "single"
    const day = opts.selectedRunDay.value ?? newestDay.value
    const rows: EngineCompareRow[] = []
    for (const { type, measure, pairs } of cells) {
      const cell = compareCell(single, day, pairs)
      if (cell == null) continue
      rows.push({
        base: project.base,
        title: project.title,
        metricType: type,
        phase: phaseOf(type),
        bucket: bucketOf(type),
        measure,
        ...cell,
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
    const executor = new DataQueryExecutor([...queryConfigurators, scenario, serverConfigurator, measureConfigurator])
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

  // A later base can bring a newer day; cells resolved to an older "latest" day must move to it.
  watch(newestDay, () => {
    if (opts.mode.value === "single" && opts.selectedRunDay.value == null) recomputeFromCache()
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
