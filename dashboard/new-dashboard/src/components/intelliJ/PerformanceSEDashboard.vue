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
    labels: ["Search Everywhere Class (slow typing)", "SE Dialog Shown Class (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-class/Kotlin/typingLetterByLetter", "community/go-to-class/Editor/typingLetterByLetter", "community/go-to-class/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere File (slow typing)", "SE Dialog Shown File (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-file/Editor/typingLetterByLetter", "community/go-to-file/Kotlin/typingLetterByLetter", "community/go-to-file/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere All (slow typing)", "SE Dialog Shown All (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-all/Editor/typingLetterByLetter", "community/go-to-all/Kotlin/typingLetterByLetter", "community/go-to-all/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere Symbol (slow typing)", "SE Dialog Shown Symbol (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-symbol/Editor/typingLetterByLetter", "community/go-to-symbol/Kotlin/typingLetterByLetter", "community/go-to-symbol/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere Text (slow typing)", "SE Dialog Shown Text (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-text/Editor/typingLetterByLetter", "community/go-to-text/Kotlin/typingLetterByLetter", "community/go-to-text/Runtime/typingLetterByLetter"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
