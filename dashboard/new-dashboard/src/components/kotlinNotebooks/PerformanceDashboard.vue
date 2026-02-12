<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlinNotebooks"
    persistent-id="kotlinNotebooks_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
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
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: ["kotlin-notebook-perf/notebook-completion-test", "kotlin-notebook-perf/data-frame-completion-test"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
