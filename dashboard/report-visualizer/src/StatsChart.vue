<!-- Copyright 2000-2019 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file. -->
<template>
  <div
    ref="chartContainer"
    style="height: 500px"
  />
</template>

<script lang="ts">
import { defineComponent, shallowRef } from "vue"
import { ChartComponent } from "./charts/ChartComponent"
import { ChartManager } from "./charts/ChartManager"
import { StatsChartManager } from "./charts/StatsChartManager"

export default defineComponent({
  name: "StatsChart",
  setup() {
    const chartContainer = shallowRef<HTMLElement | null>(null)
    new ChartComponent(async function(): Promise<ChartManager> {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      return new StatsChartManager(chartContainer.value!)
    })
    return {
      chartContainer,
    }
  }
})
</script>
