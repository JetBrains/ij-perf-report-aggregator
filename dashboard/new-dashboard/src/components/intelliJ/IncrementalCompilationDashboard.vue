<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_incremental_build_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
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
  >
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["IntelliJ Initial Build/ Incremental build", "Mappings Dir Size in Bytes 1 day"],
    measures: ["build_compilation_duration", "mappingsDirSizeInBytes"],
    projects: [
      "incremental-build-intellij/rebuild_initial",
      "incremental-build-intellij/build_incremental-1-day",
      "incremental-build-intellij/build_incremental-2-days",
      "incremental-build-intellij/build_incremental-3-days",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
