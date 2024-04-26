<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_indexing_dashboard_dev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#JAVA'"
          :title="'Indexing Speed Java (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#Kotlin'"
          :title="'Indexing Speed Kotlin (kB/s)'"
          :chart-color="'#9B51E0'"
          :value-unit="'counter'"
        />
      </div>
    </section>
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
import AggregationChart from "../charts/AggregationChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: [
      "Indexing (Big projects)",
      "Initial Scanning (Big projects)",
      "Initial Number of indexed files (Big projects)",
      "Number of indexed files with writing index value (Big projects)",
      "Number of indexed files with nothing to write (Big projects)",
    ],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles", "numberOfIndexedFilesWritingIndexValue", "numberOfIndexedFilesWithNothingToWrite"],
    projects: ["community/indexing", "intellij_commit/indexing"],
  },
  {
    labels: ["Processing Time Java", "Processing Time Kotlin"],
    measures: ["processingTime#JAVA", "processingTime#Kotlin"],
    projects: ["community/indexing", "intellij_commit/indexing"],
  },
  {
    labels: ["Second Scanning (Big projects)", "Second Number of indexed files (Big projects)"],
    measures: ["scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: ["community/second-scanning", "intellij_commit/second-scanning"],
  },
  {
    labels: ["Third Scanning (Big projects)", "Third Number of indexed files (Big projects)"],
    measures: ["scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: ["community/third-scanning", "intellij_commit/third-scanning"],
  },
  {
    labels: ["Indexing", "Scanning", "Number of indexed files", "Number of indexed files with writing index value", "Number of indexed files with nothing to write"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles", "numberOfIndexedFilesWritingIndexValue", "numberOfIndexedFilesWithNothingToWrite"],
    projects: [
      "empty_project/indexing",
      "grails/indexing",
      "java/indexing",
      "kotlin/indexing",
      "kotlin_coroutines/indexing",
      "spring_boot/indexing",
      "spring_boot_maven/indexing",
      "kotlin_petclinic/indexing",
    ],
  },
  {
    labels: ["Processing Time Java", "Processing Time Kotlin"],
    measures: ["processingTime#JAVA", "processingTime#Kotlin"],
    projects: [
      "empty_project/indexing",
      "grails/indexing",
      "java/indexing",
      "kotlin/indexing",
      "kotlin_coroutines/indexing",
      "spring_boot/indexing",
      "spring_boot_maven/indexing",
      "kotlin_petclinic/indexing",
    ],
  },
  {
    labels: ["Indexing", "Scanning", "Number of indexed files", "Number of indexed files with writing index value", "Number of indexed files with nothing to write"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles", "numberOfIndexedFilesWritingIndexValue", "numberOfIndexedFilesWithNothingToWrite"],
    projects: ["keycloak_release_20/indexing", "toolbox_enterprise/indexing", "train-ticket/indexing"],
  },
  {
    labels: ["Processing Time Java", "Processing Time Kotlin"],
    measures: ["processingTime#JAVA", "processingTime#Kotlin"],
    projects: ["keycloak_release_20/indexing", "toolbox_enterprise/indexing", "train-ticket/indexing"],
  },
  {
    labels: ["Second Scanning", "Second Number of indexed files"],
    measures: ["scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: ["keycloak_release_20/second-scanning", "toolbox_enterprise/second-scanning", "train-ticket/second-scanning"],
  },
  {
    labels: ["Third Scanning", "Third Number of indexed files"],
    measures: ["scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: ["keycloak_release_20/third-scanning", "toolbox_enterprise/third-scanning", "train-ticket/third-scanning"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
