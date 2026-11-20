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
      benchmark="Dacapo"
    />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import JbrPlatformCharts from "./JbrPlatformCharts.vue"

const metricsNames = [
  "avrora",
  "batik",
  "biojava",
  "eclipse",
  "fop",
  "graphchi",
  "h2",
  "jme",
  "jython",
  "kafka",
  "luindex",
  "lusearch",
  "pmd",
  "spring",
  "sunflow",
  "tomcat",
  "xalan",
  "zxing",
]
</script>
