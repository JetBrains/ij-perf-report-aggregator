<template>
  <StartupProductDashboard
    product="GO"
    default-project="kratos"
    initial-machine="Linux Munich i7-13700, 64 Gb"
  >
    <template #default="{ configurators }">
      <Divider label="go list" />
      <LineChart
        title="Execution time: Updating Go modules dependencies"
        :measures="['metrics.progressMetric/Progress: Updating Go modules dependencies']"
        :configurators="configurators"
      />
      <LineChart
        title="Execution time: go list -m -json"
        :measures="['metrics.progressMetric/Subtask: go list -m -json']"
        :configurators="configurators"
      />
      <LineChart
        title="Median value: go list -m -json"
        :measures="['metrics.progressMetric/Subtask: go list -m -json#mean_value']"
        :configurators="configurators"
      />
      <LineChart
        title="Number of go list -m -json runs"
        :measures="['metrics.progressMetric/Subtask: go list -m -json#count']"
        :configurators="configurators"
      />
      <Divider label="Warm Startup" />
      <LineChart
        title="Startup Completed"
        :measures="['startUpCompleted']"
        :configurators="configurators"
      />
      <LineChart
        title="Time from IDE start till File Analysis Finished"
        :measures="['metrics.totalOpeningTime/timeFromAppStartTillAnalysisFinished']"
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
