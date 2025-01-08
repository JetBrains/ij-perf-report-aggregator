<template>
  <DashboardPage
    db-name="perfintDev"
    table="phpstorm"
    persistent-id="phpstorm_product_dashboard"
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
    projects: ["laravel-io/indexing", "oro/indexing", "aggregateStitcher/indexing"],
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: ["firstCodeAnalysis"],
    projects: ["laravel-io/localInspection/HasAuthor", "laravel-io/localInspection/Tag", "mpdf/localInspection"],
  },
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: ["dql/completion", "many_classes/completion/classes"],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: ["magento2/go-to-class-with-warmup/MaAdMUser/typingLetterByLetter", "bitrix/go-to-class-with-warmup/BCCo/typingLetterByLetter"],
  },
  {
    labels: ["TypingCodeAnalysis"],
    measures: ["typingCodeAnalyzing"],
    projects: [],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["laravel-io/inspection", "akaunting/inspection", "aggregateStitcher/inspection"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
