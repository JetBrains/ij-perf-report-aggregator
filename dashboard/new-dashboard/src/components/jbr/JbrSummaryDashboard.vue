<template>
  <div class="flex flex-col gap-5">
    <StickyToolbar>
      <template #start>
        <div class="flex items-center gap-3">
          <BranchSelect
            :branch-configurator="branchConfigurator"
            :selection-limit="1"
          />
          <SelectButton
            v-model="lookbackDays"
            :options="lookbacks"
            option-label="label"
            option-value="days"
            :allow-empty="false"
            aria-label="The period whose runs the change detection searches"
          />
          <SelectButton
            v-model="selectedPlatforms"
            :options="platforms"
            multiple
            :allow-empty="false"
            aria-label="The platforms whose configurations the summary shows"
          />
          <i
            role="button"
            tabindex="0"
            aria-label="How the summary is computed"
            :aria-expanded="infoOpen"
            class="pi pi-info-circle cursor-pointer"
            @click="toggleInfo"
            @keydown.enter="toggleInfo"
            @keydown.space.prevent="toggleInfo"
          />
        </div>
      </template>
    </StickyToolbar>

    <Popover
      ref="infoPanel"
      append-to="body"
      @show="infoOpen = true"
      @hide="infoOpen = false"
    >
      <div class="flex flex-col gap-2 text-sm max-w-md">
        <span class="font-semibold">The latest significant change per benchmark and configuration</span>
        <span>Each run is summarized by an index: the geometric mean of its scenarios, each relative to its median over the period.</span>
        <span>
          Change point detection splits the period's runs where the index shifts. A cell compares the runs after the latest split with those before it (Hodges–Lehmann estimate) and
          is colored if its whole {{ levelOf(ALPHA) }} % interval lies beyond a change of {{ minEffectPercent }} %: the change is confidently at least that large. Positive is
          better. The charts of the benchmark tabs mark change points by the same rule, without waiting for more runs.
        </span>
        <span>"Worst" marks a scenario that regressed significantly on its own (Bonferroni-corrected), which the mean can hide.</span>
        <span>
          "?" marks a split with fewer than {{ MIN_RUNS_AFTER }} runs after it: too few to tell a change from noise, so the cell shows the latest change before it until more runs
          arrive.
        </span>
        <span class="text-xs opacity-80">The detector picks the split that separates the runs best, so the intervals are somewhat optimistic.</span>
      </div>
    </Popover>

    <p
      v-if="error != null"
      class="text-sm text-red-700 m-0"
    >
      {{ error }}
    </p>

    <div class="overflow-x-auto">
      <table class="summary border-collapse text-sm">
        <thead>
          <tr>
            <th rowspan="2" />
            <th
              v-for="group in columnGroups"
              :key="group.platform"
              :colspan="group.configurations.length"
              class="platform"
            >
              {{ group.platform }}
            </th>
          </tr>
          <tr>
            <th
              v-for="configuration in columns"
              :key="configuration"
              class="configuration"
            >
              {{ withBreaks(configuration) }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="row in grid"
            :key="row.definition.label"
          >
            <th class="benchmark">
              <RouterLink :to="row.definition.route">{{ row.definition.label }}</RouterLink>
              <span
                v-if="!loading && row.cells.size === 0"
                class="ml-2 text-xs font-normal text-gray-400"
                >no data in the last {{ lookbackDays }} days</span
              >
            </th>
            <td
              v-for="{ configuration, cell } in row.entries"
              :key="configuration"
            >
              <button
                v-if="cell != null"
                type="button"
                :class="['cell', classOf(cell.latest), isSelected(row.definition, configuration) ? 'selected' : '']"
                :title="describe(cell)"
                @click="select(row.definition, configuration)"
              >
                <template v-if="cell.latest != null">
                  {{ arrowOf(cell.latest) }} {{ formatChange(cell.latest.change) }}
                  <span class="date">{{ formatDate(cell.latest.time) }}</span>
                </template>
                <template v-else>≈ stable</template>
                <span
                  v-if="cell.worst != null"
                  class="worst"
                  >worst {{ formatChange(cell.worst.change) }}</span
                >
                <span
                  v-if="cell.pending != null"
                  class="pending"
                  >? {{ formatDate(cell.pending.time) }}<br />{{ cell.pending.runsAfter }} of {{ MIN_RUNS_AFTER }} runs</span
                >
              </button>
              <span
                v-else-if="loading"
                class="text-gray-400"
                >…</span
              >
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <section
      v-if="selectedCell != null"
      class="flex flex-col gap-2"
    >
      <h3 class="m-0 text-base">
        {{ selectedCell.row.label }} on {{ selectedCell.configuration }}: {{ describe(selectedCell.cell) }},
        <RouterLink :to="selectedCell.row.route">open the {{ selectedCell.row.label.split(" ")[0] }} tab</RouterLink>
      </h3>
      <p
        v-if="selectedCell.cell.changePoints.length > 0"
        class="m-0 text-sm"
      >
        Change points:
        <span
          v-for="point in selectedCell.cell.changePoints"
          :key="point.time"
          :class="['point', point === selectedCell.cell.pending ? 'pending' : classOf(point)]"
          >{{ point === selectedCell.cell.pending ? "?" : "" }} {{ formatDate(point.time) }} {{ formatChange(point.change) }}</span
        >
      </p>
      <table
        v-if="selectedCell.cell.scenarios.length > 0"
        class="scenarios border-collapse text-sm self-start"
      >
        <thead>
          <tr>
            <th>Scenario</th>
            <th>Median before</th>
            <th>Median after</th>
            <th>Change</th>
            <th>{{ levelOf(selectedCell.cell.scenarioMisrate) }} % interval</th>
            <th>Runs</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="scenario in selectedCell.cell.scenarios"
            :key="scenario.scenario"
          >
            <td>{{ scenario.scenario }}</td>
            <td class="number">{{ formatMeasureValue(scenario.medianBefore, selectedCell.row.unit) }}</td>
            <td class="number">{{ formatMeasureValue(scenario.medianAfter, selectedCell.row.unit) }}</td>
            <td :class="['number', classOf(scenario)]">{{ arrowOf(scenario) }} {{ formatChange(scenario.change) }}</td>
            <td class="number">{{ formatInterval(scenario) }}</td>
            <td class="number">{{ scenario.runsBefore }} / {{ scenario.runsAfter }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useStorage } from "@vueuse/core"
import { PopoverMethods } from "openvue/popover"
import { computed, ref, useTemplateRef, watch } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { selectedToArray } from "../../configurators/DimensionConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { refToObservable } from "../../configurators/rxjs"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import BranchSelect from "../common/BranchSelect.vue"
import { DataQuery } from "../common/dataQuery"
import { formatMeasureValue } from "../common/formatter"
import { PersistentStateManager } from "../common/PersistentStateManager"
import StickyToolbar from "../common/StickyToolbar.vue"
import { platformOf } from "./platforms"
import { ALPHA, Comparison, compareConfigurationNames, MIN_EFFECT, MIN_RUNS_AFTER, RunValue, SummaryCell, summarizeRow, SummaryRow } from "./summary"

/**
 * The first JBR tab: one row per benchmark (Render once per mode), one column per configuration, each cell the latest
 * significant change point of the benchmark there. The scenarios of a cell and the link to the benchmark's own tab are
 * one click away.
 */
const rows: SummaryRow[] = [
  { label: "Render onscreen", benchmark: "Render", route: "/jbr/renderDashboard", unit: "milliseconds", betterDirection: "lower", measureSuffix: "_onscreen" },
  { label: "Render volatile", benchmark: "Render", route: "/jbr/renderDashboard", unit: "milliseconds", betterDirection: "lower", measureSuffix: "_volatile" },
  { label: "Render buffer", benchmark: "Render", route: "/jbr/renderDashboard", unit: "milliseconds", betterDirection: "lower", measureSuffix: "_buffer" },
  { label: "J2DBench", benchmark: "J2DBench", route: "/jbr/j2dDashboard", unit: "pixelsPerSecond", betterDirection: "higher" },
  { label: "JavaDraw", benchmark: "JavaDraw", route: "/jbr/javaDrawDashboard", unit: "framesPerSecond", betterDirection: "higher" },
  { label: "MapBench", benchmark: "Mapbench", route: "/jbr/mapbenchDashboard", unit: "framesPerSecond", betterDirection: "higher" },
  { label: "SwingMark", benchmark: "SwingMark", route: "/jbr/swingmarkDashboard", unit: "milliseconds", betterDirection: "lower" },
  { label: "DaCapo", benchmark: "Dacapo", route: "/jbr/dacapoDashboard", unit: "milliseconds", betterDirection: "lower" },
  { label: "SPECjbb2015", benchmark: "SPECjbb2015", route: "/jbr/specDashboard", unit: "javaOperationsPerSecond", betterDirection: "higher" },
]

// The change detection needs at least five runs on each side of a change; Render runs about five times a week.
const lookbacks = [
  { days: 30, label: "30 days" },
  { days: 60, label: "60 days" },
  { days: 90, label: "90 days" },
]
const platforms = ["macOS", "Linux", "Windows"]
const minEffectPercent = MIN_EFFECT * 100

const lookbackDays = useStorage<number>("jbr_summary_lookback", 60)
const selectedPlatforms = useStorage<string[]>("jbr_summary_platforms", [...platforms])

const infoPanel = useTemplateRef<PopoverMethods>("infoPanel")
const infoOpen = ref(false)
function toggleInfo(event: Event) {
  infoPanel.value?.toggle(event)
}

const serverConfigurator = new ServerWithCompressConfigurator("jbr", "report")

// The branches offered are those with data in the period; dev_main holds the JBR CI runs.
const persistence = new PersistentStateManager("jbr_summary", { branch: "dev_main" }, useRouter())
const period: FilterConfigurator = {
  configureFilter(query: DataQuery): boolean {
    query.addFilter({ f: "generated_time", q: `>subtractDays(now(),${lookbackDays.value})` })
    return true
  },
  createObservable() {
    return refToObservable(lookbackDays)
  },
}
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistence, [period])
const branch = computed(() => selectedToArray(branchConfigurator.selected.value)[0] ?? "dev_main")

const values = ref(new Map<string, RunValue[]>())
const loading = ref(false)
const error = ref<string | null>(null)

/** Every run's value of every scenario of the benchmark in the period; the statistics measures ("#") left out. */
async function loadValues(benchmark: string, days: number, branch: string): Promise<RunValue[]> {
  const query = {
    db: "jbr",
    table: "report",
    fields: [{ n: "project" }, { n: "t", sql: "toUnixTimestamp(generated_time)" }, { n: "measures", subName: "name" }, { n: "measures", subName: "value" }],
    filters: [
      { f: "project", v: benchmark + "_%", o: "like" },
      { f: "branch", v: branch },
      { f: "generated_time", q: `>subtractDays(now(),${days})` },
      { f: "measures.name", q: "not like '%#%'" },
    ],
    order: "t",
  }
  const response = await fetch(serverConfigurator.computeSerializedQueryUrl(JSON.stringify(query)))
  if (!response.ok) throw new Error(`${benchmark}: ${response.status} ${await response.text()}`)
  // Without matching rows, the response holds no columns at all.
  const columns = (await response.json()) as [string[], number[], string[], number[]] | []
  if (columns.length === 0) return []
  const [projects, times, measures, measureValues] = columns
  return projects.map((project, index) => ({ project, time: times[index], measure: measures[index], value: measureValues[index] }))
}

let generation = 0
async function load(days: number, branch: string) {
  const current = ++generation
  loading.value = true
  error.value = null
  const benchmarks = [...new Set(rows.map((row) => row.benchmark))]
  const results = await Promise.allSettled(benchmarks.map((benchmark) => loadValues(benchmark, days, branch)))
  if (current !== generation) return
  const loaded = new Map<string, RunValue[]>()
  const failures: string[] = []
  for (const [index, result] of results.entries()) {
    if (result.status === "fulfilled") loaded.set(benchmarks[index], result.value)
    else failures.push(String(result.reason))
  }
  values.value = loaded
  error.value = failures.length > 0 ? "Not loaded: " + failures.join("; ") : null
  loading.value = false
}

// The compressor of the query URLs is ready once the server configurator emits.
serverConfigurator.createObservable().subscribe(() => void load(lookbackDays.value, branch.value))
watch([lookbackDays, branch], ([days, selectedBranch]) => void load(days, selectedBranch))

const tableRows = computed(() => rows.map((definition) => ({ definition, cells: summarizeRow(definition, values.value.get(definition.benchmark) ?? []) })))

const columns = computed(() => {
  const configurations = new Set(tableRows.value.flatMap((row) => [...row.cells.keys()]))
  return [...configurations].filter((configuration) => selectedPlatforms.value.includes(platformOf("_" + configuration))).toSorted(compareConfigurationNames)
})

const columnGroups = computed(() => {
  const groups: { platform: string; configurations: string[] }[] = []
  for (const configuration of columns.value) {
    const platform = platformOf("_" + configuration)
    const last = groups.at(-1)
    if (last?.platform === platform) last.configurations.push(configuration)
    else groups.push({ platform, configurations: [configuration] })
  }
  return groups
})

// The table body: each row's cell (or none) for every shown configuration.
const grid = computed(() =>
  tableRows.value.map((row) => ({ ...row, entries: columns.value.map((configuration) => ({ configuration, cell: row.cells.get(configuration) ?? null })) }))
)

const selection = ref<{ label: string; configuration: string } | null>(null)

function select(row: SummaryRow, configuration: string) {
  selection.value = isSelected(row, configuration) ? null : { label: row.label, configuration }
}

function isSelected(row: SummaryRow, configuration: string): boolean {
  return selection.value?.label === row.label && selection.value.configuration === configuration
}

const selectedCell = computed<{ row: SummaryRow; configuration: string; cell: SummaryCell } | null>(() => {
  const current = selection.value
  if (current == null) return null
  const row = tableRows.value.find((it) => it.definition.label === current.label)
  const cell = row?.cells.get(current.configuration)
  return row == null || cell == null ? null : { row: row.definition, configuration: current.configuration, cell }
})

function classOf(comparison: Comparison | null): string {
  if (comparison == null || !comparison.significant) return "noise"
  return comparison.change > 0 ? "better" : "worse"
}

function arrowOf(comparison: Comparison): string {
  if (!comparison.significant) return "≈"
  return comparison.change > 0 ? "▲" : "▼"
}

function formatChange(change: number): string {
  const percent = change * 100
  return (percent > 0 ? "+" : percent < 0 ? "−" : "") + Math.abs(percent).toFixed(1) + " %"
}

function formatInterval(comparison: Comparison): string {
  return comparison.low == null || comparison.high == null ? "too few runs" : `${formatChange(comparison.low)} … ${formatChange(comparison.high)}`
}

/** The confidence level of bounds at the misrate, in percent: 95 at 0.05, 99.7 at 0.003. */
function levelOf(misrate: number): string {
  return ((1 - misrate) * 100).toLocaleString("en-US", { maximumSignificantDigits: 3 })
}

const dateFormat = new Intl.DateTimeFormat("en-US", { month: "short", day: "numeric" })

function formatDate(time: number): string {
  return dateFormat.format(new Date(time * 1000))
}

function describe(cell: SummaryCell): string {
  const latest = cell.latest
  const judged =
    latest == null
      ? `no significant change in ${cell.runs - (cell.pending?.runsAfter ?? 0)} runs since ${formatDate(cell.firstRun)}`
      : `${formatChange(latest.change)} on ${formatDate(latest.time)} (${levelOf(ALPHA)} % interval ${formatInterval(latest)}), ` +
        `${latest.runsBefore} runs before, ${latest.runsAfter} after`
  const pending = cell.pending
  if (pending == null) return judged
  return `${judged}; a possible change on ${formatDate(pending.time)} (${formatChange(pending.change)} so far) awaits more runs, ${pending.runsAfter} of ${MIN_RUNS_AFTER}`
}

// Lets a configuration name wrap after its architecture and at its underscores, not inside "x86_64".
function withBreaks(configuration: string): string {
  return configuration.replace(/(x86_64|aarch64|x64)/, "$1​").replaceAll(/_(?!64)/g, "_​")
}
</script>

<style scoped>
.summary th,
.summary td {
  padding: 0.25rem 0.4rem;
  border-bottom: 1px solid var(--p-content-border-color, #e5e7eb);
  text-align: center;
}

.summary th.platform {
  font-weight: 600;
}

.summary th.configuration {
  font-weight: 400;
  font-size: 0.75rem;
  max-width: 7rem;
}

.summary th.benchmark {
  text-align: left;
  white-space: nowrap;
  font-weight: 500;
}

.cell {
  min-width: 5.5rem;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  border: 2px solid transparent;
  font-variant-numeric: tabular-nums;
  cursor: pointer;
}

.cell.selected {
  border-color: currentColor;
}

.date {
  display: block;
  font-size: 0.7rem;
  opacity: 0.8;
}

.worst {
  display: block;
  font-size: 0.7rem;
  color: #991b1b;
}

.cell .pending {
  display: block;
  font-size: 0.7rem;
  color: #6b7280;
}

.point.pending {
  color: #6b7280;
  border: 1px dashed currentColor;
}

.dark-mode .cell .pending,
.dark-mode .point.pending {
  color: #9ca3af;
}

.point {
  margin-left: 0.5rem;
  padding: 0 0.3rem;
  border-radius: 4px;
  font-variant-numeric: tabular-nums;
}

.better {
  color: #166534;
  background-color: #dcfce7;
}

.worse {
  color: #991b1b;
  background-color: #fee2e2;
}

.noise {
  color: #6b7280;
}

.dark-mode .better {
  color: #bbf7d0;
  background-color: #14532d;
}

.dark-mode .worse {
  color: #fecaca;
  background-color: #7f1d1d;
}

.dark-mode .noise {
  color: #9ca3af;
}

.dark-mode .worst {
  color: #fecaca;
}

.scenarios th,
.scenarios td {
  padding: 0.2rem 0.6rem;
  border-bottom: 1px solid var(--p-content-border-color, #e5e7eb);
  text-align: left;
}

.scenarios .number {
  text-align: right;
  font-variant-numeric: tabular-nums;
}
</style>
