<template>
  <DashboardPage
    db-name="perfint"
    table="webstorm"
    persistent-id="webstorm_product_dashboard_old"
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
    projects: ["ring-ui/indexing", "axios/indexing", "dxos/indexing"],
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: ["firstCodeAnalysis"],
    projects: ["allure-js/localInspection/JasmineAllureReporter.ts", "axios/localInspection/utils.js", "material-ui-react-admin/localInspection/PostEdit.tsx"],
  },
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: ["axios/completion/functions", "eslint-plugin-jest/completion/types", "pancake-frontend/completion/component", "vue3-admin-vite/completion/component"],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: ["ring-ui/go-to-class-with-warmup/Button/typingLetterByLetter"],
  },
  {
    labels: ["TypingCodeAnalysis"],
    measures: ["typingCodeAnalyzing"],
    projects: ["axios/typing/axios", "eslint-plugin-jest/typing/eslintPluginJest"],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["axios/inspection/", "gitlab/inspection/appDir", "ring-ui/inspection/src"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
