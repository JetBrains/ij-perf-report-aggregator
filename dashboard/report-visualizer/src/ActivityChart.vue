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

      const container = chartContainer.value
      if (container == null) {
        throw new Error("container is not created")
      }

      const sourceNames = descriptor.sourceNames
      if (descriptor.chartManagerProducer == null) {
        return new (await import("./charts/ActivityChartManager"))
          .ActivityChartManager(container, sourceNames ?? [type], descriptor)
      }
      else {
        return await descriptor.chartManagerProducer(container, sourceNames ?? [], descriptor)
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
<style scoped>
.activityChart {
  width: 100%;
  /*
  our data has extraordinary high values (extremes) and it makes item chart not readable (extremes are visible and others column bars are too low),
  as solution, amCharts supports breaks (https://www.amcharts.com/demos/column-chart-with-axis-break/),
  but it contradicts to our goal - to show that these items are extremes,
  so, as solution, we increase chart height to give more space to render bars.

  It is ok, as now we use UI Library (ElementUI) and can use Tabs, Collapse and any other component to group charts.
  Also, as we use Vue.js and Vue Router, it is one-line to provide dedicated view (/#/components and so on)
  */
  height: 500px;
}
</style>
