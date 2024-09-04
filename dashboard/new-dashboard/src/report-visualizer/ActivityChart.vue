<template>
  <div
    ref="chartContainer"
    class="w-full h-[450px]"
  />
  <small v-show="descriptor.id === 'serviceTimeline' || descriptor.id === 'timeline'"> Dotted border area — async preloading. Solid border area — sync preloading. </small>
</template>

<script lang="ts">
import { PropType, defineComponent, useTemplateRef, toRef, watch } from "vue"
import { debounceSync } from "../util/debounce"
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
    },
  },
  setup(props) {
    const descriptorRef = toRef(props, "descriptor")
    const chartContainer = useTemplateRef<HTMLElement>("chartContainer")

    const chartHelper = new ChartComponent(chartContainer, async function (container): Promise<ChartManager> {
      const descriptor = descriptorRef.value

      const sourceNames = descriptor.sourceNames
      if (descriptor.chartManagerProducer == null) {
        const names = sourceNames ?? [descriptor.id]
        const hasALotOfData = !names.includes("reopeningEditors")
        const { ActivityBarChartManager: ActivityBarChartManager } = await import("./charts/ActivityBarChartManager")
        return new ActivityBarChartManager(
          container,
          (dataManager) => {
            const result: GroupedItems = []
            for (const name of names) {
              const data = dataManager.data as never
              const items = data[name] as ItemV20[] | null
              if (items != null) {
                result.push({ category: name, items })
              }
            }
            return result
          },
          {
            unitConverter: UnitConverter.MILLISECONDS,
            shortenName: hasALotOfData,
            threshold: hasALotOfData ? undefined : 0,
          }
        )
      } else {
        return descriptor.chartManagerProducer(container, sourceNames ?? [], descriptor)
      }
    })

    watch(
      [descriptorRef, chartContainer],
      debounceSync(() => {
        const oldChartManager = chartHelper.chartManager
        if (oldChartManager != null) {
          oldChartManager.dispose()
          chartHelper.chartManager = null
        }

        chartHelper.requestRender()
      }, 0)
    )

    return {
      chartContainer: chartContainer,
    }
  },
})
</script>
