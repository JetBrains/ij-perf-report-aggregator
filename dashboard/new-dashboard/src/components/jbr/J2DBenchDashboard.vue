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
      benchmark="J2DBench"
    />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import JbrPlatformCharts from "./JbrPlatformCharts.vue"

const metricsNames = ["poly250", "poly250-rand_col", "poly250-AA-rand_col"]
</script>
