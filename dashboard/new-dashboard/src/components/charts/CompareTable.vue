<template>
  <div class="flex flex-col gap-3 py-3 px-5 border border-solid rounded-md">
    <div class="flex items-center justify-between gap-4 flex-wrap">
      <div class="text-sm">
        <template v-if="selection == null">
          <span class="text-gray-500"
            >Select at least two branches in the toolbar to compare. The base is chosen automatically (release branch > master > first selected) so master shows up as a regression
            candidate against the stable release.</span
          >
        </template>
        <template v-else-if="!hasRows">
          <span class="text-gray-500"
            >Comparing against base <span class="font-semibold">{{ selection.base }}</span
            >. Loading…</span
          >
        </template>
        <template v-else>
          <span class="font-semibold">{{ significantCount }} of {{ totalRows }} tests </span>
          on this dashboard differ noticeably from <span class="font-semibold">{{ selection.base }}</span> for
          <span class="font-semibold">{{ selection.compared.join(", ") }}</span>
        </template>
      </div>
      <div class="flex items-center gap-3">
        <label class="flex items-center gap-1 text-sm">
          <Checkbox
            v-model="significantOnly"
            binary
          />
          Hide noisy changes
        </label>
      </div>
    </div>

    <DataTable
      :value="visibleRows"
      :loading="isLoading"
      show-gridlines
      class="p-datatable-sm"
      sort-field="absDiffPercent"
      :sort-order="-1"
    >
      <Column
        field="sectionLabel"
        header="Section"
        :sortable="true"
      />
      <Column
        field="project"
        header="Project"
        :sortable="true"
      >
        <template #body="slotProps">
          <div class="flex items-center">
            <div>{{ slotProps.data.project }}</div>
            <div class="ml-2">
              <Button
                icon="pi pi-external-link"
                class="p-button-rounded p-button-text p-button-sm"
                @click="() => handleNavigateToTest(slotProps.data)"
              />
            </div>
          </div>
        </template>
      </Column>
      <Column
        field="metric"
        header="Metric"
        :sortable="true"
      />
      <Column
        :header="baseHeader"
        :sortable="true"
        :sort-field="baseCenterSortField"
      >
        <template #body="slotProps">
          <div>{{ formatRowValue(slotProps.data, slotProps.data.base.center) }}</div>
          <div class="text-gray-500 text-xs">{{ formatRowRange(slotProps.data, slotProps.data.base.p10, slotProps.data.base.p90) }} · n={{ slotProps.data.base.count }}</div>
        </template>
      </Column>
      <template
        v-for="(branch, index) in comparedBranches"
        :key="branch"
      >
        <Column
          :header="branch"
          :sortable="true"
          :sort-field="centerSortFns[index]"
          :pt="cellTintPt(index)"
        >
          <template #body="slotProps">
            <template v-if="Number.isFinite(slotProps.data.compared[index].center)">
              <div>{{ formatRowValue(slotProps.data, slotProps.data.compared[index].center) }}</div>
              <div class="text-gray-500 text-xs">n={{ slotProps.data.compared[index].count }}</div>
            </template>
            <span
              v-else
              class="text-gray-400"
              >—</span
            >
          </template>
        </Column>
        <Column
          header="Δ%"
          :sortable="true"
          :sort-field="diffSortFns[index]"
          :pt="cellTintPt(index)"
        >
          <template #body="slotProps">
            <span v-if="Number.isFinite(slotProps.data.compared[index].diffPercent)">
              {{ formatSigned(slotProps.data.compared[index].diffPercent, "%") }}
            </span>
            <span
              v-else
              class="text-gray-400"
              >—</span
            >
          </template>
        </Column>
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { deepEqual } from "fast-equals"
import { computed, onUnmounted, ref, shallowRef, triggerRef, watch } from "vue"
import { useRouter } from "vue-router"
import { ValueUnit } from "../common/chart"
import { DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { DataQueryExecutor, DataQueryResult } from "../common/DataQueryExecutor"
import { PredefinedMeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { DimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { injectOrError } from "../../shared/injectionKeys"
import { branchConfiguratorKey, compareSectionsRegistryKey, dashboardConfiguratorsKey, serverConfiguratorKey } from "../../shared/keys"
import { openTestDrilldown } from "../../util/testDrilldown"
import { BaseAndCompared, CompareSectionConfig, pickBaseAndCompared } from "./compareMode"
import { indexSeries, seriesKey } from "./compareQuery"
import { BaseStats, BranchStats, computeBaseStats, computeBranchStats, DISPARITY_SIGNIFICANT_THRESHOLD } from "./compareStats"
import { formatMeasureValue, resolveMeasureUnit } from "../common/formatter"
import "./compareCells.css"

interface CompareRow {
  key: string
  sectionId: string
  sectionLabel: string
  project: string
  metric: string
  valueUnit: ValueUnit
  base: BaseStats
  compared: BranchStats[] // ordered to match selection.compared
  absDiffPercent: number // max |Δ%| across compared branches; drives the default desc sort
  hasSignificant: boolean
}

// PrimeVue 4's `sortField` accepts a function returning the sort key; its .d.ts narrows the
// return to `string` but the implementation passes whatever the function returns straight to
// the comparator, so numeric keys work. Cast through the declared signature.
type ColumnSortField = (item: CompareRow) => unknown
const asSortField = (fn: ColumnSortField): ((item: any) => string) => fn as unknown as (item: any) => string

const serverConfigurator = injectOrError(serverConfiguratorKey)
const dashboardConfigurators = injectOrError(dashboardConfiguratorsKey) as DataQueryConfigurator[]
const registry = injectOrError(compareSectionsRegistryKey)
const branchConfiguratorOpt = injectOrError(branchConfiguratorKey)
if (branchConfiguratorOpt == null) {
  // DashboardPage.canCompare gates compare mode on branchConfigurator != null, so this should be unreachable.
  throw new Error("CompareTable mounted without a BranchConfigurator")
}
const branchConfigurator = branchConfiguratorOpt

const router = useRouter()
const sectionRows = shallowRef(new Map<string, CompareRow[]>())
const sectionLoading = shallowRef(new Map<string, boolean>())
const significantOnly = ref(true)

const selectedBranches = computed<string[]>(() => {
  const sel = branchConfigurator.selected.value
  if (sel == null) return []
  return Array.isArray(sel) ? sel : [sel]
})

const selection = computed<BaseAndCompared | null>(() => pickBaseAndCompared(selectedBranches.value))
const baseHeader = computed(() => selection.value?.base ?? "base")
const comparedBranches = computed<string[]>(() => selection.value?.compared ?? [])

const baseCenterSortField = asSortField((row) => row.base.center)
function sortByBranchStat<K extends keyof BranchStats>(key: K) {
  return computed(() => comparedBranches.value.map((_, i) => asSortField((row) => (row.compared[i]?.[key] as number | undefined) ?? Number.NaN)))
}
const centerSortFns = sortByBranchStat("center")
const diffSortFns = sortByBranchStat("diffPercent")

const rowsState = computed<{ all: CompareRow[]; significant: CompareRow[] }>(() => {
  const all: CompareRow[] = []
  const significant: CompareRow[] = []
  for (const section of registry.sections.value) {
    const rows = sectionRows.value.get(section.id)
    if (rows == null) continue
    for (const row of rows) {
      all.push(row)
      if (row.hasSignificant) significant.push(row)
    }
  }
  return { all, significant }
})

const visibleRows = computed(() => (significantOnly.value ? rowsState.value.significant : rowsState.value.all))
const hasRows = computed(() => rowsState.value.all.length > 0)
const totalRows = computed(() => rowsState.value.all.length)
const significantCount = computed(() => rowsState.value.significant.length)
const isLoading = computed(() => {
  for (const v of sectionLoading.value.values()) {
    if (v) return true
  }
  return false
})

interface SubscriptionEntry {
  config: CompareSectionConfig
  // Built once per section and reused across branch / filter changes. MachineConfigurator and
  // DimensionConfigurator both hold long-lived rxjs subscriptions when constructed, so we keep
  // them alive for the section's lifetime rather than recreating on each branch change.
  scenario: DimensionConfigurator
  measure: PredefinedMeasureConfigurator
  machine: MachineConfigurator | null
  // One executor fans out across all selected branches via the live branchConfigurator's own
  // queryProducer — same machinery the chart-mode pipeline uses. The executor reacts to branch
  // changes on its own; we don't rebuild it on selection swaps.
  unsubscribe: () => void
  warnedNames: Set<string>
}

const liveSubscriptions = new Map<string, SubscriptionEntry>()

function toArray(measure: string | string[]): string[] {
  return Array.isArray(measure) ? measure : [measure]
}

function buildExecutor(entry: SubscriptionEntry): DataQueryExecutor {
  // Reuse the full dashboard configurator list, including the live branchConfigurator.
  // BranchConfigurator.configureQuery pushes a queryProducer that fans the same query out
  // across every selected branch and labels each series with the branch name, so a single
  // executor produces base + compared series in one round trip.
  const configurators: DataQueryConfigurator[] = [...dashboardConfigurators]
  if (entry.machine != null) configurators.push(entry.machine)
  configurators.push(entry.scenario, serverConfigurator, entry.measure)
  return new DataQueryExecutor(configurators)
}

function startSection(section: CompareSectionConfig): void {
  stopSection(section.id)

  const measures = toArray(section.measure)
  // Bare DimensionConfigurator — no factory call, so no server fetch for valid projects and
  // no rxjs leak. The producer logic in DimensionConfigurator.configureQuery is the same one
  // chart-mode relies on, just without the dropdown-population side effect.
  const scenario = new DimensionConfigurator("project", true)
  scenario.selected.value = [...section.projects]

  const entry: SubscriptionEntry = {
    config: section,
    scenario,
    measure: new PredefinedMeasureConfigurator(ref(measures), ref(false), "line", section.valueUnit, {}, null, "item"),
    machine: section.machines != null && section.machines.length > 0 ? new MachineConfigurator(serverConfigurator, undefined, [], true, section.machines) : null,
    unsubscribe: () => {},
    warnedNames: new Set(),
  }
  liveSubscriptions.set(section.id, entry)
  setSectionLoading(section.id, true)

  const executor = buildExecutor(entry)
  entry.unsubscribe = executor.subscribe((data, configuration, loading) => {
    setSectionLoading(section.id, loading)
    if (loading || data == null) return
    recompute(entry, data, configuration)
  })
}

function stopSection(id: string): void {
  const entry = liveSubscriptions.get(id)
  if (entry != null) {
    entry.unsubscribe()
    liveSubscriptions.delete(id)
  }
  if (sectionRows.value.delete(id)) triggerRef(sectionRows)
  if (sectionLoading.value.delete(id)) triggerRef(sectionLoading)
}

function setSectionLoading(id: string, loading: boolean): void {
  if (sectionLoading.value.get(id) === loading) return
  sectionLoading.value.set(id, loading)
  triggerRef(sectionLoading)
}

function recompute(entry: SubscriptionEntry, data: DataQueryResult, configuration: DataQueryExecutorConfiguration): void {
  const section = entry.config

  const sel = selection.value
  if (sel == null) return // executor still alive but the user has narrowed the selection — nothing to render.

  const measures = toArray(section.measure)
  const branches = [sel.base, ...sel.compared]
  const indexed = indexSeries(data, configuration, branches, section.projects, measures)

  for (const name of indexed.unresolvedNames) {
    if (entry.warnedNames.has(name)) continue
    entry.warnedNames.add(name)
    console.warn(
      `[CompareTable] Could not resolve series name "${name}" against section "${section.label}"; falling back to positional match. ` +
        `If rows look mis-attributed, check that the dashboard's branch/project/measure names round-trip through the data pipeline.`
    )
  }

  const rows: CompareRow[] = []
  for (const project of section.projects) {
    for (const metric of measures) {
      const baseValues = indexed.byKey.get(seriesKey(sel.base, project, metric)) ?? []
      const comparedValues = sel.compared.map((b) => indexed.byKey.get(seriesKey(b, project, metric)) ?? [])
      const anyData = baseValues.length > 0 || comparedValues.some((v) => v.length > 0)
      if (!anyData) continue

      const baseStats = computeBaseStats(baseValues)
      const branchStats = comparedValues.map((values) => computeBranchStats(baseStats, values))

      let absDiffPercent = 0
      let hasSignificant = false
      for (const s of branchStats) {
        if (Number.isFinite(s.diffPercent)) {
          const pv = Math.abs(s.diffPercent)
          if (pv > absDiffPercent) absDiffPercent = pv
        }
        if (!hasSignificant && isSignificantBranch(s)) hasSignificant = true
      }

      rows.push({
        key: `${section.id}::${project}::${metric}`,
        sectionId: section.id,
        sectionLabel: section.label,
        project,
        metric,
        valueUnit: section.valueUnit,
        base: baseStats,
        compared: branchStats,
        absDiffPercent,
        hasSignificant,
      })
    }
  }

  sectionRows.value.set(section.id, rows)
  triggerRef(sectionRows)
}

function isSignificantBranch(stats: BranchStats): boolean {
  // Two gates, both required:
  //   - |D| ≥ threshold filters out statistical noise (tight-baseline rows that aren't real)
  //   - |Δ%| ≥ threshold filters out practically-irrelevant changes (0.0 %/0.1 % at very high D)
  // ±Infinity D paired with any finite non-zero Δ% still counts — flat baseline + any real change.
  if (Number.isNaN(stats.disparity)) return false
  if (Math.abs(stats.disparity) < DISPARITY_SIGNIFICANT_THRESHOLD) return false
  if (!Number.isFinite(stats.diffPercent)) return false
  return Math.abs(stats.diffPercent) >= DIFF_PERCENT_WARN
}

interface ImpactSeverity {
  direction: "degradation" | "improvement"
  severity: "severe" | "warn"
}

function diffPercentSeverity(diffPercent: number | undefined): ImpactSeverity | null {
  if (diffPercent == null || !Number.isFinite(diffPercent) || diffPercent === 0) return null
  const abs = Math.abs(diffPercent)
  if (abs < DIFF_PERCENT_WARN) return null
  return {
    direction: diffPercent > 0 ? "degradation" : "improvement",
    severity: abs >= DIFF_PERCENT_SEVERE ? "severe" : "warn",
  }
}

// Tint the cell's <td> directly via PrimeVue 4's pt API; `parent.props.rowData` is the row
// for this body cell. Improvements get green, degradations red; sub-threshold or missing
// Δ% leaves the cell untinted.
function cellTintPt(comparedIndex: number) {
  return {
    bodyCell: ({ parent }: { parent: { props: { rowData?: CompareRow } } }) => {
      const diffPercent = parent.props.rowData?.compared[comparedIndex]?.diffPercent
      const sev = diffPercentSeverity(diffPercent)
      return { class: sev == null ? undefined : `compare-cell-${sev.direction}-${sev.severity}` }
    },
  }
}

function formatSigned(value: number, suffix: string): string {
  return `${value >= 0 ? "+" : ""}${value.toFixed(1)}${suffix}`
}

// Mirrors LineChart's valueUnitFromMeasure: a metric ending in `.ns` or `.ms` carries its
// own unit hint and overrides the section-level default (which is `ms` unless explicitly set).
function effectiveValueUnit(row: CompareRow): ValueUnit {
  if (row.metric.endsWith(".ms")) return "ms"
  if (row.metric.endsWith(".ns")) return "ns"
  return row.valueUnit
}

function formatRowValue(row: CompareRow, value: number): string {
  if (!Number.isFinite(value)) return "—"
  return formatMeasureValue(value, resolveMeasureUnit(row.metric, { valueUnit: effectiveValueUnit(row) }))
}

function formatRowRange(row: CompareRow, lo: number, hi: number): string {
  if (!Number.isFinite(lo) || !Number.isFinite(hi)) return "—"
  return `${formatRowValue(row, lo)} – ${formatRowValue(row, hi)}`
}

function handleNavigateToTest(row: CompareRow): void {
  const sel = selection.value
  const branches = sel == null ? null : [sel.base, ...sel.compared]
  openTestDrilldown(router, { project: row.project, measure: row.metric, branches })
}

// Section lifecycle is driven only by the registry and by whether there's a usable selection.
// The executor itself reacts to branch changes via the live branchConfigurator's observable —
// no manual rebuild needed when the user picks a different branch. We still drop cached rows
// on selection change, because column positions are derived from sel.compared[] and a stale
// row would briefly show the old branch's numbers under the new branch's header.
watch(
  [() => registry.sections.value, selection],
  ([sections, sel], prev) => {
    const wantedIds = sel == null ? new Set<string>() : new Set(sections.map((s) => s.id))
    for (const id of liveSubscriptions.keys()) {
      if (!wantedIds.has(id)) stopSection(id)
    }
    if (sel == null) return

    const prevSel = prev?.[1] ?? null
    const selectionChanged = prevSel != null && !deepEqual(sel, prevSel)
    for (const section of sections) {
      const existing = liveSubscriptions.get(section.id)
      if (existing == null || !deepEqual(existing.config, section)) {
        startSection(section)
      } else if (selectionChanged) {
        // Same executor, but the comparison columns shift; clear cached rows so we don't
        // render stale (old-branch) data under the new headers until fresh data arrives.
        if (sectionRows.value.delete(section.id)) triggerRef(sectionRows)
        setSectionLoading(section.id, true)
        existing.warnedNames.clear()
      }
    }
  },
  { immediate: true }
)

onUnmounted(() => {
  for (const id of liveSubscriptions.keys()) {
    stopSection(id)
  }
})
</script>
