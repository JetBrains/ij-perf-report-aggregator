<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="lspDashboard"
    initial-machine="linux-blade-hetzner"
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
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["diagnostics (dummy files)"],
    measures: ["attempt.mean.ms"],
    projects: ["com.jetbrains.ls.lsp.performanceTest.DiagnosticPerformanceTest$huge$1.huge", "com.jetbrains.ls.lsp.performanceTest.DiagnosticPerformanceTest$tiny$1.tiny"],
  },
  {
    labels: ["diagnostics (real projects)"],
    measures: ["attempt.mean.ms"],
    projects: [
      "com.jetbrains.ls.lsp.performanceTest.PetClinicJavaPerformanceTest$diagnostics$1$1$2.diagnostics",
      "com.jetbrains.ls.lsp.performanceTest.PetClinicKotlinPerformanceTest$diagnostics$1$1$2.diagnostics",
    ],
  },
  {
    labels: ["completion"],
    measures: ["attempt.mean.ms"],
    projects: [
      "com.jetbrains.ls.lsp.performanceTest.PetClinicKotlinPerformanceTest$completion$1$1$2.completion",
      "com.jetbrains.ls.lsp.performanceTest.PetClinicJavaPerformanceTest$completion$1$1$2.completion",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
