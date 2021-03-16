<template>
  <div
    ref="chartContainer"
    class="activityChart"
  />
</template>

<script lang="ts">
import { debounceSync } from "shared/src/util/debounce"
import { PropType , defineComponent, shallowRef, toRef, watch } from "vue"
import { ActivityChartDescriptor } from "./ActivityChartDescriptor"
import { ChartComponent } from "./charts/ChartComponent"
import { ChartManager } from "./charts/ChartManager"

export default defineComponent({
  name: "ActivityChart",
  props: {
    descriptor: {
      type: Object as PropType<ActivityChartDescriptor>,
      required: true,
    }
  },
  setup(props) {
    const descriptorRef = toRef(props, "descriptor")
    const chartContainer = shallowRef<HTMLElement | null>(null)

    const chartHelper = new ChartComponent(async function(): Promise<ChartManager> {
      const descriptor = descriptorRef.value
      const container = chartContainer.value
      if (container == null) {
        throw new Error("container is not created")
      }

      const sourceNames = descriptor.sourceNames
      if (descriptor.chartManagerProducer == null) {
        return new (await import("./charts/ActivityChartManager"))
          .ActivityChartManager(container, sourceNames ?? [descriptor.id], descriptor)
      }
      else {
        return await descriptor.chartManagerProducer(container, sourceNames ?? [], descriptor)
      }
    })
    watch(descriptorRef, debounceSync(() => {
      const oldChartManager = chartHelper.chartManager
      if (oldChartManager != null) {
        oldChartManager.dispose()
        chartHelper.chartManager = null
      }

      chartHelper.renderDataIfAvailable()
    }, 0))

    return {
      chartContainer,
    }
  }
})
</script>
<style scoped>
.activityChart {
  width: 100%;
  height: 450px;
}
</style>
