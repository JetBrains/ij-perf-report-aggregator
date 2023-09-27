<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_jbr_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
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
    labels: ["Search Everywhere EditorImpl"],
    measures: ["searchEverywhere"],
    projects: [
      "intellij-with-custom-specific-jbr-1048/go-to-class/EditorImpl",
      "intellij-with-custom-specific-jbr-1050/go-to-class/EditorImpl",
      "intellij-with-custom-specific-jbr-1051/go-to-class/EditorImpl",
      "intellij-with-custom-specific-jbr-1052/go-to-class/EditorImpl",
      "intellij-with-custom-specific-jbr-1053/go-to-class/EditorImpl",
    ],
  },
  {
    labels: ["Search Everywhere Kotlin"],
    measures: ["searchEverywhere"],
    projects: [
      "intellij-with-custom-specific-jbr-1048/go-to-class/Kotlin",
      "intellij-with-custom-specific-jbr-1050/go-to-class/Kotlin",
      "intellij-with-custom-specific-jbr-1051/go-to-class/Kotlin",
      "intellij-with-custom-specific-jbr-1052/go-to-class/Kotlin",
      "intellij-with-custom-specific-jbr-1053/go-to-class/Kotlin",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
