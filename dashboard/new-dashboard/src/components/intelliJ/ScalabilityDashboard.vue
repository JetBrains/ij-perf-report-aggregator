<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_scalability_dashboard"
    initial-machine="Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
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
  </DashboardPage>>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const processorCounts = [1, 2, 4, 8, 16, 32, 64]
const chartsDeclaration: ChartDefinition[] = [{
  labels: ["Indexing (Empty Project)", "Scanning(Empty Project)", "Number of indexed files(Empty Project)"],
  measures: ["indexing", "scanning", "numberOfIndexedFiles"],
  projects: processorCounts.map(it => `empty_project/indexing/processorCount_${it}`)
},
  {
    labels: ["Indexing (IntelliJ Sources)", "Scanning(IntelliJ Sources)", "Number of indexed files(IntelliJ Sources)"],
    measures: ["indexing", "scanning", "numberOfIndexedFiles"],
    projects: processorCounts.map(it => `intellij_sources/indexing/processorCount_${it}`)
  },
  {
    labels: ["Indexing (Kotlin Coroutines)", "Scanning(Kotlin Coroutines)", "Number of indexed files(Kotlin Coroutines)"],
    measures: ["indexing", "scanning", "numberOfIndexedFiles"],
    projects: processorCounts.map(it => `kotlin_coroutines/indexing/processorCount_${it}`)
  },
  {
    labels: ["Indexing (Spring Boot)", "Scanning(Spring Boot)", "Number of indexed files(Spring Boot)"],
    measures: ["indexing", "scanning", "numberOfIndexedFiles"],
    projects: processorCounts.map(it => `spring_boot/indexing/processorCount_${it}`)
  }]

const charts = combineCharts(chartsDeclaration)
</script>