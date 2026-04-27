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

function seProjects(tab: string, patterns: string[]): string[] {
  return [
    ...patterns.map((pattern) => `community/${tab}/${pattern}/typingLetterByLetter`),
    `community/${tab}-finished-embeddings/Runtime/typingLetterByLetter`,
    `java/${tab}/Runtime/typingLetterByLetter`,

    ...patterns.map((pattern) => `intellij_commit/${tab}/${pattern}/typingLetterByLetter`),
    `intellij_commit/${tab}-finished-embeddings/Runtime/typingLetterByLetter`,

    ...patterns.map((pattern) => `intellij_commit/new-se-${tab}/${pattern}/typingLetterByLetter`),
  ]
}

function seChartsCustom(tabName: string, projectPrefix: string, patterns: string[], customLabel: string): ChartDefinition[] {
  const seMeasures: string[] = ["searchEverywhere_first_elements_added", "searchEverywhere_elements_added_10", "searchEverywhere"]

  const customLabelFixed = customLabel === "" ? "" : ` - ${customLabel}`

  return ["Cold", "Warm"].map((temp) => ({
    labels: [
      `${tabName}${customLabelFixed} - 1 element - ${temp}`,
      `${tabName}${customLabelFixed} - 10 elements - ${temp}`,
      `${tabName}${customLabelFixed} - All elements - ${temp}`,
    ],
    measures: seMeasures,
    projects: seProjects(`go-to-${projectPrefix}${temp === "Warm" ? "-with-warmup" : ""}${customLabel === "" ? "" : `-${customLabel}`}`, patterns),
  }))
}

function seCharts(tabName: string, projectPrefix: string): ChartDefinition[] {
  const patterns = ["Kotlin", "Editor", "Runtime"]
  return seChartsCustom(tabName, projectPrefix, patterns, "")
}

// Intentional misspellings for fuzzy search testing
const fuzzyFilesPatterns: string[] = ["Kotlin", "Editor", "Runtime", "Kutlin", "Edetor", "Rantime"]

const chartsDeclaration: ChartDefinition[] = [
  ...seCharts("All", "all"),
  ...seCharts("Action", "action"),
  ...seCharts("Text", "text"),
  ...seCharts("Class", "class"),
  ...seCharts("Symbol", "symbol"),
  ...seCharts("File", "file"),
  ...seCharts("Lucene Files", "lucene-files"),
  {
    labels: ["Lucene Files - Index Creation Time - Cold", "Lucene Files - Index Size (KB) - Cold"],
    measures: ["searchEverywhereLuceneFilesIndexAll", "searchEverywhereLuceneIndexSize"],
    projects: seProjects("go-to-lucene-files", ["Kotlin", "Editor", "Runtime"]),
  },
  {
    labels: ["Lucene Files - Index Creation Time - Warm", "Lucene Files - Index Size (KB) - Warm"],
    measures: ["searchEverywhereLuceneFilesIndexAll", "searchEverywhereLuceneIndexSize"],
    projects: seProjects("go-to-lucene-files-with-warmup", ["Kotlin", "Editor", "Runtime"]),
  },
  ...seChartsCustom("File", "file", fuzzyFilesPatterns, "fuzzy"),
  {
    labels: ["Warm Search Everywhere Insert", "Warm SE First Element Insert"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "intellij_commit/go-to-all-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
      "intellij_commit/go-to-file-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
      "intellij_commit/go-to-class-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
      "intellij_commit/go-to-symbol-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
      "intellij_commit/go-to-action-with-warmup/CollectLogsAndDiagnosticData/insertingTheWholeWord",
      "intellij_commit/go-to-text-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",

      "intellij_commit/new-se-go-to-all-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
      "intellij_commit/new-se-go-to-file-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
      "intellij_commit/new-se-go-to-class-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
      "intellij_commit/new-se-go-to-symbol-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
      "intellij_commit/new-se-go-to-action-with-warmup/CollectLogsAndDiagnosticData/insertingTheWholeWord",
      "intellij_commit/new-se-go-to-text-with-warmup/AppServerIntegrationsManagerImpl/insertingTheWholeWord",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
