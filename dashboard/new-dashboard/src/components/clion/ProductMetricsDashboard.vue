<template>
  <DashboardPage
    db-name="perfint"
    table="clion"
    persistent-id="clion_product_dashboard"
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
    measures: ["backendIndexingTimeMs"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/curl/indexing"],
    aliases: ["LLVM", "OpenCV", "cURL"],
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: ["firstCodeAnalysis"],
    projects: ["radler/fmtlib/typing/simple (4 lines)"],
    aliases: ["{fmt}"],
  },
  {
    labels: ["Completion"],
    measures: ["fus_time_to_show_90p"],
    projects: [
      "radler/fmtlib/completion/std.string (cold)",
      "radler/fmtlib/completion/std.string (hot)",
      "radler/fmtlib/completion/std.shared_ptr (dep) (hot)",
      "radler/fmtlib/completion/fmt.join_view (dep) (hot)",
    ],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: [],
  },
  {
    labels: ["Typing AWT Delay"],
    measures: ["test#max_awt_delay"],
    projects: ["radler/fmtlib/typing/simple (4 lines)"],
    aliases: ["{fmt}"],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["radler/fmtlib/inspection"],
    aliases: ["{fmt}"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
