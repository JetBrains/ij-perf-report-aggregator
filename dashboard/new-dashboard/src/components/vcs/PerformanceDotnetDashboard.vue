<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="vcs_idea_ultimate_dashboard"
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
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const dotnetLatest = "dotnet_clone_latest_commit/gitLogIndexing"

const chartsDeclaration: ChartDefinition[] = [
  //Indexing
  {
    labels: ["Indexing"],
    measures: [["vcs-log-indexing"]],
    projects: [dotnetLatest],
  },
  {
    labels: ["Number of collected commits"],
    measures: [["vcs-log-indexing#numberOfCommits"]],
    projects: [dotnetLatest],
  },
  {
    labels: ["Real number of commits through git rev-list --count --all"],
    measures: [["realNumberOfCommits"]],
    projects: [dotnetLatest],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
