<template>
  <div
    ref="chartElement"
    class="bg-white shadow rounded-lg"
    :style="{height: `${chartHeight}px`}"
    @mouseenter="show"
    @mouseleave="hide"
  />
  <ChartTooltip
    ref="tooltip"
  />
</template>
<script setup lang="ts">
import { inject, onMounted, onUnmounted, PropType, ref, shallowRef, toRef, watch, watchEffect } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { DEFAULT_LINE_CHART_HEIGHT } from "../chart"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { dataQueryExecutorKey } from "../injectionKeys"
import { ChartToolTipManager } from "./ChartToolTipManager"
import ChartTooltip from "./ChartTooltip.vue"
import { LineChartManager } from "./LineChartManager"

const props = defineProps({
  provider: {
    type: DataQueryExecutor,
    default: () => null,
  },
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
const providedDataQueryExecutor = inject(dataQueryExecutorKey, null)
const skipZeroValues = toRef(props, "skipZeroValues")
const chartToolTipManager = new ChartToolTipManager()
const tooltip = ref<typeof ChartTooltip>()

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
  dataQueryExecutor = props.provider ?? providedDataQueryExecutor
  if (dataQueryExecutor == null) {
    throw new Error("Neither `provider` property is set, nor `dataQueryExecutor` is provided")
  }

  // static list of measures is provided - create sub data query executor
  if (props.measures != null) {
    const configurators: Array<DataQueryConfigurator> = [
      new PredefinedMeasureConfigurator(props.measures, skipZeroValues),
    ]
    const infoFields = chartToolTipManager.reportInfoProvider?.infoFields ?? []
    if (infoFields.length !== 0) {
      configurators.push({
        configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
          for (const infoField of infoFields) {
            query.addField(infoField)
          }
          return true
        },
      })
    }
    dataQueryExecutor = dataQueryExecutor.createSub(configurators)
    dataQueryExecutor.scheduleLoad()
  }

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

  watch(skipZeroValues, () => {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    chartManager!.dataQueryExecutor.scheduleLoad()
  })
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