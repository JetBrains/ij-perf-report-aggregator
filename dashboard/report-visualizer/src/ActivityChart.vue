<template>
  <div
    ref="chartContainer"
    class="activityChart"
  />
  <small v-show="descriptor.id === 'serviceTimeline' || descriptor.id === 'timeline'">
    Dotted border area — async preloading. Solid border area — sync preloading.
  </small>
</template>

<script lang="ts">
import { debounceSync } from "shared/src/util/debounce"
import { PropType , defineComponent, shallowRef, toRef, watch } from "vue"
import { ActivityChartDescriptor } from "./ActivityChartDescriptor"
import { GroupedItems } from "./DataManager"
import { ChartComponent, ChartManager } from "./charts/ChartComponent"
import { ItemV20, UnitConverter } from "./data"

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
        const names = sourceNames ?? [descriptor.id]
        const hasALotOfData = !names.some(it => it === "reopeningEditors")
        return new ((await import("./charts/ActivityBarChartManager")).ActivityBarChartManager)(container, dataManager => {
          const result: GroupedItems = []
          for (const name of names) {
            const data = dataManager.data as never
            const items = data[name] as Array<ItemV20> | null
            if (items != null) {
              result.push({category: name, items})
            }
          }
          return result
        }, {
          unitConverter: UnitConverter.MILLISECONDS,
          shortenName: hasALotOfData,
          threshold: hasALotOfData ? undefined : 0,
        })
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
