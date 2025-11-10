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
    labels: ["Lagging during indexing"],
    measures: [["ui.lagging#average", "ui.lagging#max", "ui.lagging#sum"]],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/curl/indexing", "radler/big_project_50k_10k/indexing"],
    aliases: ["LLVM", "OpenCV", "cURL", "Big Project"],
  },
  {
    labels: ["Lags count during indexing"],
    measures: ["ui.lagging#count"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/curl/indexing", "radler/big_project_50k_10k/indexing"],
    aliases: ["LLVM", "OpenCV", "cURL", "Big Project"],
  },
]

const laggingCompletionCharts: ChartDefinition[] = [
  {
    labels: ["Lagging during completion"],
    measures: [["ui.lagging#average", "ui.lagging#max", "ui.lagging#sum"]],
    projects: [
      "radler/fmtlib/completion/fmt.join_view (dep) (hot)",
      "radler/fmtlib/completion/std.shared_ptr (dep) (hot)",
      "radler/fmtlib/completion/std.string (cold)",
      "radler/fmtlib/completion/std.string (hot)",
    ],
    aliases: ["fmt.join_view", "std.shared_ptr (dep) (hot)", "std.string (cold)", "std.string (hot)"],
  },
  {
    labels: ["Lags count during completion"],
    measures: ["ui.lagging#count"],
    projects: [
      "radler/fmtlib/completion/fmt.join_view (dep) (hot)",
      "radler/fmtlib/completion/std.shared_ptr (dep) (hot)",
      "radler/fmtlib/completion/std.string (cold)",
      "radler/fmtlib/completion/std.string (hot)",
    ],
    aliases: ["fmt.join_view", "std.shared_ptr (dep) (hot)", "std.string (cold)", "std.string (hot)"],
  },
]

const laggingNavigationCharts: ChartDefinition[] = [
  {
    labels: ["Lagging during navigation"],
    measures: [["ui.lagging#average", "ui.lagging#max", "ui.lagging#sum"]],
    projects: [
      "radler/luau/findUsages/class template (DenseHashTable)",
      "radler/luau/findUsages/enumerable (LuauOpcode)",
      "radler/luau/gotoDeclaration/time.h",
      "radler/luau/gotoDeclaration/TypeChecker.getScopes",
    ],
    aliases: ["class template (DenseHashTable)", "enumerable (LuauOpcode)", "time.h", "TypeChecker.getScopes"],
  },
  {
    labels: ["Lags count during navigation"],
    measures: ["ui.lagging#count"],
    projects: [
      "radler/luau/findUsages/class template (DenseHashTable)",
      "radler/luau/findUsages/enumerable (LuauOpcode)",
      "radler/luau/gotoDeclaration/time.h",
      "radler/luau/gotoDeclaration/TypeChecker.getScopes",
    ],
    aliases: ["class template (DenseHashTable)", "enumerable (LuauOpcode)", "time.h", "TypeChecker.getScopes"],
  },
]

const laggingIndexingChartsCombined = combineCharts(laggingIndexingCharts)
const laggingCompletionChartsCombined = combineCharts(laggingCompletionCharts)
const laggingNavigationChartsCombined = combineCharts(laggingNavigationCharts)
</script>
