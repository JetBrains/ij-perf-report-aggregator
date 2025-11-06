<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_lagging_latency_dashboard"
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
    labels: ["ui.Lagging during indexing - average"],
    measures: ["ui.lagging#average"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/curl/indexing", "radler/big_project_50k_10k/indexing"],
    aliases: ["LLVM", "OpenCV", "cURL", "Big Project"],
  },
  {
    labels: ["ui.Lagging during indexing - max"],
    measures: ["ui.lagging#max"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/curl/indexing", "radler/big_project_50k_10k/indexing"],
    aliases: ["LLVM", "OpenCV", "cURL", "Big Project"],
  },
  {
    labels: ["ui.Lagging during indexing - count"],
    measures: ["ui.lagging#count"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/curl/indexing", "radler/big_project_50k_10k/indexing"],
    aliases: ["LLVM", "OpenCV", "cURL", "Big Project"],
  },
  {
    labels: ["ui.Lagging during indexing - sum"],
    measures: ["ui.lagging#sum"],
    projects: ["radler/llvm/indexing", "radler/opencv/indexing", "radler/curl/indexing", "radler/big_project_50k_10k/indexing"],
    aliases: ["LLVM", "OpenCV", "cURL", "Big Project"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
