<template>
  <div class="flex flex-col gap-y-2.5 py-3 px-5 border border-solid rounded-md border-zinc-200">
    <h3
      v-if="props.title != null"
      class="m-0 text-sm"
    >
      {{ props.title }}
    </h3>
    <div
      ref="chartElement"
      class="bg-white"
      :style="{ height: `${chartHeight}px` }"
      @mouseenter="show"
      @mouseleave="hide"
    />
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, shallowRef, toRef, watchEffect } from "vue"
import { PredefinedMeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { chartToolTipKey, configuratorListKey, injectOrError } from "../../shared/injectionKeys"
import { containerKey, sidebarStartupKey } from "../../shared/keys"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "../common/chart"
import { DataQuery, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { SeriesNameConfigurator } from "../startup/SeriesNameConfigurator"
import { StartupLineChartManager } from "./StartupLineChartManager"
import { StartupTooltipManager } from "./StartupTooltipManager"

interface StartupLineChartProps {
  title?: string | undefined
  measures: string[]
  dataZoom?: boolean
  skipZeroValues?: boolean
  chartType?: ChartType
  valueUnit?: ValueUnit
}

const props = withDefaults(defineProps<StartupLineChartProps>(), {
  skipZeroValues: true,
  dataZoom: false,
  valueUnit: "ms",
  chartType: "line",
  title: undefined,
})

const chartElement = shallowRef<HTMLElement>()
let chartManager: StartupLineChartManager | null = null

const skipZeroValues = toRef(props, "skipZeroValues")
const chartToolTipManager = new StartupTooltipManager(props.valueUnit)
const container = injectOrError(containerKey)
const tooltip = injectOrError(chartToolTipKey)
const sidebarVm = injectOrError(sidebarStartupKey)

const show = () => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  tooltip.value["show"](chartToolTipManager)
}
const hide = () => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-call
  tooltip.value["hide"]()
}

let dataQueryExecutor: DataQueryExecutor | null

const providedConfigurators = injectOrError(configuratorListKey)
let unsubscribe: (() => void) | null = null

watchEffect(function () {
  const configurators = [...providedConfigurators, new PredefinedMeasureConfigurator(props.measures, skipZeroValues, props.chartType, props.valueUnit, {}, null)]
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
      },
    })
  }
  if (props.measures.length == 1) {
    configurators.push(new SeriesNameConfigurator(props.measures[0]))
  }
  dataQueryExecutor = new DataQueryExecutor(configurators)

  chartToolTipManager.dataQueryExecutor = dataQueryExecutor
})

onMounted(() => {
  if (chartElement.value) {
    chartManager = new StartupLineChartManager(
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      chartElement.value,
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      dataQueryExecutor!,
      toRef(props, "dataZoom"),
      chartToolTipManager,
      sidebarVm,
      container.value
    )
    unsubscribe = chartManager.subscribe()
  }
})

onUnmounted(() => {
  unsubscribe?.()
  chartManager?.dispose()
})

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
</script>
