<template>
  <DashboardPage
    db-name="perfintDev"
    table="rust"
    persistent-id="rust_product_dashboard"
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
import { rustCompletionCases, rustGlobalInspectionProjects, rustLocalInspectionCases } from "./RustTestCases"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: ["indexingTimeWithoutPauses"],
    projects: rustGlobalInspectionProjects.map((project) => `${project}/indexing`),
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: ["firstCodeAnalysis"],
    projects: rustLocalInspectionCases,
  },
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: rustCompletionCases,
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: ["searchEverywhere/cargo/go-to-all-with-warmup/Display/typingLetterByLetter"],
  },
  {
    labels: ["TypingCodeAnalysis"],
    measures: ["typingCodeAnalyzing"],
    projects: rustLocalInspectionCases.concat(rustLocalInspectionCases.map((testCase) => `${testCase}-top-level-typing`)),
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: rustGlobalInspectionProjects.map((project) => `global-inspection/${project}-inspection`),
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
