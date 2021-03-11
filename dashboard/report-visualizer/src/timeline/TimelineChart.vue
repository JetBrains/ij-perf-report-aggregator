<template>
  <div
    ref="chartContainer"
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
    new ChartComponent(function() {
      const value = chartContainer.value
      if (value == null) {
        throw new Error("container is not created")
      }
      return Promise.resolve(new TimelineChartManager(value))
    })
    return {
      chartContainer,
    }
  }
})
</script>
