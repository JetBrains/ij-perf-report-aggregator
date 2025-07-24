<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="ideaDev_wsl_performance_dashboard"
    initial-machine="Windows Munich i7-13700, 64 Gb"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["community/indexing", "spring-pet-clinic-gradle/indexing", "spring-pet-clinic-maven/indexing"],
  },
  {
    labels: ["Scanning"],
    measures: ["scanningTimeWithoutPauses"],
    projects: ["community/indexing", "spring-pet-clinic-gradle/indexing", "spring-pet-clinic-maven/indexing"],
  },
  {
    labels: ["Number of indexed files"],
    measures: ["numberOfIndexedFiles"],
    projects: ["community/indexing", "spring-pet-clinic-gradle/indexing", "spring-pet-clinic-maven/indexing"],
  },
  {
    labels: ["Rebuild"],
    measures: ["build_compilation_duration"],
    projects: ["community/rebuild"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
