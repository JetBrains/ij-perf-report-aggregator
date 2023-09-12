<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="vcs_idea_ultimate_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts">
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
    projects: ["intellij_clone_specific_commit/gitLogIndexing", "intellij_clone_specific_commit/gitLogIndexing-sql"]
  },
  {
    labels: ["Number of collected commits"],
    measures: [["vcs-log-indexing#numberOfCommits"]],
    projects: ["intellij_clone_specific_commit/gitLogIndexing", "intellij_clone_specific_commit/gitLogIndexing-sql"]
  },
  {
    labels: ["LoadingDetails - the time spent reading  batch of commits from git  (git log command)"],
    measures: [["LoadingDetails"]],
    projects: ["intellij_clone_specific_commit/gitLogIndexing", "intellij_clone_specific_commit/gitLogIndexing-sql"]
  },
  {
    labels: ["Show file history"],
    measures: [["showFileHistory"]],
    projects: ["intellij_clone_specific_commit/showFileHistory/EditorImpl-phm", "intellij_clone_specific_commit/showFileHistory/EditorImpl-sql", "intellij_clone_specific_commit/showFileHistory/EditorImpl-noindex"]
  },
  {
    labels: ["Checkout"],
    measures: [["git-checkout"]],
    projects: ["intellij_clone_specific_commit/git-checkout"]
  },
  {
    labels: ["Checkout FUS duration"],
    measures: [["git-checkout#fusCheckoutDuration"]],
    projects: ["intellij_clone_specific_commit/git-checkout"]
  },
  {
    labels: ["Checkout FUS VFS refresh duration "],
    measures: [["git-checkout#fusCheckoutDuration"]],
    projects: ["intellij_clone_specific_commit/git-checkout"]
  }
]

const charts = combineCharts(chartsDeclaration)

</script>