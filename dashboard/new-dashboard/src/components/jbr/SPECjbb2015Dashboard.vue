<template>
  <DashboardPage
    db-name="jbr"
    table="report"
    persistent-id="jbr_team_dashboard"
    :with-installer="false"
    :is-build-number-exists="true"
  >
    <JbrPlatformCharts
      v-for="metric in metricsNames"
      :key="metric"
      :measure="metric"
      benchmark="SPECjbb2015"
    />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import JbrPlatformCharts from "./JbrPlatformCharts.vue"

const metricsNames = ["jbb2015.result.metric.max-jOPS", "jbb2015.result.metric.critical-jOPS"]
</script>
