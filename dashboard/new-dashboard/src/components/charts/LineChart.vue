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
import { CallbackDataParams } from "echarts/types/dist/shared"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "shared/src/chart"
import { PredefinedMeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, QueryProducer, SimpleQueryProducer } from "shared/src/dataQuery"
import { reportInfoProviderKey } from "shared/src/injectionKeys"
import { inject, onMounted, onUnmounted, shallowRef, toRef, withDefaults } from "vue"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import { getInfoDataFrom, InfoData } from "../InfoSidebarVm"
import { ChartManager } from "./ChartManager"
import { LineChartVM } from "./LineChartVM"
import { filter, Observable, shareReplay } from "rxjs"
import { FilterConfigurator } from "shared/src/configurators/filter"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"

interface LineChartProps {
  title: string
  measures: Array<string>
  configurators: Array<DataQueryConfigurator>
  skipZeroValues?: boolean
  chartType?: ChartType
  valueUnit?: ValueUnit
}

const props = withDefaults(defineProps<LineChartProps>(), {
  skipZeroValues: true,
  valueUnit: "ms",
  chartType: "line",
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
let unsubscribe: (() => void)|null  = null

function fetchChangesFromInstaller(infoData: InfoData) {
  if(infoData.installerId == undefined) {
    sidebarVm?.show(infoData)
  }
  const serverUrlObservable = refToObservable(shallowRef(ServerConfigurator.DEFAULT_SERVER_URL)).pipe(
    filter((it: string | null): it is string => it !== null && it.length > 0),
    shareReplay(1),
  )
  new DataQueryExecutor([new ServerConfigurator("perfint", "installer", serverUrlObservable), new class implements DataQueryConfigurator, FilterConfigurator {
    configureFilter(query: DataQuery): boolean {
      return true
    }
    configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
      configuration.queryProducers.push(new SimpleQueryProducer())
      query.addField({n: "changes", sql: "concat(toString(arrayElement(changes, 1)),' .. ',toString(arrayElement(changes, -1)))"})
      query.addFilter({f: "id", v: infoData.installerId})
      query.order = "changes"
      return true
    }

    createObservable(): Observable<unknown> | null {
      return null
    }
  }]).subscribe((data, configuration) => {
    infoData.changes = data.flat(3)[0] as string
    sidebarVm?.show(infoData)
  })
}

onMounted(() => {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  chartManager = new ChartManager(chartElement.value!, container?.value)
  chartVm = new LineChartVM(
    chartManager,
    dataQueryExecutor,
    props.valueUnit,
  )

  unsubscribe = chartVm.subscribe()

  chartManager.chart.on("click", (params: CallbackDataParams) => {
    if (params.dataIndex != undefined) {
      const infoData = getInfoDataFrom(params, props.valueUnit)
      fetchChangesFromInstaller(infoData)
    }
  })
})

onUnmounted(() => {
  if(unsubscribe != null) unsubscribe()
  chartManager.dispose()
})

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
</script>