<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_performance-k2_dashboard_dev"
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
        :aliases="chart.aliases"
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
    labels: ["FindUsages PsiManager_getInstance"],
    measures: ["findUsages"],
    projects: ["intellij_commit-changedDefault/findUsages/PsiManager_getInstance_firstCall", "intellij_commit/findUsages/PsiManager_getInstance_firstCall"],
    aliases: ["findUsages-k2", "findUsages-k1"],
  },
  {
    labels: ["Local inspections .kt Kotlin Serialization"],
    measures: ["localInspections"],
    projects: ["kotlin-changedDefault/localInspection", "kotlin/localInspection"],
    aliases: ["localInspections-k2", "localInspections-k1"],
  },
  {
    labels: ["Completion .java IntelliJ"],
    measures: ["completion"],
    projects: ["intellij_commit-changedDefault/completion/java_file", "intellij_commit/completion/java_file"],
    aliases: ["java-completion-k2", "java-completion-k1"],
  },
  {
    labels: ["Search Everywhere Go to All"],
    measures: ["searchEverywhere"],
    projects: ["community-changedDefault/go-to-all/Editor/typingLetterByLetter", "community/go-to-all/Editor/typingLetterByLetter"],
    aliases: ["SE-go-to-all-k2", "SE-go-to-all-k1"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
