<template>
  <StartupProductDashboard
    product="GO"
    default-project="opentelemetry-go-list-parallel"
    initial-machine="Linux Munich i7-13700, 64 Gb"
  >
    <template #default="{ configurators }">
      <LineChart
        title="Updating Go modules dependencies"
        :measures="['metrics.progressMetric/Progress: Updating Go modules dependencies']"
        :configurators="configurators"
      />
    </template>
  </StartupProductDashboard>
</template>
<script setup lang="ts">
import LineChart from "../charts/LineChart.vue"
import StartupProductDashboard from "../common/StartupProductDashboard.vue"
</script>
