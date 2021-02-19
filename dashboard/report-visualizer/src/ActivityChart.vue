<template>
  <div
    ref="chartContainer"
    class="activityChart"
  />
</template>

<script lang="ts">
import { defineComponent, shallowRef, toRef, watch } from "vue"
import { chartDescriptors } from "./charts/ActivityChartDescriptor"
import { ChartComponent } from "./charts/ChartComponent"
import { ChartManager } from "./charts/ChartManager"

export default defineComponent({
  name: "ActivityChart",
  props: {
    type: {
      type: String,
      required: true,
    }
  },
  setup(props) {
    const typeRef = toRef(props, "type")
    const chartContainer = shallowRef<HTMLElement | null>(null)

    const chartHelper = new ChartComponent(async function(): Promise<ChartManager> {
      const type = typeRef.value
      const descriptor = chartDescriptors.find(it => it.id === type)
      if (descriptor === undefined) {
        throw new Error(`Unknown chart type: ${type}`)
      }

      const sourceNames = descriptor.sourceNames
      if (descriptor.chartManagerProducer == null) {
        return new (await import("./charts/ActivityChartManager"))
          // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
          .ActivityChartManager(chartContainer.value!, sourceNames ?? [type], descriptor)
      }
      else {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        return await descriptor.chartManagerProducer(chartContainer.value!, sourceNames!, descriptor)
      }
    })
    watch(typeRef, () => {
      const oldChartManager = chartHelper.chartManager
      if (oldChartManager != null) {
        oldChartManager.dispose()
        chartHelper.chartManager = null
      }

      chartHelper.renderDataIfAvailable()
    })

    return {
      chartContainer,
    }
  }
})
</script>
