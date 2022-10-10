<template>
  <div
    ref="chartElement"
    class="bg-white"
    :style="{height: `${chartHeight}px`}"
  />
</template>
<script setup lang="ts">
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "shared/src/chart"
import { PredefinedMeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { DataQueryConfigurator } from "shared/src/dataQuery"
import { onMounted, onUnmounted, shallowRef, toRef, withDefaults } from "vue"
import { LineChartVM } from "./LineChartVM"

interface LineChartProps {
  measures: Array<string>
  configurators: Array<DataQueryConfigurator>
  skipZeroValues?: boolean
  chartType?: ChartType
  valueUnit?: ValueUnit
}

const props = withDefaults(defineProps<LineChartProps>(), {
  skipZeroValues: true,
  valueUnit: "ms",
})

const chartElement = shallowRef<HTMLElement>()
const skipZeroValues = toRef(props, "skipZeroValues")
const measureConfigurator = new PredefinedMeasureConfigurator(
  props.measures,
  skipZeroValues,
  props.chartType,
  props.valueUnit,
  {
    symbolSize: 7,
    showSymbol: false,
  },
)

const dataQueryExecutor = new DataQueryExecutor( [
  ...props.configurators,
  measureConfigurator,
])


let chart: ChartManagerHelper
let chartVm: LineChartVM

onMounted(() => {
  chart = new ChartManagerHelper(chartElement.value!)
  chartVm = new LineChartVM(
    chart,
    dataQueryExecutor,
    props.valueUnit,
  )

  chartVm.subscribe()
})

onUnmounted(() => {
  // TODO: Make them lifetimed for auto-dispose
  chart.dispose()
  chartVm.dispose()
})

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
</script>