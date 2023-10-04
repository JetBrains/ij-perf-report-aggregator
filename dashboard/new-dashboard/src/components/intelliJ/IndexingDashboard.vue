<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="idea"
    persistent-id="idea_indexing_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#JAVA'"
          :title="'Indexing Java (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#Kotlin'"
          :title="'Indexing Kotlin (kB/s)'"
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
    labels: ["Indexing (Big projects)", "Scanning (Big projects)", "Number of indexed files (Big projects)"],
    measures: [["indexingTimeWithoutPauses", "indexing"], ["scanningTimeWithoutPauses", "scanning"], "numberOfIndexedFiles"],
    projects: ["community/indexing", "intellij_sources/indexing", "space/indexing"],
  },
  {
    labels: ["Reopening Scanning (Big projects)", "Reopening Number of indexed files (Big projects)"],
    measures: ["scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: ["community/scanning", "intellij_sources/scanning", "space/scanning"],
  },
  {
    labels: [
      "Indexing with the new record storages (IntelliJ project)",
      "Scanning with the new record storages (IntelliJ project)",
      "Number of indexed files with the new record storages (IntelliJ project)",
    ],
    measures: [["indexingTimeWithoutPauses", "indexing"], ["scanningTimeWithoutPauses", "scanning"], "numberOfIndexedFiles"],
    projects: [
      "vfs-record-storage/in-memory-intellij_sources/indexing",
      "vfs-record-storage/over-mmapped-file-intellij_sources/indexing",
      "vfs-record-storage/attributes-over-lock-free-file-page-cache-intellij_sources/indexing",
      "vfs-record-storage/attributes-over-mmapped-file-intellij_sources/indexing",
      "vfs-record-storage/content-hashes-and-attributes-over-mmapped-file-intellij_sources/indexing",
      "vfs-record-storage/content-and-attributes-over-lock-free-file-page-cache-intellij_sources/indexing",
    ],
  },
  {
    labels: ["Indexing", "Scanning", "Number of indexed files"],
    measures: [["indexingTimeWithoutPauses", "indexing"], ["scanningTimeWithoutPauses", "scanning"], "numberOfIndexedFiles"],
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
    labels: ["Indexing", "Scanning", "Number of indexed files"],
    measures: [["indexingTimeWithoutPauses", "indexing"], ["scanningTimeWithoutPauses", "scanning"], "numberOfIndexedFiles"],
    projects: ["keycloak_release_20/indexing", "toolbox_enterprise/indexing", "train-ticket/indexing"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
