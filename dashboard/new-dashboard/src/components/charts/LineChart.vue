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
      <i
        v-if="resolvedDescription"
        role="button"
        tabindex="0"
        aria-label="Show metric details"
        :aria-expanded="infoOpen"
        class="pi pi-info-circle text-sm cursor-pointer"
        @click="toggleInfo"
        @keydown.enter="toggleInfo"
        @keydown.space.prevent="toggleInfo"
      />
      <a
        v-show="labelHovered"
        :href="'#' + anchor"
        class="ml-2"
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
    <Popover
      ref="infoPanel"
      append-to="body"
      @show="infoOpen = true"
      @hide="infoOpen = false"
    >
      <div class="flex flex-col gap-1 text-sm w-fit max-w-md">
        <span class="font-semibold whitespace-nowrap">{{ metricLabel }}</span>
        <hr class="my-1 w-full border-gray-200 dark:border-gray-600" />
        <span v-if="resolvedDescription">{{ resolvedDescription }}</span>
        <span class="text-xs opacity-80">unit: {{ resolvedUnit }}</span>
        <span
          v-if="betterDirectionLabel"
          class="text-xs opacity-80"
          >{{ betterDirectionLabel.arrow }} {{ betterDirectionLabel.label }}</span
        >
        <a
          v-if="metricInfo?.url"
          :href="metricInfo.url"
          target="_blank"
          rel="noopener"
          class="mt-1 text-xs underline decoration-dotted hover:no-underline"
          >Open documentation ↗</a
        >
      </div>
    </Popover>
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
import { PopoverMethods } from "primevue/popover"
import { PredefinedMeasureConfigurator, TooltipTrigger } from "../../configurators/MeasureConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { injectOrError, reportInfoProviderKey } from "../../shared/injectionKeys"
import { accidentsConfiguratorKey, containerKey, sidebarVmKey } from "../../shared/keys"
import { DataQueryExecutor } from "../common/DataQueryExecutor"
import { ChartType, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "../common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"
import { resolveMeasureUnit } from "../common/formatter"
import { useSettingsStore } from "../settings/settingsStore"
import { BetterDirection } from "../../shared/changeDetector/algorithm"
import { measureNameToLabel } from "../../shared/metricsMapping"
import { getMetricDescription } from "../../shared/metricsDescription"
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
  description?: string
  betterDirection?: BetterDirection
}

const {
  title,
  measures,
  configurators,
  skipZeroValues = true,
  chartType = "line",
  valueUnit = "auto",
  tooltipTrigger = "item",
  legendFormatter = (name: string) => name,
  withMeasureName = false,
  canBeClosed = false,
  description,
  betterDirection: betterDirectionProp,
} = defineProps<LineChartProps>()

// One central lookup feeds both the trend-direction rule and the header info popup; the matching prop
// stays an override so a chart can still state a context-specific direction or description.
const metricInfo = computed(() => getMetricDescription(measures[0]))
const resolvedBetterDirection = computed<BetterDirection>(() => betterDirectionProp ?? metricInfo.value?.betterDirection ?? "lower")
const resolvedDescription = computed(() => description ?? metricInfo.value?.description)

const valueUnitFromMeasure: Ref<ValueUnit> = computed(() => {
  if (measures.every((m) => m.endsWith(".ms"))) {
    return "ms"
  } else if (measures.every((m) => m.endsWith(".ns"))) {
    return "ns"
  }
  return valueUnit
})

function describeBetterDirection(direction: BetterDirection): { arrow: string; label: string } | null {
  switch (direction) {
    case "higher":
      return { arrow: "↑", label: "Higher is better" }
    case "lower":
      return { arrow: "↓", label: "Lower is better" }
    case "stable":
      return { arrow: "→", label: "Should stay stable" }
    // "none": changes are detected but never judged good or bad, so no direction is shown
    case "none":
      return null
  }
}

// The structured info popup binds these directly in the template, so the panel content stays a real Vue
// tree (reactive, escaped, keyboard-reachable) rather than an injected HTML string.
const metricLabel = computed(() => measureNameToLabel(measures[0]))
const resolvedUnit = computed(() => resolveMeasureUnit(measures[0], { valueUnit: valueUnitFromMeasure.value, scaling: false }))
const betterDirectionLabel = computed(() => describeBetterDirection(resolvedBetterDirection.value))

// A Popover (not a tooltip) holds the doc link: it is dismissable, so it stays open until the user clicks
// away or presses Escape, which lets the link be clicked. toggle(event) aligns it to the clicked icon.
const infoPanel = useTemplateRef<PopoverMethods>("infoPanel")
const infoOpen = ref(false)
function toggleInfo(event: Event) {
  infoPanel.value?.toggle(event)
}

const anchor = computed(() => {
  return title.replaceAll(/[^\dA-Za-z]/g, "")
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
  tooltipTrigger,
  resolvedBetterDirection.value
)

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
    chartVm = new LineChartVM(chartManager, dataQueryExecutor, valueUnitFromMeasure.value, measures, accidentsConfigurator, legendFormatter, title, (hasDataValue: boolean) => {
      // If transitioning from no data to having data, recreate the chart, otherwise empty chart is shown
      if (!previousHasData.value && hasDataValue) {
        setTimeout(() => {
          createChart()
        }, 0)
      }
      hasData.value = hasDataValue
      previousHasData.value = hasDataValue
    })
    chartVm.enableSidebarAutoOpen({
      sidebarVm,
      valueUnit: valueUnitFromMeasure.value,
      accidentsConfigurator,
    })
    unsubscribe = chartVm.subscribe()
    chartManager.chart.on("click", chartVm.getOnClickHandler(sidebarVm, chartManager, valueUnitFromMeasure.value, accidentsConfigurator))
  }
}

const emit = defineEmits<{
  chartClosed: [measures: string[]]
}>()

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
  setTimeout(() => {
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
