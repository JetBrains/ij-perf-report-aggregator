<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_search_everywhere_dashboard_dev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsWithClientChart
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
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Search Everywhere Action (slow typing)", "SE First Element Added Action (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-action/Kotlin/typingLetterByLetter",
      "community/go-to-action/Editor/typingLetterByLetter",
      "community/go-to-action/Runtime/typingLetterByLetter",
      "community/go-to-action-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-action/Runtime/typingLetterByLetter",
      "community/go-to-action-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-action-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-action-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-action-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-action-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-action/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-action/Editor/typingLetterByLetter",
      "intellij_commit/go-to-action/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-action-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-action-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-action-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-action-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Search Everywhere Class (slow typing)", "SE First Element Added Class (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-class/Kotlin/typingLetterByLetter",
      "community/go-to-class/Editor/typingLetterByLetter",
      "community/go-to-class/Runtime/typingLetterByLetter",
      "community/go-to-class-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-class/Runtime/typingLetterByLetter",
      "community/go-to-class-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-class-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-class-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-class-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-class-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-class/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-class/Editor/typingLetterByLetter",
      "intellij_commit/go-to-class/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-class-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-class-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-class-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-class-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Search Everywhere File (slow typing)", "SE First Element Added File(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-file/Editor/typingLetterByLetter",
      "community/go-to-file/Kotlin/typingLetterByLetter",
      "community/go-to-file/Runtime/typingLetterByLetter",
      "community/go-to-file-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-file/Runtime/typingLetterByLetter",
      "community/go-to-file-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-file-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-file-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-file-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-file-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-file/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-file/Editor/typingLetterByLetter",
      "intellij_commit/go-to-file/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-file-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-file-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-file-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-file-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Search Everywhere All (slow typing)", "SE First Element Added All(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-all/Editor/typingLetterByLetter",
      "community/go-to-all/Kotlin/typingLetterByLetter",
      "community/go-to-all/Runtime/typingLetterByLetter",
      "community/go-to-all-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-all/Runtime/typingLetterByLetter",
      "community/go-to-all-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-all-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-all-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-all-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-all-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-all/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-all/Editor/typingLetterByLetter",
      "intellij_commit/go-to-all/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-all-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-all-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-all-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-all-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Search Everywhere Symbol (slow typing)", "SE First Element Symbol (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-symbol/Editor/typingLetterByLetter",
      "community/go-to-symbol/Kotlin/typingLetterByLetter",
      "community/go-to-symbol/Runtime/typingLetterByLetter",
      "community/go-to-symbol-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-symbol/Runtime/typingLetterByLetter",
      "community/go-to-symbol-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-symbol-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-symbol-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-symbol/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-symbol/Editor/typingLetterByLetter",
      "intellij_commit/go-to-symbol/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-symbol-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-symbol-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-symbol-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Search Everywhere Text (slow typing)", "SE First Element Text (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-text/Editor/typingLetterByLetter",
      "community/go-to-text/Kotlin/typingLetterByLetter",
      "community/go-to-text/Runtime/typingLetterByLetter",
      "community/go-to-text-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-text/Runtime/typingLetterByLetter",
      "community/go-to-text-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-text-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-text-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-text-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-text-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-text/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-text/Editor/typingLetterByLetter",
      "intellij_commit/go-to-text/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-text-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-text-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-text-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-text-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
