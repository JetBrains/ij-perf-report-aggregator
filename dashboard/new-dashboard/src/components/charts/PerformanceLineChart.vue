<template>
  <div class="flex flex-col gap-y-2.5 py-3 px-5 border border-solid rounded-md border-zinc-200">
    <h3 class="m-0 text-sm">
      {{ props.title }}
    </h3>
    <div
      ref="chartElement"
      class="bg-white"
      :style="{ height: `${chartHeight}px` }"
    />
  </div>
</template>
<script setup lang="ts">
import { useElementVisibility } from "@vueuse/core"
import { CallbackDataParams } from "echarts/types/src/util/types"
import { computed, inject, onMounted, onUnmounted, Ref, shallowRef, toRef, watch } from "vue"
import { PredefinedMeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { injectOrError, reportInfoProviderKey } from "../../shared/injectionKeys"
import { accidentsConfiguratorKey, containerKey, sidebarVmKey } from "../../shared/keys"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "../common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { getInfoDataFrom } from "../common/sideBar/InfoSidebarPerformance"
import { SeriesNameConfigurator } from "../startup/SeriesNameConfigurator"
import { PerformanceChartManager } from "./PerformanceChartManager"
import { PerformanceLineChartVM } from "./PerformanceLineChartVM"

interface LineChartProps {
  title: string
  measures: string[]
  configurators: (DataQueryConfigurator | FilterConfigurator)[]
  skipZeroValues?: boolean
  chartType?: ChartType
  valueUnit?: ValueUnit
  legendFormatter?: (name: string) => string
  withMeasureName?: boolean
}

const props = withDefaults(defineProps<LineChartProps>(), {
  skipZeroValues: true,
  valueUnit: "ms",
  chartType: "line",
  legendFormatter(name: string): string {
    return name
  },
  withMeasureName: false,
})

const accidentsConfigurator = inject(accidentsConfiguratorKey, null)
const chartElement = shallowRef<HTMLElement>()

const chartIsVisible = useElementVisibility(chartElement)

const skipZeroValues = toRef(props, "skipZeroValues")
const reportInfoProvider = inject(reportInfoProviderKey, null)

const infoFieldsConfigurator =
  reportInfoProvider && reportInfoProvider.infoFields.length > 0
    ? {
        createObservable() {
          return null
        },
        configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
          for (const infoField of reportInfoProvider.infoFields) {
            query.addField(infoField)
          }
          return true
        },
      }
    : null

const container = injectOrError(containerKey)
const sidebarVm = injectOrError(sidebarVmKey)

let chartManager: PerformanceChartManager | null
let chartVm: PerformanceLineChartVM | null = null
let unsubscribe: (() => void) | null = null

const measures: Ref<string[]> = computed(() => {
  return props.measures
})

const measureConfigurator = new PredefinedMeasureConfigurator(
  measures,
  skipZeroValues,
  props.chartType,
  props.valueUnit,
  {
    symbolSize: 7,
    showSymbol: false,
  },
  accidentsConfigurator
)

if (measures.value.length == 1) {
  new SeriesNameConfigurator(measures.value[0])
}

const configurators = [...props.configurators, measureConfigurator, infoFieldsConfigurator]
if (props.withMeasureName) {
  configurators.push(new SeriesNameConfigurator(measures.value[0]))
}

const dataQueryExecutor = new DataQueryExecutor([...configurators].filter((item): item is DataQueryConfigurator => item != null))

function createChart() {
  if (chartVm != null) {
    return
  }
  if (chartElement.value) {
    chartManager?.dispose()
    unsubscribe?.()
    chartManager = new PerformanceChartManager(chartElement.value, container.value)
    chartVm = new PerformanceLineChartVM(chartManager, dataQueryExecutor, props.valueUnit, accidentsConfigurator, props.legendFormatter)
    unsubscribe = chartVm.subscribe()
    chartManager.chart.on("click", (params: CallbackDataParams) => {
      const infoData = getInfoDataFrom(sidebarVm.type, params, props.valueUnit, accidentsConfigurator)
      sidebarVm.show(infoData)
    })
  } else {
    console.error("Dom was not yet initialized")
  }
}

function setupChartOnVisibility() {
  watch(
    chartIsVisible,
    (isVisible) => {
      if (isVisible) {
        createChart()
      }
    },
    { immediate: true }
  )
  // If the chart is not visible, still try to create it after a delay
  setTimeout(createChart, 1000)
}

onMounted(() => {
  setupChartOnVisibility()
})

onUnmounted(() => {
  unsubscribe?.()
  chartManager?.dispose()
})

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
</script>
