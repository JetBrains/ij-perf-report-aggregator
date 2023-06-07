<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_gradle_dashboard"
    initial-machine="Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsDeclaration = [
  "gradle.sync.duration",
  "GRADLE_CALL",
  "PROJECT_RESOLVERS",
  "DATA_SERVICES",
  "WORKSPACE_MODEL_APPLY",

  "AWTEventQueue.dispatchTimeTotal",
  "CPU | Load |Total % 95th pctl",
  "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
  "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
  "gcPause",
  "gcPauseCount",
  "fullGCPause",
  "freedMemoryByGC",
  "totalHeapUsedMax",
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      "project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup",
      "project-import-gradle-micronaut/measureStartup",
      "project-import-gradle-hibernate-orm/measureStartup",
      "project-import-gradle-cas/measureStartup",
      "project-reimport-gradle-cas/measureStartup",
      "project-import-from-cache-gradle-cas/measureStartup",
      "project-import-gradle-500-modules/measureStartup",
      "project-import-gradle-1000-modules/measureStartup",
      "project-import-gradle-5000-modules/measureStartup",
      "project-import-gradle-android-extra-large/measureStartup",
      "project-reimport-space/measureStartup",
      "project-import-space/measureStartup",
      "project-import-open-telemetry/measureStartup",
      "project-import-gradle-openliberty/measureStartup",
    ],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
