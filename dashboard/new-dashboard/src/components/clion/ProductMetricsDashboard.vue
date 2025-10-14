<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_product_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
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
    labels: ["First Code Analysis", "File Openings: code loaded", "File Openings: tab shown"],
    measures: ["firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"],
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
    projects: [
      "radler/luau/go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter",
      "radler/luau/go-to-class-with-warmup/CompileOptions/typingLetterByLetter",
      "radler/luau/go-to-file-with-warmup/TableShape.cpp/typingLetterByLetter",
      "radler/luau/go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter",
      "radler/luau/go-to-action-with-warmup/RCP/typingLetterByLetter",
      "radler/luau/go-to-text-with-warmup/LUAU_BUILD_TESTS/typingLetterByLetter",
    ],
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
    // TODO: remove radler/fmtlib/inspection
    projects: ["radler/fmtlib/inspection", "radler/fmtlib/globalInspection"],
    aliases: ["{fmt}"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
