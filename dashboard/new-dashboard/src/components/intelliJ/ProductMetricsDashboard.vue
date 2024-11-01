<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_product_dashboard_dev"
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
    measures: ["indexingTimeWithoutPauses"],
    projects: ["community/indexing", "intellij_commit/indexing", "kotlin/indexing"],
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: ["firstCodeAnalysis"],
    projects: ["intellij_commit/localInspection/java_file", "intellij_commit/localInspection/kotlin_file", "kotlin/localInspection", "kotlin_coroutines/localInspection"],
  },
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: ["intellij_commit/completion/kotlin_file", "intellij_commit/completion/java_file"],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: ["community/go-to-all/Editor/typingLetterByLetter", "community/go-to-all-with-warmup/Editor/typingLetterByLetter"],
  },
  {
    labels: ["TypingCodeAnalysis"],
    measures: ["typingCodeAnalyzing"],
    projects: ["toolbox_enterprise/ultimateCase/SecurityTests", "keycloak_release_20/ultimateCase/JpaUserProvider", "train-ticket/ultimateCase/InsidePaymentServiceImpl"],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["kotlin_coroutines/inspection", "spring_boot_maven/inspection"],
  },
  {
    labels: ["Highlighting - remove symbol"],
    measures: ["typing_EditorBackSpace_duration"],
    projects: ["intellij_commit/editor-highlighting"],
  },
  {
    labels: ["Highlighting - type symbol"],
    measures: ["typing_}_duration"],
    projects: ["intellij_commit/editor-highlighting"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
