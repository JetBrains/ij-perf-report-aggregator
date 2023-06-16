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
import { isDefined } from "@vueuse/core"
import { CallbackDataParams } from "echarts/types/src/util/types"
import { inject, onMounted, onUnmounted, shallowRef, toRef, watch } from "vue"
import { PredefinedMeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { reportInfoProviderKey } from "../../shared/injectionKeys"
import { accidentsKeys, containerKey, sidebarVmKey } from "../../shared/keys"
import { calculateChanges } from "../../util/changes"
import { getInfoDataFrom, InfoData, InfoSidebarVm } from "../InfoSidebarVm"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "../common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { ChartManager } from "./ChartManager"
import { LineChartVM } from "./LineChartVM"

interface LineChartProps {
  title: string
  measures: string[]
  configurators: DataQueryConfigurator[]
  skipZeroValues?: boolean
  chartType?: ChartType
  valueUnit?: ValueUnit
}

const props = withDefaults(defineProps<LineChartProps>(), {
  skipZeroValues: true,
  valueUnit: "ms",
  chartType: "line",
})

const accidents = inject(accidentsKeys)
const chartElement = shallowRef<HTMLElement>()
const skipZeroValues = toRef(props, "skipZeroValues")
const reportInfoProvider = inject(reportInfoProviderKey, null)
const measureConfigurator = new PredefinedMeasureConfigurator(
  props.measures,
  skipZeroValues,
  props.chartType,
  props.valueUnit,
  {
    symbolSize: 7,
    showSymbol: false,
  },
  accidents
)

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
const dataQueryExecutor = new DataQueryExecutor([...props.configurators, measureConfigurator, infoFieldsConfigurator].filter((item): item is DataQueryConfigurator => item != null))

const container = inject(containerKey)
const sidebarVm = inject(sidebarVmKey)

let chartManager: ChartManager | null
let chartVm: LineChartVM
let unsubscribe: (() => void) | null = null

function initializePlot() {
  if (chartElement.value) {
    chartManager?.dispose()
    unsubscribe?.()
    chartManager = new ChartManager(chartElement.value, container?.value)
    chartVm = new LineChartVM(chartManager, dataQueryExecutor, props.valueUnit, accidents)
    unsubscribe = chartVm.subscribe()
    chartManager.chart.on("click", (params: CallbackDataParams) => {
      const infoData = getInfoDataFrom(params, props.valueUnit, accidents)
      sidebarVm?.show(infoData)
    })
  } else {
    console.error("Dom was not yet initialized")
  }
}

onMounted(() => {
  if (isDefined(accidents)) {
    watch(accidents, () => {
      if (isDefined(accidents)) {
        initializePlot()
      }
    })
  } else {
    initializePlot()
  }
})

onUnmounted(() => {
  unsubscribe?.()
  chartManager?.dispose()
})

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
</script>
