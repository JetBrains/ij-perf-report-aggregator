<template>
  <div
    :id="anchor"
    class="flex flex-col gap-y-2.5 py-3 px-5 border border-solid rounded-md mb-2"
  >
    <h3
      class="m-0 flex items-center"
      @mouseover="labelHovered = true"
      @mouseleave="labelHovered = false"
    >
      {{ title + (settingStore.scaling ? " (scaled)" : "") + (settingStore.removeOutliers ? " (outliers removed)" : "") }}&nbsp;
      <a
        v-show="labelHovered"
        :href="'#' + anchor"
      >
        <LinkIcon class="w-4 h-4" />
      </a>
      <span class="ml-auto flex items-center">
        <span
          v-if="!hasData"
          class="text-sm text-gray-500"
        >
          Missing data. Please check selectors: machine, branch, time range
        </span>
        <span
          v-if="canBeClosed"
          class="text-sm pi pi-plus rotate-45 cursor-pointer transition"
          @click="closeChart"
        />
      </span>
    </h3>
    <div
      v-if="hasData"
      ref="chartElement"
      :style="{ height: `${chartHeight}px` }"
    />
  </div>
</template>
<script setup lang="ts">
import { useElementVisibility } from "@vueuse/core"
import { computed, inject, onMounted, onUnmounted, ref, Ref, useTemplateRef, watch } from "vue"
import { PredefinedMeasureConfigurator, TooltipTrigger } from "../../configurators/MeasureConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { injectOrError, reportInfoProviderKey } from "../../shared/injectionKeys"
import { accidentsConfiguratorKey, containerKey, sidebarVmKey } from "../../shared/keys"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "../common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { useSettingsStore } from "../settings/settingsStore"
import { SeriesNameConfigurator } from "../startup/SeriesNameConfigurator"
import { ChartManager } from "./ChartManager"
import { LineChartVM } from "./LineChartVM"
import { useDarkModeStore } from "../../shared/useDarkModeStore"

interface LineChartProps {
  title: string
  measures: string[]
  configurators: (DataQueryConfigurator | FilterConfigurator)[]
  skipZeroValues?: boolean
  chartType?: ChartType
  valueUnit?: ValueUnit
  tooltipTrigger?: TooltipTrigger
  legendFormatter?: (name: string) => string
  withMeasureName?: boolean
  canBeClosed?: boolean
}

const {
  title,
  measures,
  configurators,
  skipZeroValues = true,
  chartType = "line",
  valueUnit = "ms",
  tooltipTrigger = "item",
  legendFormatter = (name: string) => name,
  withMeasureName = false,
  canBeClosed = false,
} = defineProps<LineChartProps>()

const anchor = computed(() => {
  return title.replaceAll(/[^\dA-Za-z]/g, "")
})

const valueUnitFromMeasure: Ref<ValueUnit> = computed(() => {
  if (measures.every((m) => m.endsWith(".ms"))) {
    return "ms"
  } else if (measures.every((m) => m.endsWith(".ns"))) {
    return "ns"
  } else {
    return valueUnit
  }
})

const settingStore = useSettingsStore()

const accidentsConfigurator = inject(accidentsConfiguratorKey, null)
const chartElement = useTemplateRef<HTMLElement>("chartElement")

const chartIsVisible = useElementVisibility(chartElement)

const skipZeroValuesRef = computed(() => {
  return skipZeroValues
})
const reportInfoProvider = inject(reportInfoProviderKey, null)

const labelHovered = ref(false)
const hasData = ref(true)
const previousHasData = ref(true)

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

let chartManager: ChartManager | null
let chartVm: LineChartVM | null = null
let unsubscribe: (() => void) | null = null

const measuresRef: Ref<string[]> = computed(() => {
  return measures
})

const measureConfigurator = new PredefinedMeasureConfigurator(
  measuresRef,
  skipZeroValuesRef,
  chartType,
  valueUnitFromMeasure.value,
  {
    showSymbol: false,
  },
  accidentsConfigurator,
  tooltipTrigger
)

if (measuresRef.value.length == 1) {
  new SeriesNameConfigurator(measuresRef.value[0])
}

const lineConfigurators = [...configurators, measureConfigurator, infoFieldsConfigurator]
if (withMeasureName) {
  lineConfigurators.push(new SeriesNameConfigurator(measuresRef.value[0]))
}

const dataQueryExecutor = new DataQueryExecutor([...lineConfigurators].filter((item): item is DataQueryConfigurator => item != null))

useDarkModeStore().$subscribe(() => {
  createChart()
})

function createChart() {
  if (chartElement.value) {
    chartManager?.dispose()
    unsubscribe?.()
    chartManager = new ChartManager(chartElement.value, container.value)
    chartVm = new LineChartVM(chartManager, dataQueryExecutor, valueUnitFromMeasure.value, accidentsConfigurator, legendFormatter, (hasDataValue: boolean) => {
      // If transitioning from no data to having data, recreate the chart, otherwise empty chart is shown
      if (!previousHasData.value && hasDataValue) {
        setTimeout(() => {
          createChart()
        }, 0)
      }
      hasData.value = hasDataValue
      previousHasData.value = hasDataValue
    })
    unsubscribe = chartVm.subscribe()
    chartManager.chart.on("click", chartVm.getOnClickHandler(sidebarVm, chartManager, valueUnitFromMeasure.value, accidentsConfigurator))
  }
}

const emit = defineEmits(["chartClosed"])

function closeChart() {
  emit("chartClosed", measuresRef.value)
}

function setupChartOnVisibility() {
  watch(
    chartIsVisible,
    (isVisible) => {
      if (isVisible && chartVm == null) {
        createChart()
      }
    },
    { immediate: true }
  )
  // If the chart is not visible, still try to create it after a delay
  setTimeout(function () {
    if (chartVm == null) {
      createChart()
    }
  }, 1000)
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
