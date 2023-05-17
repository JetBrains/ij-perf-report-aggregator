<template>
  <div class="flex flex-col gap-y-2.5 py-3 px-5 border border-solid rounded-md border-zinc-200">
    <h3 class="m-0 text-sm">
      {{ props.title }}
    </h3>
    <div
      ref="chartElement"
      class="bg-white"
      :style="{height: `${chartHeight}px`}"
    />
  </div>
</template>
<script setup lang="ts">
import { CallbackDataParams } from "echarts/types/src/util/types"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "shared/src/chart"
import { PredefinedMeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration} from "shared/src/dataQuery"
import { reportInfoProviderKey } from "shared/src/injectionKeys"
import { Accident } from "shared/src/meta"
import { calculateChanges } from "shared/src/util/changes"
import { inject, onMounted, onUnmounted, shallowRef, toRef, withDefaults } from "vue"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import { getInfoDataFrom, InfoData, InfoSidebarVm } from "../InfoSidebarVm"
import { ChartManager } from "./ChartManager"
import { LineChartVM } from "./LineChartVM"

interface LineChartProps {
  title: string
  measures: Array<string>
  configurators: Array<DataQueryConfigurator>
  skipZeroValues?: boolean
  chartType?: ChartType
  valueUnit?: ValueUnit
  accidents?: Array<Accident>|null
}

const props = withDefaults(defineProps<LineChartProps>(), {
  skipZeroValues: true,
  valueUnit: "ms",
  chartType: "line",
  accidents: null
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

const infoFieldsConfigurator = reportInfoProvider && reportInfoProvider.infoFields.length > 0 ?
  {
    createObservable() {
      return null
    },
    configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
      for (const infoField of reportInfoProvider.infoFields) {
        query.addField(infoField)
      }
      return true
    },
  } : null
const dataQueryExecutor = new DataQueryExecutor([
  ...props.configurators,
  measureConfigurator,
  infoFieldsConfigurator,
].filter((item): item is DataQueryConfigurator => item != null))

const container = inject(containerKey)
const sidebarVm = inject(sidebarVmKey)

let chartManager: ChartManager
let chartVm: LineChartVM
let unsubscribe: (() => void) | null = null

onMounted(() => {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  chartManager = new ChartManager(chartElement.value!, container?.value)
  chartVm = new LineChartVM(
    chartManager,
    dataQueryExecutor,
    props.valueUnit,
    props.accidents
  )

  unsubscribe = chartVm.subscribe()

  chartManager.chart.on("click", (params: CallbackDataParams) => {
    if (params.dataIndex != undefined) {
      const infoData = getInfoDataFrom(params, props.valueUnit, props.accidents)
      showSideBar(sidebarVm, infoData)
    }
  })
})

function showSideBar(sidebarVm: InfoSidebarVm | undefined, infoData: InfoData) {
  const separator = ".."
  const db = infoData.installerId ? "perfint" : "perfintDev"
  const id = infoData.installerId ?? infoData.buildId
  calculateChanges(db, id, (decodedChanges: string|null) => {
    if(decodedChanges != null) {
      infoData.changes = decodedChanges
    }
    sidebarVm?.show(infoData)
  })
}

onUnmounted(() => {
  if(unsubscribe != null) unsubscribe()
  chartManager.dispose()
})

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
</script>