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

const spaceSpecific = "space/gitLogIndexing"
const spaceSpecificSql = "space/gitLogIndexing-sql"
const spaceCommit = "space/git-commit"

const chartsDeclaration: ChartDefinition[] = [
  //Indexing
  {
    labels: ["Indexing"],
    measures: [["vcs-log-indexing"]],
    projects: [spaceSpecific, spaceSpecificSql],
  },
  {
    labels: ["Number of collected commits"],
    measures: [["vcs-log-indexing#numberOfCommits"]],
    projects: [spaceSpecific, spaceSpecificSql],
  },
  {
    labels: ["Real number of commits through git rev-list --count --all"],
    measures: [["realNumberOfCommits"]],
    projects: [spaceSpecific, spaceSpecificSql],
  },
  {
    labels: ["'vcs-log' directory size in bytes"],
    measures: [["vcs-log-size-bytes"]],
    projects: [spaceSpecific, spaceSpecificSql],
  },
  //Commit
  {
    labels: ["Commit FUS duration"],
    measures: [["git-commit#fusCommitDuration"]],
    projects: [spaceCommit],
  },
  {
    labels: ["Refreshing VCS Log when repositories change (on commit, rebase, checkout branch, etc.)"],
    measures: [["vcs-log-refreshing"]],
    projects: [spaceCommit],
  },
  {
    labels: ["Building a [com.intellij.vcs.log.graph.PermanentGraph] for the list of commits"],
    measures: [["vcs-log-building-graph"]],
    projects: [spaceCommit],
  },
  {
    labels: ["Loading full VCS Log (all commits and references)"],
    measures: [["vcs-log-loading-full-log"]],
    projects: [spaceCommit],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
