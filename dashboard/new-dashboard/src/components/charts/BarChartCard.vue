<template>
  <div
    ref="chartElement"
    class="bg-white overflow-hidden shadow rounded-lg w-full"
    :style="{ height: `${height}px` }"
  />
</template>
<script setup lang="ts">
import { inject, onMounted, onUnmounted, useTemplateRef } from "vue"
import { PredefinedGroupingMeasureConfigurator } from "../../configurators/PredefinedGroupingMeasureConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey, configuratorListKey, injectOrError, timeRangeKey } from "../../shared/injectionKeys"
import { BarChartManager } from "../common/BarChartManager"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { chartDefaultStyle } from "../common/chart"

const { height = 440, measures = [] } = defineProps<{
  height?: number
  measures: string[]
}>()

const chartElement = useTemplateRef<HTMLElement>("chartElement")
let chartManager: BarChartManager | null = null

const timeRange = injectOrError(timeRangeKey)
const aggregationOperatorConfigurator = injectOrError(aggregationOperatorConfiguratorKey)

const chartStyle = inject(chartStyleKey, chartDefaultStyle)
const measureConfigurator = new PredefinedGroupingMeasureConfigurator(measures, timeRange, chartStyle)
const dataQueryExecutor = new DataQueryExecutor([...injectOrError(configuratorListKey), aggregationOperatorConfigurator, measureConfigurator])
onMounted(() => {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  chartManager = new BarChartManager(chartElement.value!, dataQueryExecutor, chartStyle)
})
onUnmounted(() => {
  const it = chartManager
  if (it != null) {
    chartManager = null
    it.dispose()
  }
})
</script>
