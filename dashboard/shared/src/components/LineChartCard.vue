<template>
  <div
    ref="chartElement"
    class="bg-white shadow rounded-lg"
    :style="{height: `${chartHeight}px`}"
    @mouseenter="show"
    @mouseleave="hide"
  />
</template>
<script setup lang="ts">
import { inject, onMounted, onUnmounted, PropType, shallowRef, toRef, watchEffect } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { DEFAULT_LINE_CHART_HEIGHT } from "../chart"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import { DataQuery, DataQueryExecutorConfiguration } from "../dataQuery"
import { chartToolTipKey, configuratorListKey } from "../injectionKeys"
import { ChartToolTipManager } from "./ChartToolTipManager"
import { LineChartManager } from "./LineChartManager"

const props = defineProps({
  skipZeroValues: {
    type: Boolean,
    default: true,
  },
  dataZoom: {
    type: Boolean,
    default: false,
  },
  measures: {
    type: Array as PropType<Array<string> | null>,
    default: () => null,
  },
})

const chartElement = shallowRef<HTMLElement | null>(null)
let chartManager: LineChartManager | null = null
const providedConfigurators = inject(configuratorListKey, null)
const skipZeroValues = toRef(props, "skipZeroValues")
const chartToolTipManager = new ChartToolTipManager()
// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
const tooltip = inject(chartToolTipKey)!

const show = (event: Event) => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  tooltip.value?.["show"](event, chartToolTipManager)
}
const hide = () => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  tooltip.value?.["hide"]()
}

let dataQueryExecutor: DataQueryExecutor | null

watchEffect(function () {
  let configurators = providedConfigurators
  if (configurators == null) {
    throw new Error("`dataQueryExecutor` is not provided")
  }

  // static list of measures is provided - create sub data query executor
  if (props.measures != null) {
    configurators = configurators.concat(new PredefinedMeasureConfigurator(props.measures, skipZeroValues))
    const infoFields = chartToolTipManager.reportInfoProvider?.infoFields ?? []
    if (infoFields.length !== 0) {
      configurators.push({
        createObservable() {
          return null
        },
        configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
          for (const infoField of infoFields) {
            query.addField(infoField)
          }
          return true
        }
      })
    }
  }
  dataQueryExecutor = new DataQueryExecutor(configurators)
  chartToolTipManager.dataQueryExecutor = dataQueryExecutor
  if (chartManager != null) {
    chartManager.dataQueryExecutor = dataQueryExecutor
  }
})

onMounted(() => {
  chartManager = new LineChartManager(
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    chartElement.value!,
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    dataQueryExecutor!,
    toRef(props, "dataZoom"),
    chartToolTipManager.formatArrayValue.bind(chartToolTipManager),
  )
})
onUnmounted(() => {
  const it = chartManager
  if (it != null) {
    chartManager = null
    it.dispose()
  }
})

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
</script>
<style scoped>

a {
  text-decoration: none;
}

</style>