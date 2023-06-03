<template>
  <div
    v-if="label == null"
    ref="chartElement"
    class="bg-white shadow rounded-lg"
    :style="{height: `${chartHeight}px`}"
    @mouseenter="show"
    @mouseleave="hide"
  />
  <div
    v-else
    class="bg-white shadow rounded-lg"
    @mouseenter="show"
    @mouseleave="hide"
  >
    <div class="flex justify-center mt-3 mb-2">
      <h3
        class="px-2 bg-white text-sm text-gray-900 capitalize"
      >
        {{ label }}
      </h3>
    </div>
    <div
      ref="chartElement"
      :style="{height: `${chartHeight}px`}"
    />
  </div>
</template>
<script setup lang="ts">
import { inject, onMounted, onUnmounted, shallowRef, toRef, watchEffect } from "vue"
import { PredefinedMeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { chartToolTipKey, configuratorListKey } from "../../shared/injectionKeys"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "../common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { ChartToolTipManager } from "./ChartToolTipManager"
import { LineChartManager, PopupTrigger } from "./LineChartManager"

const props = withDefaults(defineProps<{
  skipZeroValues?: boolean
  compoundTooltip?: boolean
  dataZoom?: boolean
  measures?: string[] | null
  chartType?: ChartType
  valueUnit?: ValueUnit
  configurators?: DataQueryConfigurator[]|null
  trigger?: PopupTrigger
  label?: string
  aggregatedMeasure?: string | null
}>(), {
  skipZeroValues: true,
  compoundTooltip: true,
  dataZoom: false,
  measures: null,
  chartType: "line",
  valueUnit: "ms",
  configurators: null,
  trigger: "axis",
  aggregatedMeasure: null,
  label: undefined,
})

const chartElement = shallowRef<HTMLElement | null>(null)
let chartManager: LineChartManager | null = null

const skipZeroValues = toRef(props, "skipZeroValues")
const chartToolTipManager = new ChartToolTipManager(props.valueUnit)
// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
const tooltip = inject(chartToolTipKey)!

const show = () => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  tooltip.value["show"](chartToolTipManager)
}
const hide = () => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  tooltip.value["hide"]()
}

let dataQueryExecutor: DataQueryExecutor | null

const providedConfigurators = inject(configuratorListKey)
watchEffect(function () {
  let configurators = props.configurators ?? providedConfigurators
  if (configurators == null) {
    throw new Error(`${configurators} is not provided`)
  }

  // static list of measures is provided - create sub data query executor
  if (props.measures != null) {
    configurators = [...configurators, new PredefinedMeasureConfigurator(props.measures, skipZeroValues, props.chartType, props.valueUnit)]
    const infoFields = chartToolTipManager.reportInfoProvider?.infoFields ?? []
    if (infoFields.length > 0) {
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

  if (props.aggregatedMeasure != null) {
    configurators = [...configurators]
    configurators.push({
      configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
        if (props.aggregatedMeasure != null) {
          query.addFilter({f: "measures.name", v: props.aggregatedMeasure})
        }
        return true
      }, createObservable() {
        return null
      },
    })
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
    props.compoundTooltip ? chartToolTipManager : null,
    props.valueUnit,
    props.trigger,
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