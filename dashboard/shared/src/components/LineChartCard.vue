<template>
  <el-popover
    v-model:visible="infoIsVisible"
    placement="top"
    trigger="manual"
    :auto-close="1000"
    width="fit-content"
  >
    <template #reference>
      <div
        ref="chartElement"
        class="bg-white overflow-hidden shadow rounded-lg w-full"
        :style="{height: `${chartHeight}px`}"
      />
    </template>
    <template #default>
      <div
        @mouseenter="tooltipEnter"
        @mouseleave="tooltipLeave"
      >
        <el-link
          type="default"
          style="float: right"
          :underline="false"
          icon="el-icon-close"
          @click.prevent="hideTooltipOnCloseLink"
        />

        <el-space>
          <el-link
            v-if="reportTooltipData.linkUrl != null"
            :href="reportTooltipData.linkUrl"
            target="_blank"
            type="primary"
          >
            {{ reportTooltipData.linkText }}
          </el-link>
          <span v-else>{{ reportTooltipData.linkText }}</span>

          <el-link
            v-if="reportTooltipData.firstSeriesData.length >= 3"
            title="Changes"
            :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${reportTooltipData.firstSeriesData[2]}&tab=buildChangesDiv`"
            target="_blank"
            type="info"
          >
            changes
          </el-link>

          <el-link
            v-if="reportTooltipData.firstSeriesData.length >= 4"
            title="Test Artifacts"
            :href="`https://buildserver.labs.intellij.net/viewLog.html?buildId=${reportTooltipData.firstSeriesData[3]}&tab=artifacts`"
            target="_blank"
            type="info"
          >
            artifacts
          </el-link>
        </el-space>
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
    </template>
  </el-popover>
</template>
<script lang="ts">
import { defineComponent, inject, onMounted, onUnmounted, PropType, Ref, shallowRef, toRef, watch, watchEffect } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { DEFAULT_LINE_CHART_HEIGHT } from "../chart"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { dataQueryExecutorKey } from "../injectionKeys"
import { ChartToolTipManager } from "./ChartToolTipManager"
import { LineChartManager } from "./LineChartManager"

export default defineComponent({
  name: "LineChartCard",
  props: {
    provider: {
      type: DataQueryExecutor,
      default: () => null
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
            }
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
      infoIsVisible: chartToolTipManager.infoIsVisible,
      hideTooltipOnCloseLink() {
        chartToolTipManager.infoIsVisible.value = false
      },
      tooltipEnter: chartToolTipManager.scheduleTooltipHide.clear,
      tooltipLeave: chartToolTipManager.scheduleTooltipHide,
    }
  }
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
</style>