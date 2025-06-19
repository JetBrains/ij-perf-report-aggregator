<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="lagging_latency_dashboard_dev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsWithClientChart
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
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["performance.ui.Lagging", "performance.ui.latency", "editor.typing.latency.max", "editor.typing.latency.90", "performance.popup.latency"],
    measures: ["ui.lagging#max_value", "ui.latency#max_value", "editor.typing.latency#max", "editor.typing.latency#90", "popup.latency#max_value"],
    projects: [
      "popups-performance-test/test-popups",
      "typingInJavaFile_16Threads/typing",
      "typingInJavaFile_4Threads/typing",
      "typingInKotlinFile_16Threads/typing",
      "typingInKotlinFile_4Threads/typing",
    ],
  },
  {
    labels: ["Typing during indexing (average awt delay)", "Typing during indexing (max awt delay)"],
    measures: ["test#average_awt_delay", "test#max_awt_delay"],
    projects: ["typingInJavaFile_16Threads/typing", "typingInJavaFile_4Threads/typing", "typingInKotlinFile_16Threads/typing", "typingInKotlinFile_4Threads/typing"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
