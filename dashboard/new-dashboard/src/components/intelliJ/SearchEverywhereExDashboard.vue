<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_java_dashboard_devserver"
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

function seProjects(tab: string): string[] {
  const patterns = ["Kotlin", "Editor", "Runtime"]
  return [
    ...patterns.map((pattern) => `community/${tab}/${pattern}/typingLetterByLetter`),
    `community/${tab}-finished-embeddings/Runtime/typingLetterByLetter`,
    `java/${tab}/Runtime/typingLetterByLetter`,

    ...patterns.map((pattern) => `intellij_commit/${tab}/${pattern}/typingLetterByLetter`),
    `intellij_commit/${tab}-finished-embeddings/Runtime/typingLetterByLetter`,

    ...patterns.map((pattern) => `intellij_commit/new-se-${tab}/${pattern}/typingLetterByLetter`),
  ]
}

function seCharts(tabName: string, projectPrefix: string): ChartDefinition[] {
  const seMeasures: string[] = [
    "searchEverywhere_first_elements_added",
    "searchEverywhere_elements_added_5",
    "searchEverywhere_elements_added_10",
    "searchEverywhere_elements_added_15",
    "searchEverywhere",
  ]

  return ["Cold", "Warm"].map((temp) => ({
    labels: [
      `${temp} SE Elements Added ${tabName} (slow typing) - 1`,
      `${temp} SE Elements Added ${tabName} (slow typing) - 5`,
      `${temp} SE Elements Added ${tabName} (slow typing) - 10`,
      `${temp} SE Elements Added ${tabName} (slow typing) - 15`,
      `${temp} SE Elements Added ${tabName} (slow typing) - ALL`,
    ],
    measures: seMeasures,
    projects: seProjects(`go-to-${projectPrefix}${temp === "Warm" ? "-with-warmup" : ""}`),
  }))
}

const chartsDeclaration: ChartDefinition[] = [
  ...seCharts("Action", "action"),
  ...seCharts("Class", "class"),
  ...seCharts("File", "file"),
  ...seCharts("All", "all"),
  ...seCharts("Symbol", "symbol"),
  ...seCharts("Text", "text"),
  ...seCharts("Lucene Files", "lucene-files"),
]

const charts = combineCharts(chartsDeclaration)
</script>
