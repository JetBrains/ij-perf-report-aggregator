<template>
  <el-card
    shadow="never"
    :body-style="{ padding: '0px' }"
  >
    <div
      ref="chartElement"
      style="width: 100%; height: 300px;"
    />
  </el-card>
</template>
<script lang="ts">
import { defineComponent, onMounted, onUnmounted, ref, Ref } from "vue"
import { ChartManager } from "../ChartManager"
import { DataQueryExecutor } from "../DataQueryExecutor"

export default defineComponent({
  name: "ChartCard",
  props: {
    provider: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const chartElement: Ref<HTMLElement | null> = ref(null)
    let chartManager: ChartManager | null = null
    const dataQueryExecutor = props.provider as DataQueryExecutor
    onMounted(() => {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      chartManager = new ChartManager(chartElement.value!, dataQueryExecutor)
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