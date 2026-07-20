<template>
  <div class="flex flex-col gap-1">
    <div class="text-sm font-medium">{{ title }}</div>
    <div
      v-if="!hasData"
      class="text-xs text-gray-500 py-2"
    >
      No comparable scenarios.
    </div>
    <div
      v-show="hasData"
      ref="chartElement"
      :style="{ height: `${chartHeight}px` }"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, useTemplateRef, watch } from "vue"
import { BarChart } from "echarts/charts"
import { GridComponent, TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import type { CallbackDataParams } from "echarts/types/dist/shared"
import { ChartManagerHelper } from "../../common/ChartManagerHelper"
import { BarChartOptions } from "../../common/echarts"
import { useDarkModeStore } from "../../../shared/useDarkModeStore"
import { EngineCompareRow, formatEngineCell } from "./useEngineComparison"
import { bucketLabel } from "./highlightingMetrics"
import { formatSignedPercent } from "../../charts/compareStats"

use([BarChart, GridComponent, TooltipComponent, CanvasRenderer])

// Matches the diverging tints in compareCells.css: red for slower (Δ% > 0), green for faster (Δ% < 0).
const SLOWER_COLOR = "#dc2626"
const FASTER_COLOR = "#16a34a"

const { rows, title } = defineProps<{ rows: EngineCompareRow[]; title: string }>()

// Ascending by Δ% so the most negative (best improvement) sits at the bottom and the largest positive
// (worst regression) at the top — ECharts places category index 0 at the bottom of a horizontal bar.
const sortedRows = computed<EngineCompareRow[]>(() => rows.filter((row) => Number.isFinite(row.diffPercent)).sort((a, b) => a.diffPercent - b.diffPercent))

const hasData = computed(() => sortedRows.value.length > 0)
const chartHeight = computed(() => Math.max(120, sortedRows.value.length * 24 + 40))

const chartElement = useTemplateRef<HTMLElement>("chartElement")
let chartManager: ChartManagerHelper | null = null

function tooltipFor(row: EngineCompareRow): string {
  const { before, after } = formatEngineCell(row)
  return `${row.title} · ${bucketLabel(row.bucket)}<br/>LEGACY ${before} → NEW ${after}<br/>Δ ${formatSignedPercent(row.diffPercent)}`
}

function renderChart(): void {
  if (chartManager == null || !hasData.value) return
  const dark = useDarkModeStore().darkMode
  const textColor = dark ? "#cccccc" : "#333333"
  const lineColor = dark ? "#3a3a3a" : "#e5e7eb"
  const data = sortedRows.value

  chartManager.chart.setOption<BarChartOptions>(
    {
      animation: false,
      textStyle: { color: textColor },
      grid: { top: 6, left: 5, right: 48, bottom: 22, containLabel: true },
      tooltip: {
        trigger: "item",
        formatter: (params) => tooltipFor(data[(params as CallbackDataParams).dataIndex]),
      },
      xAxis: {
        type: "value",
        axisLabel: { color: textColor, formatter: (value: number) => `${value}%` },
        splitLine: { lineStyle: { color: lineColor } },
      },
      yAxis: {
        type: "category",
        data: data.map((row) => `${row.title} · ${bucketLabel(row.bucket)}`),
        axisLabel: { color: textColor },
        axisLine: { lineStyle: { color: lineColor } },
      },
      series: [
        {
          type: "bar",
          barMaxWidth: 16,
          data: data.map((row) => ({
            value: row.diffPercent,
            itemStyle: { color: row.diffPercent > 0 ? SLOWER_COLOR : FASTER_COLOR },
            label: { position: row.diffPercent >= 0 ? "right" : "left" },
          })),
          label: {
            show: true,
            color: textColor,
            formatter: (params: CallbackDataParams) => formatSignedPercent(params.value as number),
          },
        },
      ],
    },
    { replaceMerge: ["series", "yAxis"] }
  )
  chartManager.chart.resize()
}

onMounted(() => {
  if (chartElement.value == null) return
  chartManager = new ChartManagerHelper(chartElement.value)
  renderChart()
})

// Re-render on data or theme change; the height binding updates first, so resize after the DOM settles.
watch([sortedRows, () => useDarkModeStore().darkMode], () => {
  void nextTick().then(renderChart)
})

onUnmounted(() => {
  chartManager?.dispose()
  chartManager = null
})
</script>
