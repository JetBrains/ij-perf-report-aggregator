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
import { CallbackDataParams } from "echarts/types/src/util/types"
import { inject, onMounted, onUnmounted, shallowRef, toRef } from "vue"
import { PredefinedMeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { reportInfoProviderKey } from "../../shared/injectionKeys"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import { calculateChanges } from "../../util/changes"
import { Accident } from "../../util/meta"
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
  accidents?: Accident[] | null
}

const props = withDefaults(defineProps<LineChartProps>(), {
  skipZeroValues: true,
  valueUnit: "ms",
  chartType: "line",
  accidents: null,
})

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
  props.accidents
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

let chartManager: ChartManager
let chartVm: LineChartVM
let unsubscribe: (() => void) | null = null

onMounted(() => {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  chartManager = new ChartManager(chartElement.value!, container?.value)
  chartVm = new LineChartVM(chartManager, dataQueryExecutor, props.valueUnit, props.accidents)

  unsubscribe = chartVm.subscribe()

  chartManager.chart.on("click", (params: CallbackDataParams) => {
    const infoData = getInfoDataFrom(params, props.valueUnit, props.accidents)
    showSideBar(sidebarVm, infoData)
  })
})

function showSideBar(sidebarVm: InfoSidebarVm | undefined, infoData: InfoData) {
  const db = infoData.installerId ? "perfint" : "perfintDev"
  const id = infoData.installerId ?? infoData.buildId
  calculateChanges(db, id, (decodedChanges: string | null) => {
    if (decodedChanges != null) {
      infoData.changes = decodedChanges
    }
    sidebarVm?.show(infoData)
  })
}

onUnmounted(() => {
  if (unsubscribe != null) unsubscribe()
  chartManager.dispose()
})

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
</script>
