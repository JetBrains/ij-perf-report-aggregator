<template>
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
<script lang="ts">
import { defineComponent, onMounted, onUnmounted, ref, Ref, watch, toRef } from "vue"
import { ChartManager } from "../ChartManager"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { PredefinedMeasureConfigurator } from "../configurators/MeasureConfigurator"

export default defineComponent({
  name: "ChartCard",
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
    onMounted(() => {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      chartManager = new ChartManager(chartElement.value!, dataQueryExecutor, toRef(props, "dataZoom"))
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