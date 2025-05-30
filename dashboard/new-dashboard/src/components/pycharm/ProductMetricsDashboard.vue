<template>
  <DashboardPage
    db-name="perfintDev"
    table="pycharm"
    persistent-id="pycharm_product_dashboard_dev"
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
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["django/indexing", "empty/indexing", "flask/indexing", "keras/indexing", "mypy/indexing"],
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: ["firstCodeAnalysis"],
    projects: [
      "django/findUsages/ForeignKey",
      "django/findUsages/Form",
      "django/findUsages/Model",
      "flask/findUsages/Flask",
      "flask/findUsages/request",
      "keras/findUsages/Sequential",
      "mypy/findUsages/Errors",
    ],
  },
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: ["edx-platform (Django)/completion/model", "edx-platform (Django)/completion/view"],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: [
      "mypy/go-to-all-with-warmup/class/typingLetterByLetter",
      "empty/go-to-all-with-warmup/class/typingLetterByLetter",
      "flask/go-to-all-with-warmup/class/typingLetterByLetter",
      "keras/go-to-all-with-warmup/class/typingLetterByLetter",
    ],
  },
  {
    labels: ["TypingCodeAnalysis"],
    measures: ["typingCodeAnalyzing"],
    projects: [],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["django/inspection", "flask/inspection", "keras/inspection", "mypy/inspection"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
