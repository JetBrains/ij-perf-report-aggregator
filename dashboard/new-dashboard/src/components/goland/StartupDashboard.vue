<template>
  <StartupProductDashboard
    product="GO"
    default-project="kratos"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :hide-additional-metrics="false"
  >
    <template #default="{ configurators }">
      <Divider label="go list" />
      <LineChart
        title="Updating Go modules dependencies"
        :measures="[
          'metrics.progressMetric/Progress: Updating Go modules dependencies',
          'metrics.progressMetric/Progress: Updating Go modules dependencies#mean_value',
          'metrics.progressMetric/Progress: Updating Go modules dependencies#standard_deviation',
        ]"
        :configurators="configurators"
      />
      <LineChart
        title="go list -m -json (median value)"
        :measures="['metrics.progressMetric/Subtask: go list -m -json#mean_value', 'metrics.progressMetric/Subtask: go list -m -json#standard_deviation']"
        :configurators="configurators"
      />
    </template>
  </StartupProductDashboard>
</template>
<script setup lang="ts">
import LineChart from "../charts/LineChart.vue"
import Divider from "../common/Divider.vue"
import StartupProductDashboard from "../common/StartupProductDashboard.vue"
</script>
