<template>
  <OverlayPanel
    ref="op"
    :show-close-icon="true"
  >
    <div
      @mouseenter="tooltipEnter"
      @mouseleave="tooltipLeave"
    >
      <a
        v-if="reportTooltipData.linkUrl != null"
        :href="reportTooltipData.linkUrl"
        target="_blank"
        type="primary"
      >
        {{ reportTooltipData.linkText }}
      </a>
      <span v-else>{{ reportTooltipData.linkText }}</span>

      <a
        v-if="reportTooltipData.firstSeriesData.length >= 3"
        title="Changes"
        :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${reportTooltipData.firstSeriesData[2]}&tab=buildChangesDiv`"
        target="_blank"
        class="info"
      >
        changes
      </a>

      <a
        v-if="reportTooltipData.firstSeriesData.length >= 4"
        title="Test Artifacts"
        :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${reportTooltipData.firstSeriesData[3]}&tab=artifacts`"
        target="_blank"
        class="info"
      >
        artifacts
      </a>

      <div
        v-for="item in reportTooltipData.items"
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
        v-if="reportTooltipData.firstSeriesData.length >= 7"
        style="margin: 10px 0 0;white-space: nowrap"
      >
        <span
          class="tooltipNameMarker"
          :style='{"background-color": "green"}'
        />
        <span style="margin-left:2px;">Build Number</span>
        <span class="tooltipValue">{{ reportTooltipData.firstSeriesData[4] }}.{{
            reportTooltipData.firstSeriesData[5]
          }}{{ reportTooltipData.firstSeriesData[6] !== 0 ? "." + reportTooltipData.firstSeriesData[6] : "" }}</span>
      </div>
      <div
        v-if="reportTooltipData.firstSeriesData.length >= 8"
        style="margin: 10px 0 0;white-space: nowrap"
      >
        <span
          class="tooltipNameMarker"
          :style='{"background-color": "blue"}'
        />
        <span style="margin-left:2px;">Machine</span>
        <span class="tooltipValue">{{ reportTooltipData.firstSeriesData[7] }}</span>
      </div>
    </div>
  </OverlayPanel>
  <div
    ref="chartElement"
    class="bg-white overflow-hidden shadow rounded-lg w-full"
    :style="{height: `${chartHeight}px`}"
    @mouseenter="show"
    @mouseleave="hide"
  />
</template>
<script lang="ts">
import { defineComponent, inject, onMounted, onUnmounted, PropType, ref, Ref, shallowRef, toRef, watch, watchEffect } from "vue"
import { DEFAULT_LINE_CHART_HEIGHT } from "../chart"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { dataQueryExecutorKey } from "../injectionKeys"
import { debounceSync } from "../util/debounce"
import { ChartToolTipManager } from "./ChartToolTipManager"
import { LineChartManager } from "./LineChartManager"

export interface Tooltip {
  show(event: any, ref: HTMLElement | null): void

  hide(event: any, ref: HTMLElement | null): void
}

export default defineComponent({
  name: "LineChartCard",
  props: {
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
  },
  setup(props) {
    const chartElement: Ref<HTMLElement | null> = shallowRef(null)
    let chartManager: LineChartManager | null = null
    const providedDataQueryExecutor = inject(dataQueryExecutorKey, null)
    const skipZeroValues = toRef(props, "skipZeroValues")
    const chartToolTipManager = new ChartToolTipManager()
    const op = ref<Tooltip>()
    const show = event => {
      op.value.show(event, chartElement.value)
    }
    const hide = event => {
      debounceSync(() => {
        op.value.hide(event, chartElement.value)
      }, 2_000)()
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
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      chartManager = new LineChartManager(chartElement.value!, chartToolTipManager.dataQueryExecutor, toRef(props, "dataZoom"),
        chartToolTipManager.formatArrayValue.bind(chartToolTipManager))

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

    return {
      chartElement,
      chartHeight: DEFAULT_LINE_CHART_HEIGHT,
      reportTooltipData: chartToolTipManager.reportTooltipData,
      op, show, hide,
    }
  },
})

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
  font-family: Menlo, Monaco, Consolas, Courier, monospace;
  float: right;
  margin-left: 20px;
}
a {
  text-decoration: none;
}
a.info {
  color: gray;
}
</style>