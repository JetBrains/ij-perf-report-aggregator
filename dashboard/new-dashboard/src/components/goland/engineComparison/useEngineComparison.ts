import { computed, inject, onUnmounted, Ref, shallowRef, triggerRef, watch } from "vue"
import { DimensionConfigurator, selectedToArray } from "../../../configurators/DimensionConfigurator"
import { PredefinedMeasureConfigurator } from "../../../configurators/MeasureConfigurator"
import { defaultModeName, TestModeConfigurator } from "../../../configurators/TestModeConfigurator"
import { injectOrError } from "../../../shared/injectionKeys"
import { branchConfiguratorKey, dashboardConfiguratorsKey, serverConfiguratorKey, testModeConfiguratorKey } from "../../../shared/keys"
import { DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../common/dataQuery"
import { DataQueryExecutor, DataQueryResult } from "../../common/DataQueryExecutor"
import { formatMeasureValue, resolveMeasureUnit } from "../../common/formatter"
import { indexSeriesRuns, SeriesRun, seriesKey } from "../../charts/compareQuery"
import { BaseStats, BranchStats } from "../../charts/compareStats"
import { bucketOf, buildMeasure, HighlightingProject, phaseOf, projects, Quantity } from "./highlightingMetrics"
import { computeEngineAggregates, EngineAggregates, EngineComparePoint, EngineGap } from "./engineCompareAggregates"
import { compareCell, countMissing, pairRuns } from "./engineRunSelection"

// How the comparison reduces the runs in the window: "aggregated" = robust center over every run
// (SPEC-style verdict), "single" = one run's before/after (see EngineComparisonMode below).
export type EngineComparisonMode = "aggregated" | "single"

// A selectable run in single-run mode. `t` is the build timestamp (ms) shared by every project the
// build published; `label` is its human date and time.
export interface Run {
  t: number
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
  // In single-run mode, the build timestamp to show; null means the newest build in any project.
  selectedRun: Readonly<Ref<number | null>>
}

interface EngineComparison {
  rows: Ref<EngineCompareRow[]>
  aggregates: Ref<EngineAggregates>
  runs: Ref<Run[]>
  loading: Ref<boolean>
}

// The runs of one (base, metric type) cell for each engine, as the last query returned them.
interface CellRuns {
  type: string
  measure: string
  legacy: SeriesRun[]
  new: SeriesRun[]
}

// One live query per project base, mirroring today's per-project GroupProjectsChart ownership. Measure
// names embed the project's scenario file, so the bases cannot share a single executor.
interface BaseEntry {
  project: HighlightingProject
  legacyProject: string
  newProject: string
  measures: Ref<string[]>
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

  // Replaced on each result, so a narrower time range or another branch drops the builds it no longer has.
  const cellsByBase = shallowRef(new Map<string, CellRuns[]>())
  const loadingByBase = shallowRef(new Map<string, boolean>())
  // Every build that published any selected cell of any project for either engine. A build in which one
  // engine failed everywhere still belongs here, so its missing results are counted.
  const runTimes = computed(() => new Set([...cellsByBase.value.values()].flatMap((cells) => cells.flatMap((cell) => [...cell.legacy, ...cell.new].map((run) => run.t)))))
  // The build a null run selection resolves to: one build for every cell, so a "latest" summary is one snapshot.
  const newestRun = computed<number | null>(() => (runTimes.value.size === 0 ? null : Math.max(...runTimes.value)))
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

  function storeCells(entry: BaseEntry, data: DataQueryResult, configuration: DataQueryExecutorConfiguration): void {
    const branches = selectedBranches()
    const measures = measuresFor(entry.project)
    const indexed = indexSeriesRuns(data, configuration, branches, [entry.legacyProject, entry.newProject], measures)
    const runsOf = (variant: string, measure: string) => branches.flatMap((branch) => indexed.get(seriesKey(branch, variant, measure)) ?? [])
    // types[i] and measures[i] are parallel (measuresFor maps the former to the latter).
    const cells = opts.selectedMetricTypes.value.map((type, i) => ({
      type,
      measure: measures[i],
      legacy: runsOf(entry.legacyProject, measures[i]),
      new: runsOf(entry.newProject, measures[i]),
    }))
    cellsByBase.value.set(entry.project.base, cells)
    triggerRef(cellsByBase)
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
    const entry: BaseEntry = { project, legacyProject, newProject, measures, unsubscribe: () => {} }
    setLoading(project.base, true)
    entry.unsubscribe = executor.subscribe((data, configuration, loading) => {
      setLoading(project.base, loading)
      if (loading || data == null) return
      storeCells(entry, data, configuration)
    })
    entries.set(project.base, entry)
  }

  // Measures change (metric-type or quantity selection) triggers a fresh query per base. Empty metric
  // selection makes each measure configurator produce no query, so clear the cells explicitly.
  watch(
    [opts.selectedMetricTypes, opts.quantity],
    () => {
      if (opts.selectedMetricTypes.value.length === 0) {
        if (cellsByBase.value.size > 0) {
          cellsByBase.value.clear()
          triggerRef(cellsByBase)
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

  onUnmounted(() => {
    for (const entry of entries.values()) entry.unsubscribe()
    entries.clear()
  })

  // Rows and gaps in the canonical project order regardless of which executor resolved first. The
  // scope is the selected build in single mode and every build of the window in aggregated mode.
  const comparison = computed<{ rows: EngineCompareRow[]; gaps: EngineGap[] }>(() => {
    const single = opts.mode.value === "single"
    const run = opts.selectedRun.value ?? newestRun.value
    const scope = single ? (run == null ? [] : [run]) : runTimes.value
    const rows: EngineCompareRow[] = []
    const gaps: EngineGap[] = []
    for (const project of projects) {
      for (const cell of cellsByBase.value.get(project.base) ?? []) {
        const missing = countMissing(scope, cell.legacy, cell.new)
        if (missing.new > 0) gaps.push({ base: project.base, title: project.title, metricType: cell.type, engine: "NEW", builds: missing.new })
        if (missing.legacy > 0) gaps.push({ base: project.base, title: project.title, metricType: cell.type, engine: "LEGACY", builds: missing.legacy })
        const stats = compareCell(single, run, pairRuns(cell.legacy, cell.new))
        if (stats == null) continue
        rows.push({
          base: project.base,
          title: project.title,
          metricType: cell.type,
          phase: phaseOf(cell.type),
          bucket: bucketOf(cell.type),
          measure: cell.measure,
          ...stats,
        })
      }
    }
    return { rows, gaps }
  })

  const rows = computed<EngineCompareRow[]>(() => comparison.value.rows)
  // The metric types of the stored results, not of the live selection: until a refetch lands, the
  // rows and gaps still describe the previous selection.
  const storedTypes = computed(() => new Set([...cellsByBase.value.values()].flatMap((cells) => cells.map((cell) => cell.type))))
  const aggregates = computed<EngineAggregates>(() => computeEngineAggregates(comparison.value.rows, comparison.value.gaps, projects.length * storedTypes.value.size))
  const runs = computed<Run[]>(() =>
    [...runTimes.value]
      .toSorted((a, b) => b - a)
      .map((t) => ({ t, label: new Date(t).toLocaleString(undefined, { year: "numeric", month: "short", day: "numeric", hour: "2-digit", minute: "2-digit" }) }))
  )
  const loading = computed<boolean>(() => [...loadingByBase.value.values()].some(Boolean))

  return { rows, aggregates, runs, loading }
}
