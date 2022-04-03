<template>
  <div
    ref="chartElement"
    class="bg-white shadow rounded-lg"
    :style="{height: `${chartHeight}px`}"
    @mouseenter="show"
    @mouseleave="hide"
  />
  <OverlayPanel
    ref="tooltip"
    :show-close-icon="true"
  >
    <div>
      <a
        v-if="tooltipData.linkUrl != null"
        :href="tooltipData.linkUrl"
        target="_blank"
      >
        {{ tooltipData.linkText }}
      </a>
      <span v-else>{{ tooltipData.linkText }}</span>

      <a
        v-if="tooltipData.firstSeriesData.length >= 3"
        title="Changes"
        :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[2]}&tab=buildChangesDiv`"
        target="_blank"
        class="info"
      >
        changes
      </a>

      <a
        v-if="tooltipData.firstSeriesData.length >= 4"
        title="Test Artifacts"
        :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${tooltipData.firstSeriesData[3]}&tab=artifacts`"
        target="_blank"
        class="info"
      >
        artifacts
      </a>

      <div
        v-for="item in tooltipData.items"
        :key="item.name"
        style="margin: 10px 0 0;white-space: nowrap"
      >
        <span
          class="tooltipNameMarker"
          :style='{"background-color": item.color}'
        />
        <span style="margin-left:2px;">{{ item.name }}</span>
        <span class="tooltipValue">{{ item.value }}</span>
      </div>
      <div
        v-if="tooltipData.firstSeriesData.length >= 8"
        style="margin: 10px 0 0;white-space: nowrap"
      >
        <span
          class="tooltipNameMarker"
          :style='{"background-color": "blue"}'
        />
        <span style="margin-left:2px;">machine</span>
        <span class="tooltipValue">{{ tooltipData.firstSeriesData[7] }}</span>
      </div>
    </div>
  </OverlayPanel>
</template>
<script setup lang="ts">
import { computed, inject, onMounted, onUnmounted, PropType, ref, Ref, shallowRef, toRef, watch, watchEffect } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { DEFAULT_LINE_CHART_HEIGHT } from "../chart"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { dataQueryExecutorKey } from "../injectionKeys"
import { debounceSync } from "../util/debounce"
import { ChartToolTipManager } from "./ChartToolTipManager"
import { LineChartManager } from "./LineChartManager"

export interface Tooltip {
  show(event: Event, ref: HTMLElement | null): void

  hide(): void
}

const props = defineProps({
  provider: {
    type: DataQueryExecutor,
    default: () => null,
  },
  skipZeroValues: {
    type: Boolean,
    default: true,
  },
  dataZoom: {
    type: Boolean,
    default: false,
  },
  measures: {
    type: Array as PropType<Array<string> | null>,
    default: () => null,
  },
})

const chartElement: Ref<HTMLElement | null> = shallowRef(null)
let chartManager: LineChartManager | null = null
const providedDataQueryExecutor = inject(dataQueryExecutorKey, null)
const skipZeroValues = toRef(props, "skipZeroValues")
const chartToolTipManager = new ChartToolTipManager()
const tooltip = ref<Tooltip>()

const tooltipData = chartToolTipManager.reportTooltipData

let lastShowEvent: Event | null = null
const debouncedHide = debounceSync(() => {
  lastShowEvent = null
  debouncedShow.clear()
  tooltip.value?.hide()
}, 2_000)
const debouncedShow = debounceSync(() => {
  if (lastShowEvent != null) {
    tooltip.value?.show(lastShowEvent, chartElement.value)
  }
}, 300)
const show = (event: Event) => {
  debouncedHide?.clear()
  lastShowEvent = event
  debouncedShow()
}
const hide = () => {
  debouncedShow.clear()
  lastShowEvent = null
  debouncedHide()
}

watchEffect(function () {
  let dataQueryExecutor = props.provider ?? providedDataQueryExecutor
  if (dataQueryExecutor == null) {
    throw new Error("Neither `provider` property is set, nor `dataQueryExecutor` is provided")
  }

  // static list of measures is provided - create sub data query executor
  if (props.measures != null) {
    const configurators: Array<DataQueryConfigurator> = [
      new PredefinedMeasureConfigurator(props.measures, skipZeroValues),
    ]
    const infoFields = chartToolTipManager.reportInfoProvider?.infoFields ?? []
    if (infoFields.length !== 0) {
      configurators.push({
        configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
          for (const infoField of infoFields) {
            query.addField(infoField)
          }
          return true
        },
      })
    }
    dataQueryExecutor = dataQueryExecutor.createSub(configurators)
    dataQueryExecutor.scheduleLoad()
  }
  chartToolTipManager.dataQueryExecutor = dataQueryExecutor
  if (chartManager != null) {
    chartManager.dataQueryExecutor = dataQueryExecutor
  }
})

onMounted(() => {
  chartManager = new LineChartManager(
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    chartElement.value!,
    chartToolTipManager.dataQueryExecutor,
    toRef(props, "dataZoom"),
    chartToolTipManager.formatArrayValue.bind(chartToolTipManager),
  )

  watch(skipZeroValues, () => {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    chartManager!.dataQueryExecutor.scheduleLoad()
  })
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
<style scoped>
.tooltipNameMarker {
  display: inline-block;
  margin-right: 4px;
  border-radius: 10px;
  width: 10px;
  height: 10px;
}

.tooltipValue {
  @apply font-mono;
  float: right;
  margin-left: 20px;
}

a {
  text-decoration: none;
}

a.info {
  @apply text-gray-600;
}
</style>