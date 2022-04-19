<template>
  <div
    ref="chartElement"
    class="bg-white overflow-hidden shadow rounded-lg w-full"
    :style="{height: `${height}px`}"
  />
</template>
<script setup lang="ts">
import { inject, onMounted, onUnmounted, Ref, shallowRef } from "vue"
import { BarChartManager } from "../BarChartManager"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { chartDefaultStyle } from "../chart"
import { PredefinedGroupingMeasureConfigurator } from "../configurators/PredefinedGroupingMeasureConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey, configuratorListKey, timeRangeKey } from "../injectionKeys"

const props = withDefaults(defineProps<{
  height?: number
  measures: Array<string>
}>(), {
  height: 440,
  measures: () => [],
})

const chartElement = shallowRef<HTMLElement | null>(null)
let chartManager: BarChartManager | null = null
// eslint-disable-next-line vue/no-setup-props-destructure
const measures = props.measures

const timeRange = inject(timeRangeKey)
if (timeRange === undefined) {
  throw new Error("timeRange is not injected but required")
}

const aggregationOperatorConfigurator = inject(aggregationOperatorConfiguratorKey)
if (aggregationOperatorConfigurator === undefined) {
  throw new Error("aggregationOperatorConfigurator is not injected but required")
}

const measureConfigurator = new PredefinedGroupingMeasureConfigurator(measures, timeRange, inject(chartStyleKey, chartDefaultStyle))
// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
const dataQueryExecutor = new DataQueryExecutor(inject(configuratorListKey)!.concat(aggregationOperatorConfigurator, measureConfigurator))
onMounted(() => {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  chartManager = new BarChartManager(chartElement.value!, dataQueryExecutor)
})
onUnmounted(() => {
  const it = chartManager
  if (it != null) {
    chartManager = null
    it.dispose()
  }
})
</script>