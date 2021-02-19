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
          :href="reportTooltipData.linkUrl"
          target="_blank"
          type="primary"
        >
          {{ reportTooltipData.linkText }}
        </el-link>
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
/* eslint-disable */
import { defineComponent, onMounted, onUnmounted, reactive, Ref, ref, toRef, watch } from "vue"
import { ChartManager, numberFormat, timeFormat } from "../ChartManager"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"
import ElPopover from "element-plus/es/el-popover"
import ElLink from "element-plus/es/el-link"
import { CallbackDataParams } from "echarts/types/src/util/types"
import { debounceSync } from "../util/debounce"
import { encodeQuery } from "../dataQuery"

export default defineComponent({
  name: "ChartCard",
  components: {
    ElPopover,
    ElLink,
  },
  props: {
    provider: {
      type: Object,
      required: true,
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
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const infoIsVisible = ref(false)

    const scheduleTooltipHide = debounceSync(() => {
      infoIsVisible.value = false
    }, 2_000)

    const chartElement: Ref<HTMLElement | null> = ref(null)
    let chartManager: ChartManager | null = null
    let dataQueryExecutor = props.provider as DataQueryExecutor
    const measures = props.measures as Array<string>
    if (measures.length !== 0) {
      // static list of measures is provided - create sub data query executor
      const measureConfigurator = new PredefinedMeasureConfigurator(measures, props.skipZeroValues)
      dataQueryExecutor = dataQueryExecutor.createSub([measureConfigurator])

      watch(() => props.skipZeroValues, value => {
        measureConfigurator.skipZeroValues = value
        dataQueryExecutor.scheduleLoad()
      })

      dataQueryExecutor.scheduleLoad()
    }

    const reportTooltipData = reactive<TooltipData>({items: [], linkText: "", linkUrl: ""})

    function formatArrayValue(params: CallbackDataParams | Array<CallbackDataParams>, _ticket: string) {
      const query = dataQueryExecutor.lastQuery
      if (query == null) {
        return null
      }

      const p = params as Array<CallbackDataParams>
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
      reportTooltipData.linkUrl = `/api/v1/report/${encodeQuery(query)}`
      infoIsVisible.value = true
      scheduleTooltipHide()
      return null
    }

    onMounted(() => {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      chartManager = new ChartManager(chartElement.value!, dataQueryExecutor, toRef(props, "dataZoom"), formatArrayValue)
    })
    onUnmounted(() => {
      const it = chartManager
      if (it != null) {
        chartManager = null
        it.dispose()
      }
    })

    return {
      reportTooltipData,
      infoIsVisible,
      chartElement,
      hideTooltipOnCloseLink() {
        infoIsVisible.value = false
      },
      tooltipEnter: scheduleTooltipHide.clear,
      tooltipLeave: scheduleTooltipHide,
    }
  }
})

interface TooltipData {
  linkText: string
  linkUrl: string
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