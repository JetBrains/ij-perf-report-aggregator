<!-- Copyright 2000-2019 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
<template>
  <div
    ref="chartContainer"
    class="timeLineChart"
  />
</template>

<script lang="ts">
import { defineComponent, shallowRef } from "vue"
import {ChartComponent} from "../charts/ChartComponent"
import {ServiceTimeLineChartManager} from "./ServiceTimeLineChartManager"

export default defineComponent({
  name: "ServiceTimelineChart",
  setup() {
    const chartContainer = shallowRef<HTMLElement | null>(null)
    new ChartComponent(function() {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      return Promise.resolve(new ServiceTimeLineChartManager(chartContainer.value!))
    })
    return {
      chartContainer,
    }
  }
})
</script>
