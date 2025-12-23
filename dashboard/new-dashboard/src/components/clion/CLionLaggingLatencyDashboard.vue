<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_lagging_latency_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section>
      <Divider title="Lagging during indexing" />
      <GroupProjectsChart
        v-for="chart in laggingIndexingChartsCombined"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
      <Divider title="Lagging during completion" />
      <GroupProjectsChart
        v-for="chart in laggingCompletionChartsCombined"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
      <Divider title="Lagging during navigation" />
      <GroupProjectsChart
        v-for="chart in laggingNavigationChartsCombined"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
      <Divider title="Lagging during browsing files" />
      <GroupProjectsChart
        v-for="chart in laggingHighlightingChartsCombined"
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
import Divider from "../common/Divider.vue"

const laggingIndexingCharts: ChartDefinition[] = [
  {
    labels: ["Lagging during indexing - average, max"],
    measures: [["ui.lagging#average", "ui.lagging#max"]],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/big_project_50k_10k/indexing", "radler/big_project_50k_10k_many_symbols/indexing"],
    aliases: ["LLVM", "OpenCV", "Big Project", "Big Project Many Symbols"],
  },
  {
    labels: ["Lagging during indexing - sum"],
    measures: ["ui.lagging#sum"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/big_project_50k_10k/indexing", "radler/big_project_50k_10k_many_symbols/indexing"],
    aliases: ["LLVM", "OpenCV", "Big Project", "Big Project Many Symbols"],
  },
  {
    labels: ["Lagging during indexing - count"],
    measures: ["ui.lagging#count"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/big_project_50k_10k/indexing", "radler/big_project_50k_10k_many_symbols/indexing"],
    aliases: ["LLVM", "OpenCV", "Big Project", "Big Project Many Symbols"],
  },
]

const laggingCompletionCharts: ChartDefinition[] = [
  {
    labels: ["Lagging during completion - average, max"],
    measures: [["ui.lagging#average", "ui.lagging#max"]],
    projects: ["radler/fmtlib/completion/fmt.join_view (dep) (hot)", "radler/fmtlib/completion/std.shared_ptr (dep) (hot)", "radler/fmtlib/completion/std.string (hot)"],
    aliases: ["fmt.join_view (dep) (hot)", "std.shared_ptr (dep) (hot)", "std.string (hot)"],
  },
  {
    labels: ["Lagging during completion - sum"],
    measures: ["ui.lagging#sum"],
    projects: ["radler/fmtlib/completion/fmt.join_view (dep) (hot)", "radler/fmtlib/completion/std.shared_ptr (dep) (hot)", "radler/fmtlib/completion/std.string (hot)"],
    aliases: ["fmt.join_view (dep) (hot)", "std.shared_ptr (dep) (hot)", "std.string (hot)"],
  },
  {
    labels: ["Lagging during completion - count"],
    measures: ["ui.lagging#count"],
    projects: ["radler/fmtlib/completion/fmt.join_view (dep) (hot)", "radler/fmtlib/completion/std.shared_ptr (dep) (hot)", "radler/fmtlib/completion/std.string (hot)"],
    aliases: ["fmt.join_view (dep) (hot)", "std.shared_ptr (dep) (hot)", "std.string (hot)"],
  },
]

const laggingNavigationCharts: ChartDefinition[] = [
  {
    labels: ["Lagging during navigation - average, max"],
    measures: [["ui.lagging#average", "ui.lagging#max"]],
    projects: ["radler/luau/findUsages/class template (DenseHashTable)", "radler/luau/gotoDeclaration/time.h", "radler/luau/gotoDeclaration/TypeChecker.getScopes"],
    aliases: ["class template (DenseHashTable)", "time.h", "TypeChecker.getScopes"],
  },
  {
    labels: ["Lagging during navigation - sum"],
    measures: ["ui.lagging#sum"],
    projects: ["radler/luau/findUsages/class template (DenseHashTable)", "radler/luau/gotoDeclaration/time.h", "radler/luau/gotoDeclaration/TypeChecker.getScopes"],
    aliases: ["class template (DenseHashTable)", "time.h", "TypeChecker.getScopes"],
  },
  {
    labels: ["Lagging during navigation - count"],
    measures: ["ui.lagging#count"],
    projects: ["radler/luau/findUsages/class template (DenseHashTable)", "radler/luau/gotoDeclaration/time.h", "radler/luau/gotoDeclaration/TypeChecker.getScopes"],
    aliases: ["class template (DenseHashTable)", "time.h", "TypeChecker.getScopes"],
  },
]

const laggingHighlightingCharts: ChartDefinition[] = [
  {
    labels: ["Lagging during browsing - average, max"],
    measures: [["ui.lagging#average", "ui.lagging#max", "ui.lagging#percentage_share"]],
    projects: ["radler/opencv/syntaxHighlighting/opencv"],
    aliases: ["syntaxHighlighting opencv"],
  },
  {
    labels: ["Lagging during browsing - lagging percentage share"],
    measures: [["ui.lagging#percentage_share"]],
    projects: ["radler/opencv/syntaxHighlighting/opencv"],
    aliases: ["syntaxHighlighting opencv"],
  },
]

const laggingIndexingChartsCombined = combineCharts(laggingIndexingCharts)
const laggingCompletionChartsCombined = combineCharts(laggingCompletionCharts)
const laggingNavigationChartsCombined = combineCharts(laggingNavigationCharts)
const laggingHighlightingChartsCombined = combineCharts(laggingHighlightingCharts)
</script>
