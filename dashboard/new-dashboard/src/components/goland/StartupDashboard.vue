<template>
  <StartupProductDashboard
    product="GO"
    default-project="kratos"
  >
    <template #default="{ configurators }">
      <PerformanceLineChart
        title="Updating Go modules dependencies"
        :measures="[
          'metrics.progressMetric/Progress: Updating Go modules dependencies',
          'metrics.progressMetric/Progress: Updating Go modules dependencies#mean_value',
          'metrics.progressMetric/Progress: Updating Go modules dependencies#standard_deviation',
        ]"
        :configurators="configurators"
      />
      <PerformanceLineChart
        title="go list -m -json (median value)"
        :measures="['metrics.progressMetric/Subtask: go list -m -json#mean_value', 'metrics.progressMetric/Subtask: go list -m -json#standard_deviation']"
        :configurators="configurators"
      />
    </template>
  </StartupProductDashboard>
</template>
<script setup lang="ts">
import PerformanceLineChart from "../charts/PerformanceLineChart.vue"
import StartupProductDashboard from "../common/StartupProductDashboard.vue"
</script>
