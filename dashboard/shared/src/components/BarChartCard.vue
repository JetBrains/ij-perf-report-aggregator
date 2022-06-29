<template>
  <div
    ref="chartElement"
    class="bg-white overflow-hidden shadow rounded-lg w-full"
    :style="{height: `${height}px`}"
  />
</template>
<script setup lang="ts">
import { inject, onMounted, onUnmounted, shallowRef } from "vue"
import { BarChartManager } from "../BarChartManager"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { chartDefaultStyle } from "../chart"
import { PredefinedGroupingMeasureConfigurator } from "../configurators/PredefinedGroupingMeasureConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey, configuratorListKey, injectOrError } from "../injectionKeys"

const props = withDefaults(defineProps<{
  height?: number
  measures: Array<string>
}>(), {
  height: 440,
  valueUnit: "ms",
  measures: () => [],
})

const chartElement = shallowRef<HTMLElement | null>(null)
let chartManager: BarChartManager | null = null
// eslint-disable-next-line vue/no-setup-props-destructure
const measures = props.measures

const aggregationOperatorConfigurator = injectOrError(aggregationOperatorConfiguratorKey)

const chartStyle = inject(chartStyleKey, chartDefaultStyle)
const measureConfigurator = new PredefinedGroupingMeasureConfigurator(measures, chartStyle)
// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
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