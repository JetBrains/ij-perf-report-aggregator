<template>
  <DashboardPage
    db-name="perfint"
    table="datagrip"
    persistent-id="datagrip_product_dashboard"
    initial-machine="linux-blade-hetzner"
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
    measures: ["indexingTimeWithoutPauses"],
    projects: [],
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: ["firstCodeAnalysis"],
    projects: [],
  },
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: [],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: [],
  },
  {
    labels: ["TypingCodeAnalysis"],
    measures: ["typingCodeAnalyzing"],
    projects: [],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: [],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
