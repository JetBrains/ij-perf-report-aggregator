<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="vcs_space_ultimate_dashboard"
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

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: [["vcs-log-indexing"]],
    projects: ["space/gitLogIndexing", "space/gitLogIndexing-sql"],
  },
  {
    labels: ["Number of collected commits"],
    measures: [["vcs-log-indexing#numberOfCommits"]],
    projects: ["space/gitLogIndexing", "space/gitLogIndexing-sql"],
  },
  {
    labels: ["LoadingDetails - the time spent reading  batch of commits from git  (git log command)"],
    measures: [["LoadingDetails"]],
    projects: ["space/gitLogIndexing", "space/gitLogIndexing-sql"],
  },
  {
    labels: ["Commit FUS duration"],
    measures: [["git-commit#fusCommitDuration"]],
    projects: ["space/git-commit"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
