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
    labels: ["LoadingDetails - the time spent reading  batch of commits from git  (git log command)"],
    measures: [["LoadingDetails"]],
    projects: [spaceSpecific, spaceSpecificSql],
  },
  //Commit
  {
    labels: ["Commit FUS duration"],
    measures: [["git-commit#fusCommitDuration"]],
    projects: [spaceCommit],
  },
  {
    labels: ["refresh - time spent on reading latest commits from git and adding them to existing"],
    measures: [["refresh"]],
    projects: [spaceCommit],
  },
  {
    labels: ["building graph - time spent on building commit graph(PermanentGraph)"],
    measures: [["building graph"]],
    projects: [spaceCommit],
  },
  {
    labels: ["full log reload - time spent on building commit graph(PermanentGraph) + reading whole repository from git"],
    measures: [["full log reload"]],
    projects: [spaceCommit],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
