<template>
  <StartupProductDashboard
    product="GO"
    default-project="opentelemetry-collector-contrib"
    initial-machine="Mac Mini M2 Pro (10 vCPU, 32 GB)"
  >
    <template #default="{ configurators }">
      <Divider label="Warm Startup" />
      <LineChart
        title="Updating Go modules dependencies"
        :measures="['metrics.progressMetric/Progress: Updating Go modules dependencies']"
        :configurators="configurators"
        tooltip-trigger="axis"
      />
      <LineChart
        title="Startup and File Analysis"
        :measures="['startUpCompleted', 'metrics.totalOpeningTime/timeFromAppStartTillAnalysisFinished']"
        :configurators="configurators"
        tooltip-trigger="axis"
      />
    </template>
  </StartupProductDashboard>
</template>
<script setup lang="ts">
import LineChart from "../charts/LineChart.vue"
import Divider from "../common/Divider.vue"
import StartupProductDashboard from "../common/StartupProductDashboard.vue"
</script>
