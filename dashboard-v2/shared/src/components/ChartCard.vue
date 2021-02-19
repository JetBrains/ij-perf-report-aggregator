<template>
  <el-popover
    v-model:visible="infoIsVisible"
    placement="top"
    trigger="manual"
    :auto-close="1000"
    :width="240"
  >
    <template #reference>
      <el-card
        shadow="never"
        :body-style="{ padding: '0px' }"
      >
        <div
          ref="chartElement"
          style="width: 100%; height: 340px;"
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
          style="margin: 10px 0 0;"
        >
          <span
            class="tooltipNameMarker"
            :style='{"background-color": item.color}'
          />
          <span style="margin-left:2px">{{ item.name }}</span>
          <span class="tooltipValue">{{ item.value }}</span>
        </div>
      </div>
    </template>
  </el-popover>
</template>
<script lang="ts">
import { CallbackDataParams } from "echarts/types/src/util/types"
import { defineComponent, inject, onMounted, onUnmounted, PropType, reactive, Ref, ref, shallowRef, toRef, watch } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { LineChartManager } from "../LineChartManager"
import { numberFormat, timeFormat } from "../chart"
import { dataQueryExecutorKey, tooltipUrlProviderKey } from "../componentKeys"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import { encodeQuery } from "../dataQuery"
import { debounceSync } from "../util/debounce"

export default defineComponent({
  name: "ChartCard",
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
    // not reactive - change of initial value is ignored by intention
    measures: {
      type: Array as PropType<Array<string>>,
      default: () => [],
    },
  },
  setup(props) {
    const chartElement: Ref<HTMLElement | null> = shallowRef(null)
    let chartManager: LineChartManager | null = null
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    let dataQueryExecutor = props.provider ?? inject(dataQueryExecutorKey)!
    // eslint-disable-next-line vue/no-setup-props-destructure
    const measures = props.measures
    if (measures.length !== 0) {
      // static list of measures is provided - create sub data query executor
      const measureConfigurator = new PredefinedMeasureConfigurator(measures, props.skipZeroValues)
      watch(() => props.skipZeroValues, value => {
        (measureConfigurator as PredefinedMeasureConfigurator).skipZeroValues = value
        dataQueryExecutor.scheduleLoad()
      })
      dataQueryExecutor = dataQueryExecutor.createSub([measureConfigurator])
      dataQueryExecutor.scheduleLoad()
    }

    // inject is used instead of prop because on dashboard page there are a lot of ChartCards and it is tedious to set property for each
    const chartToolTipManager = new ChartToolTipManager(dataQueryExecutor)
    onMounted(() => {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      chartManager = new LineChartManager(chartElement.value!, dataQueryExecutor, toRef(props, "dataZoom"),
        chartToolTipManager.formatArrayValue.bind(chartToolTipManager))
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

class ChartToolTipManager {
  private readonly tooltipUrlProvider = inject(tooltipUrlProviderKey, null)
  readonly infoIsVisible = ref(false)

  readonly reportTooltipData = reactive<TooltipData>({items: [], linkText: "", linkUrl: null})

  readonly scheduleTooltipHide = debounceSync(() => {
    this.infoIsVisible.value = false
  }, 2_000)

  constructor(private readonly dataQueryExecutor: DataQueryExecutor) {
  }

  formatArrayValue(params: CallbackDataParams | Array<CallbackDataParams>, _ticket: string) {
    const query = this.dataQueryExecutor.lastQuery
    if (query == null) {
      return null
    }

    const p = params as Array<CallbackDataParams>
    const reportTooltipData = this.reportTooltipData
    reportTooltipData.items = p.map(function (measure) {
      return {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        name: measure.seriesName!,
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        value: numberFormat.format((measure.value as Array<number>)[measure.encode!["y"][0]]),
        color: measure.color as string,
      }
    })
    reportTooltipData.linkText = timeFormat.format((p[0].value as Array<number>)[0])
    reportTooltipData.linkUrl = this.tooltipUrlProvider == null ? null : `/api/v1/report/${encodeQuery(query)}`
    this.infoIsVisible.value = true
    this.scheduleTooltipHide()
    return null
  }
}

interface TooltipData {
  linkText: string
  linkUrl: string | null
  items: Array<TooltipDataItem>
}

interface TooltipDataItem {
  name: string
  value: string
  color: string
}

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