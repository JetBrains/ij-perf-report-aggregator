<template>
  <StartupProductDashboard
    product="GO"
    default-project="kratos"
    initial-machine="Linux Munich i7-13700, 64 Gb"
  >
    <template #default="{ configurators }">
      <Divider label="go list" />
      <LineChart
        title="Updating Go modules dependencies"
        :measures="['metrics.progressMetric/Progress: Updating Go modules dependencies', 'metrics.progressMetric/Progress: Updating Go modules dependencies#mean_value']"
        :configurators="configurators"
      />
      <LineChart
        title="go list -m -json"
        :measures="['metrics.progressMetric/Subtask: go list -m -json']"
        :configurators="configurators"
      />
      <LineChart
        title="go list -m -json (median value)"
        :measures="['metrics.progressMetric/Subtask: go list -m -json#mean_value']"
        :configurators="configurators"
      />
      <LineChart
        title="go list -m -json (count)"
        :measures="['metrics.progressMetric/Subtask: go list -m -json#count']"
        :configurators="configurators"
      />
      <Divider label="Warm Startup" />
      <LineChart
        title="Startup Completed"
        :measures="['startUpCompleted']"
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
