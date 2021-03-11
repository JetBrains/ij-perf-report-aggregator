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
import { CallbackDataParams } from "echarts/types/src/util/types"
import { defineComponent, inject, onMounted, onUnmounted, PropType, reactive, Ref, ref, shallowRef, toRef, watch, watchEffect } from "vue"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { formatDuration, LineChartManager } from "../LineChartManager"
import { DEFAULT_LINE_CHART_HEIGHT, timeFormat } from "../chart"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import { DataQueryFilter, encodeQuery } from "../dataQuery"
import { dataQueryExecutorKey, serverUrlKey, tooltipUrlProviderKey } from "../injectionKeys"
import { debounceSync } from "../util/debounce"

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
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const serverUrl = inject(serverUrlKey)!
    const skipZeroValues = toRef(props, "skipZeroValues")
    const chartToolTipManager = new ChartToolTipManager(serverUrl)
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

class ChartToolTipManager {
  constructor(private readonly serverUrl: Ref<string>) {
  }

  public dataQueryExecutor!: DataQueryExecutor

  private readonly tooltipUrlProvider = inject(tooltipUrlProviderKey, null)
  readonly infoIsVisible = ref(false)

  readonly reportTooltipData = reactive<TooltipData>({items: [], linkText: "", linkUrl: null})

  readonly scheduleTooltipHide = debounceSync(() => {
    this.infoIsVisible.value = false
  }, 2_000)

  formatArrayValue(params: Array<CallbackDataParams>) {
    const query = this.dataQueryExecutor.lastQuery
    if (query == null) {
      return null
    }

    const reportTooltipData = this.reportTooltipData
    reportTooltipData.items = params.map(function (measure) {
      return {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        name: measure.seriesName!,
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        value: formatDuration((measure.value as Array<number>)[measure.encode!["y"][0]]),
        color: measure.color as string,
      }
    })
    const generatedTime = (params[0].value as Array<number>)[0]
    reportTooltipData.linkText = timeFormat.format(generatedTime)
    if (this.tooltipUrlProvider == null) {
      reportTooltipData.linkUrl = null
    }
    else {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
      let q: { filters: Array<DataQueryFilter> } = JSON.parse(JSON.stringify(query, function replacer(key, value) {
        if (key === "fields" || key === "order" || key === "flat") {
          return undefined
        }
        else {
          // eslint-disable-next-line @typescript-eslint/no-unsafe-return
          return value
        }
      }))
      q.filters = q.filters.filter(it => it.field !== "measures.name")
      for (const filter of q.filters) {
        if (filter.field === "generated_time") {
          filter.value = generatedTime / 1000
          delete filter.sql
          break
        }
      }
      reportTooltipData.linkUrl = `/report?reportUrl=${encodeURIComponent(this.serverUrl.value)}/api/v1/report/${encodeQuery(q as never)}`
    }
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