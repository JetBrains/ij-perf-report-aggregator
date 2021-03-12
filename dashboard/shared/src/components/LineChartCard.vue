<template>
  <el-popover
    v-model:visible="infoIsVisible"
    placement="top"
    trigger="manual"
    :auto-close="1000"
    width="fit-content"
  >
    <template #reference>
      <el-card
        shadow="never"
        :body-style="{ padding: '0px' }"
      >
        <div
          ref="chartElement"
          :style='{width: "100%", height: `${chartHeight}px`}'
        />
      </el-card>
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

        <el-link
          v-if="reportTooltipData.linkUrl != null"
          :href="reportTooltipData.linkUrl"
          target="_blank"
          type="primary"
        >
          {{ reportTooltipData.linkText }}
        </el-link>
        <span v-else>{{ reportTooltipData.linkText }}</span>

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
      </div>
    </template>
  </el-popover>
</template>
<script lang="ts">
import { defineComponent, inject, onMounted, onUnmounted, PropType, Ref, shallowRef, toRef, watch, watchEffect } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { LineChartManager } from "../LineChartManager"
import { DEFAULT_LINE_CHART_HEIGHT } from "../chart"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import { dataQueryExecutorKey } from "../injectionKeys"
import { ChartToolTipManager } from "./ChartToolTipManager"

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
        dataQueryExecutor = dataQueryExecutor.createSub([new PredefinedMeasureConfigurator(props.measures, skipZeroValues)])
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