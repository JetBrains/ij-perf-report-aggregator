<template>
  <div
    ref="chartElement"
    class="bg-white"
    :style="{height: `${chartHeight}px`}"
    @mouseenter="show"
    @mouseleave="hide"
  />
</template>
<script setup lang="ts">
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "shared/src/chart"
import { ChartToolTipManager } from "shared/src/components/ChartToolTipManager"
import { PopupTrigger } from "shared/src/components/LineChartManager"
import { ChartType, PredefinedMeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "shared/src/dataQuery"
import { chartToolTipKey, configuratorListKey } from "shared/src/injectionKeys"
import { inject, onMounted, onUnmounted, shallowRef, toRef, watchEffect, withDefaults } from "vue"
import { LineChartVm } from "./LineChartVm"

interface LineChartProps {
  skipZeroValues?: boolean
  compoundTooltip?: boolean
  dataZoom?: boolean
  measures?: Array<string> | null
  chartType?: ChartType
  valueUnit?: ValueUnit
  configurators?: Array<DataQueryConfigurator> | null
  trigger?: PopupTrigger
  aggregatedMeasure: string | null
}

const props = withDefaults(defineProps<LineChartProps>(), {
  skipZeroValues: true,
  compoundTooltip: true,
  dataZoom: false,
  measures: null,
  chartType: "line",
  valueUnit: "ms",
  configurators: null,
  trigger: "axis",
  aggregatedMeasure: null,
})

let chart: ChartManagerHelper
let chartVm: LineChartVm

const chartElement = shallowRef<HTMLElement>()
const skipZeroValues = toRef(props, "skipZeroValues")
const chartToolTipManager = new ChartToolTipManager(props.valueUnit)
// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
const tooltip = inject(chartToolTipKey)!

const show = () => {
  chartVm.showPoints()
}
const hide = () => {
  chartVm.hidePoints()
}
const providedConfigurators = inject(configuratorListKey)

function getConfigurators(): DataQueryConfigurator[] {
  const configurators = props.configurators ?? providedConfigurators

  if (configurators == null) {
    throw new Error(`${configurators} is not provided`)
  }

// static list of measures is provided - create sub data query executor
  if (props.measures != null) {
    configurators.push(
        new PredefinedMeasureConfigurator(
            props.measures,
            skipZeroValues,
            props.chartType,
            props.valueUnit,
            {
              // symbol: "none",
              // symbolSize: 6,
              showSymbol: false,
            },
        )
    )

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
  } else if (props.aggregatedMeasure != null) {
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

  return configurators
}

onMounted(() => {
  const dataQueryExecutor = new DataQueryExecutor(getConfigurators())
  chartToolTipManager.dataQueryExecutor = dataQueryExecutor

  chart = new ChartManagerHelper(chartElement.value!)

  chartVm = new LineChartVm(
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