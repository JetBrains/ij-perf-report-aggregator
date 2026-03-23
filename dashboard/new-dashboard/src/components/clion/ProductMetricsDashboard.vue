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
        v-for="chart in mainCharts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
    <Divider title="SearchEverywhere - New vs Old" />
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div
        v-for="chart in seNewVsOldCharts"
        :key="chart.definition.label"
        class="flex-1 min-w-0"
      >
        <GroupProjectsChart
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

const mainChartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: ["backendIndexingTimeMs"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/curl/indexing", "radler/zephyr_bap_broadcast_sink/indexing"],
    aliases: ["LLVM", "OpenCV", "cURL", "Zephyr Bap Broadcast Sink"],
  },
  {
    labels: ["First Code Analysis", "File Openings: code loaded", "File Openings: tab shown"],
    measures: ["firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"],
    projects: ["radler/fmtlib/typing/simple (4 lines)"],
    aliases: ["{fmt}"],
  },

  {
    labels: ["Typing latency"],
    measures: ["typing#latency#mean_value"],
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
    labels: ["SearchEverywhere - Wait Until Full or Done"],
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
    labels: ["SearchEverywhere - First Elements Added"],
    measures: ["searchEverywhere_first_elements_added"],
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
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["radler/fmtlib/globalInspection"],
    aliases: ["{fmt}"],
  },
]

const seNewVsOldChartsDeclaration: ChartDefinition[] = [
  {
    labels: ["SearchEverywhere - New vs Old - All"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter", "radler/luau/new-se-go-to-all-with-warmup/AstJsonEncoder/typingLetterByLetter"],
  },
  {
    labels: ["SearchEverywhere - New vs Old - Symbol"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter", "radler/luau/new-se-go-to-symbol-with-warmup/Type_Boolean/typingLetterByLetter"],
  },
  {
    labels: ["SearchEverywhere - New vs Old - Class"],
    measures: ["searchEverywhere"],
    projects: ["radler/luau/go-to-class-with-warmup/CompileOptions/typingLetterByLetter", "radler/luau/new-se-go-to-class-with-warmup/CompileOptions/typingLetterByLetter"],
  },
]

const mainCharts = combineCharts(mainChartsDeclaration)
const seNewVsOldCharts = combineCharts(seNewVsOldChartsDeclaration)
const charts = combineCharts([...mainChartsDeclaration, ...seNewVsOldChartsDeclaration])
</script>
