<template>
  <div
    ref="chartElement"
    class="bg-white overflow-hidden shadow rounded-lg w-full"
    :style="{height: `${height}px`}"
  />
</template>
<script lang="ts">
import { PropType, defineComponent, inject, onMounted, onUnmounted, Ref, shallowRef } from "vue"
import { BarChartManager } from "../BarChartManager"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { chartDefaultStyle } from "../chart"
import { PredefinedGroupingMeasureConfigurator } from "../configurators/PredefinedGroupingMeasureConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyle, dataQueryExecutorKey, timeRangeKey } from "../injectionKeys"

export default defineComponent({
  name: "BarChartCard",
  props: {
    provider: {
      type: DataQueryExecutor,
      default: () => null,
    },
    height: {
      type: Number,
      default: 440,
    },
    // not reactive - change of initial value is ignored by intention
    measures: {
      type: Array as PropType<Array<string>>,
      default: () => [],
    },
  },
  setup(props) {
    const chartElement: Ref<HTMLElement | null> = shallowRef(null)
    let chartManager: BarChartManager | null = null
    // eslint-disable-next-line vue/no-setup-props-destructure
    const measures = props.measures

    const timeRange = inject(timeRangeKey)
    if (timeRange === undefined) {
      throw new Error("timeRange is not injected but required")
    }

    const aggregationOperatorConfigurator = inject(aggregationOperatorConfiguratorKey)
    if (aggregationOperatorConfigurator === undefined) {
      throw new Error("aggregationOperatorConfigurator is not injected but required")
    }

    const measureConfigurator = new PredefinedGroupingMeasureConfigurator(measures, timeRange, inject(chartStyle, chartDefaultStyle))
    let dataQueryExecutor = props.provider ?? inject(dataQueryExecutorKey)
    dataQueryExecutor = dataQueryExecutor.createSub([aggregationOperatorConfigurator, measureConfigurator])
    dataQueryExecutor.init()

    onMounted(() => {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      chartManager = new BarChartManager(chartElement.value!, dataQueryExecutor)
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
    }
  }
})
</script>