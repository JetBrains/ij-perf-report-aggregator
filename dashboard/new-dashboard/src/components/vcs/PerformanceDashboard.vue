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

const intellijSpecific = "intellij_clone_specific_commit/gitLogIndexing"
const intellijSpecificSql = "intellij_clone_specific_commit/gitLogIndexing-sql"
const intellijLatest = "intellij_latest_master/gitLogIndexing"

const showFileHistoryEditorPhm = "intellij_clone_specific_commit/showFileHistory/EditorImpl-phm"
const showFileHistoryEditorSql = "intellij_clone_specific_commit/showFileHistory/EditorImpl-sql"
const showFileHistoryEditorNoIndex = "intellij_clone_specific_commit/showFileHistory/EditorImpl-noindex"

const chartsDeclaration: ChartDefinition[] = [
  //Indexing
  {
    labels: ["Indexing"],
    measures: [["vcs-log-indexing"]],
    projects: [intellijSpecific, intellijSpecificSql, intellijLatest],
  },
  {
    labels: ["Number of collected commits"],
    measures: [["vcs-log-indexing#numberOfCommits"]],
    projects: [intellijSpecific, intellijSpecificSql, intellijLatest],
  },
  {
    labels: ["Real number of collected commits through git rev-list --count --all"],
    measures: [["realNumberOfCommits"]],
    projects: [intellijSpecific, intellijSpecificSql, intellijLatest],
  },
  //Show file history
  {
    labels: ["Show file history (test metric)"],
    measures: [["showFileHistory"]],
    projects: [showFileHistoryEditorPhm, showFileHistoryEditorSql, showFileHistoryEditorNoIndex],
  },
  {
    labels: ["Show file history - showing first pack of data (test metric)"],
    measures: [["showFirstPack"]],
    projects: [showFileHistoryEditorPhm, showFileHistoryEditorSql, showFileHistoryEditorNoIndex],
  },
  {
    labels: [
      "Computing - time spent on computing a peace of history. If index - time of computing before the first rename. " +
        "If git - time of computing before timeout of operation occurred",
    ],
    measures: [["file-history-computing"]],
    projects: [showFileHistoryEditorPhm, showFileHistoryEditorSql, showFileHistoryEditorNoIndex],
  },
  {
    labels: ["Refreshing VCS Log when repositories change (on commit, rebase, checkout branch, etc.)"],
    measures: [["vcs-log-refreshing"]],
    projects: [showFileHistoryEditorPhm, showFileHistoryEditorSql, showFileHistoryEditorNoIndex],
  },
  {
    labels: ["Building a [com.intellij.vcs.log.graph.PermanentGraph] for the list of commits"],
    measures: [["vcs-log-building-graph"]],
    projects: [showFileHistoryEditorPhm, showFileHistoryEditorSql, showFileHistoryEditorNoIndex],
  },
  {
    labels: ["Loading full VCS Log (all commits and references)"],
    measures: [["vcs-log-loading-full-log"]],
    projects: [showFileHistoryEditorPhm, showFileHistoryEditorSql, showFileHistoryEditorNoIndex],
  },
  //Checkout
  {
    labels: ["Checkout"],
    measures: [["git-checkout"]],
    projects: ["intellij_clone_specific_commit/git-checkout"],
  },
  {
    labels: ["Checkout FUS duration"],
    measures: [["git-checkout#fusCheckoutDuration"]],
    projects: ["intellij_clone_specific_commit/git-checkout"],
  },
  {
    labels: ["Checkout FUS VFS refresh duration "],
    measures: [["git-checkout#fusVfsRefreshDuration"]],
    projects: ["intellij_clone_specific_commit/git-checkout"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
