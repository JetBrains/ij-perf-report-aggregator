<template>
  <div class="flex flex-col gap-2 py-3 px-5 border border-solid rounded-md overflow-x-auto">
    <div class="text-sm font-semibold">Δ% by scenario <span class="font-normal text-gray-500">(NEW vs LEGACY; red = slower, green = faster; hover for exact values)</span></div>
    <table class="border-collapse text-sm">
      <thead>
        <tr>
          <th
            rowspan="2"
            class="text-left font-medium px-2 py-1 border-b border-gray-200 dark:border-gray-700"
          >
            Project
          </th>
          <th
            v-for="group in headerGroups"
            :key="group.phase"
            :colspan="group.types.length"
            class="text-center font-medium px-2 py-1 border-b border-l border-gray-200 dark:border-gray-700"
          >
            {{ group.label }}
          </th>
        </tr>
        <tr>
          <th
            v-for="type in selectedMetricTypes"
            :key="type"
            class="text-center font-normal text-gray-500 px-2 py-1 border-b border-l border-gray-200 dark:border-gray-700"
          >
            {{ bucketLabel(bucketOf(type)) }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="row in gridRows"
          :key="row.base"
        >
          <td class="text-left font-medium px-2 py-1 whitespace-nowrap border-b border-gray-100 dark:border-gray-800">{{ row.title }}</td>
          <td
            v-for="(cell, index) in row.cells"
            :key="index"
            v-tooltip="cell.tooltip"
            class="text-center px-2 py-1 tabular-nums border-b border-l border-gray-100 dark:border-gray-800"
            :class="cell.class"
          >
            {{ cell.text }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { formatSignedPercent, severityCellClass } from "../../charts/compareStats"
import { EngineCompareRow, EngineComparisonMode, formatEngineCell } from "./useEngineComparison"
import { phaseLabel, phaseOf, projects, bucketOf, bucketLabel } from "./highlightingMetrics"
import "../../charts/compareCells.css"

const { rows, selectedMetricTypes, mode } = defineProps<{ rows: EngineCompareRow[]; selectedMetricTypes: string[]; mode: EngineComparisonMode }>()

// (base, metric type) -> row, so a cell lookup is O(1) across the 9 × up-to-9 grid.
const rowByCell = computed(() => {
  const map = new Map<string, EngineCompareRow>()
  for (const row of rows) map.set(`${row.base}::${row.metricType}`, row)
  return map
})

interface HeaderGroup {
  phase: string
  label: string
  types: string[]
}

// Groups the selected columns by phase for the two-level header. METRIC_TYPES orders all buckets of a
// phase contiguously, and the selection preserves that order, so consecutive columns share a phase.
const headerGroups = computed<HeaderGroup[]>(() => {
  const groups: HeaderGroup[] = []
  for (const type of selectedMetricTypes) {
    const phase = phaseOf(type)
    const last = groups.at(-1)
    if (last != null && last.phase === phase) last.types.push(type)
    else groups.push({ phase, label: phaseLabel(phase), types: [type] })
  }
  return groups
})

interface CellData {
  text: string
  class: string
  tooltip: string
}

function buildCell(base: string, type: string): CellData {
  const row = rowByCell.value.get(`${base}::${type}`)
  if (row == null || !Number.isFinite(row.diffPercent)) return { text: "—", class: "text-gray-400", tooltip: "" }
  const { before, after } = formatEngineCell(row)
  // Aggregated over many runs shows the sample sizes; a single run has n=1 and omits them.
  const counts = mode === "aggregated" ? ` · n=${row.legacy.count}/${row.new.count}` : ""
  return {
    text: formatSignedPercent(row.diffPercent, 0),
    class: severityCellClass(row.diffPercent) ?? "",
    tooltip: `LEGACY ${before} → NEW ${after} · Δ ${formatSignedPercent(row.diffPercent)}${counts}`,
  }
}

// One row per project (all 9 shown even when a base has no data), each cell aligned to selectedMetricTypes.
const gridRows = computed(() =>
  projects.map((project) => ({
    base: project.base,
    title: project.title,
    cells: selectedMetricTypes.map((type) => buildCell(project.base, type)),
  }))
)
</script>
