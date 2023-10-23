<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_search_everywhere_dashboard"
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
    labels: ["Search Everywhere (insert whole word)", "SE Dialog Shown (insert whole word)", "SE Items Loaded (insert whole word)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
    projects: [
      "community/go-to-action/SharedIndex",
      "community/go-to-class/EditorImpl",
      "community/go-to-class/SharedIndex",
      "community/go-to-file/EditorImpl",
      "community/go-to-file/SharedIndex",
      "community/go-to-action/SharedIndex/insertingTheWholeWord",
      "community/go-to-class/EditorImpl/insertingTheWholeWord",
      "community/go-to-class/SharedIndex/insertingTheWholeWord",
      "community/go-to-file/EditorImpl/insertingTheWholeWord",
      "community/go-to-file/SharedIndex/insertingTheWholeWord",
    ],
  },
  {
    labels: ["Search Everywhere Action (slow typing)", "SE Dialog Shown Action (slow typing)", "SE Items Loaded Action (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
    projects: [
      "community/go-to-action/Runtime",
      "community/go-to-action/Runtime/typingLetterByLetter",
      "community/go-to-action/Editor/typingLetterByLetter",
      "community/go-to-action/Kotlin/typingLetterByLetter",
    ],
  },
  {
    labels: ["Search Everywhere Class (slow typing)", "SE Dialog Shown Class (slow typing)", "SE Items Loaded Class (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
    projects: ["community/go-to-class/Kotlin", "community/go-to-class/Kotlin/typingLetterByLetter", "community/go-to-class/Editor/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere File (slow typing)", "SE Dialog Shown File (slow typing)", "SE Items Loaded File (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
    projects: [
      "community/go-to-file/properties",
      "community/go-to-file/properties/typingLetterByLetter",
      "community/go-to-file/Editor/typingLetterByLetter",
      "community/go-to-file/Kotlin/typingLetterByLetter",
    ],
  },
  {
    labels: ["Search Everywhere All (slow typing)", "SE Dialog Shown All (slow typing)", "SE Items Loaded All (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
    projects: ["community/go-to-all/Editor/typingLetterByLetter", "community/go-to-all/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere Symbol (slow typing)", "SE Dialog Shown Symbol (slow typing)", "SE Items Loaded Symbol (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
    projects: ["community/go-to-symbol/Editor/typingLetterByLetter", "community/go-to-symbol/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere Text (slow typing)", "SE Dialog Shown Text (slow typing)", "SE Items Loaded Text (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
    projects: ["community/go-to-text/Editor/typingLetterByLetter", "community/go-to-text/Kotlin/typingLetterByLetter"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
