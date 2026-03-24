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

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Cold Search Everywhere Action (slow typing)", "Cold SE First Element Added Action (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-action"),
  },
  {
    labels: ["Warm Search Everywhere Action (slow typing)", "Warm SE First Element Added Action (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-action-with-warmup"),
  },
  {
    labels: ["Cold Search Everywhere Class (slow typing)", "Cold SE First Element Added Class (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-class"),
  },
  {
    labels: ["Warm Search Everywhere Class (slow typing)", "Warm SE First Element Added Class (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-class-with-warmup"),
  },
  {
    labels: ["Cold Search Everywhere File (slow typing)", "Cold SE First Element Added File(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-file"),
  },
  {
    labels: ["Warm Search Everywhere File (slow typing)", "Warm SE First Element Added File(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-file-with-warmup"),
  },
  {
    labels: ["Cold Search Everywhere All (slow typing)", "Cold SE First Element Added All(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-all"),
  },
  {
    labels: ["Warm Search Everywhere All (slow typing)", "Warm SE First Element Added All(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-all-with-warmup"),
  },
  {
    labels: ["Cold Search Everywhere Symbol (slow typing)", "Cold SE First Element Symbol (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-symbol"),
  },
  {
    labels: ["Warm Search Everywhere Symbol (slow typing)", "Warm SE First Element Symbol (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-symbol-with-warmup"),
  },
  {
    labels: ["Cold Search Everywhere Text (slow typing)", "Cold SE First Element Text (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-text"),
  },
  {
    labels: ["Warm Search Everywhere Text (slow typing)", "Warm SE First Element Text (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-text-with-warmup"),
  },
  {
    labels: ["EXPERIMENT: Cold Search Everywhere Lucene Files (slow typing)", "EXPERIMENT: Cold SE First Element Lucene Files (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-lucene-files"),
  },
  {
    labels: ["EXPERIMENT: Warm Search Everywhere Lucene Files (slow typing)", "EXPERIMENT: Warm SE First Element Lucene Files (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: seProjects("go-to-lucene-files-with-warmup"),
  },
  {
    labels: ["performance.ui.lagging", "performance.ui.latency", "performance.popup.latency"],
    measures: ["ui.lagging#max_value", "ui.latency#max_value", "popup.latency#max_value"],
    projects: [
      "popups-performance-test/test-popups",
      "typingInJavaFile_16Threads/typing",
      "typingInJavaFile_4Threads/typing",
      "typingInKotlinFile_16Threads/typing",
      "typingInKotlinFile_4Threads/typing",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
