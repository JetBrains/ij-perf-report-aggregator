<template>
  <div
    :ref="chartContainer"
    class="timeLineChart"
  />
</template>

<script lang="ts">
import { defineComponent, shallowRef } from "vue"
import {ChartComponent} from "../charts/ChartComponent"
import {TimelineChartManager} from "./TimeLineChartManager"

export default defineComponent({
  name: "TimelineChart",
  setup() {
    const chartContainer = shallowRef<HTMLElement | null>(null)
    new ChartComponent(async function() {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      return new TimelineChartManager(chartContainer.value!)
    })
    return {
      chartContainer,
    }
  }
})
</script>
