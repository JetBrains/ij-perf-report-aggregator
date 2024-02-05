<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_indexing_dashboard"
    initial-machine="linux-blade-hetzner"
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
import { ChartDefinition, combineCharts } from "../../charts/DashboardCharts"
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: ["embeddingIndexing"],
    projects: ["check-semantic-indexing-works"],
  },
  {
    labels: ["Scanning"],
    measures: ["embeddingFilesScanning"],
    projects: ["check-semantic-indexing-works"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
