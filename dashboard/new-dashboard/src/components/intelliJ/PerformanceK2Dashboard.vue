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
    projects: ["intellij_commit/findUsages/PsiManager_getInstance_firstCall"],
  },
  {
    labels: ["Local inspections .kt Kotlin Serialization"],
    measures: ["localInspections"],
    projects: ["kotlin/localInspection"],
  },
  {
    labels: ["Completion .java IntelliJ"],
    measures: ["completion"],
    projects: ["intellij_commit/completion/java_file"],
  },
  {
    labels: ["Search Everywhere Go to All"],
    measures: ["searchEverywhere"],
    projects: ["community/go-to-all/Editor/typingLetterByLetter"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
