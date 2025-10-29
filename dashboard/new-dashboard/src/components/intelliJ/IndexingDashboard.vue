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
import { INDEXING_PROJECTS } from "./indexingProjects"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["VFS Refresh"],
    measures: ["vfs_initial_refresh"],
    projects: INDEXING_PROJECTS.vfsRefresh.default,
    aliases: ["default", "1 thread", "git status"],
  },
  {
    labels: ["VFS Refresh after Mass Changes"],
    measures: [["vfsRefreshAfterMassCreate", "vfsRefreshAfterMassModify", "vfsRefreshAfterMassDelete"]],
    projects: INDEXING_PROJECTS.vfsRefresh.afterMassChanges,
    aliases: ["txt", "java", "kotlin"],
  },
  {
    labels: ["Indexing (big projects)"],
    measures: [["indexingTimeWithoutPauses", "fus_dumb_indexing_time"]],
    projects: INDEXING_PROJECTS.indexing.big,
  },
  {
    labels: ["Indexing (middle projects)"],
    measures: [["indexingTimeWithoutPauses", "fus_dumb_indexing_time"]],
    projects: INDEXING_PROJECTS.indexing.medium,
  },
  {
    labels: ["Indexing (small projects)"],
    measures: [["indexingTimeWithoutPauses", "fus_dumb_indexing_time"]],
    projects: INDEXING_PROJECTS.indexing.small,
  },
  {
    labels: ["First Scanning (big projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.firstScanning.big,
  },
  {
    labels: ["First Scanning (middle projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.firstScanning.medium,
  },
  {
    labels: ["First Scanning (small projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.firstScanning.small,
  },
  {
    labels: ["Second Scanning (big projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.secondScanning.big,
  },
  {
    labels: ["Second Scanning (middle projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.secondScanning.medium,
  },
  {
    labels: ["Second Scanning (small projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.secondScanning.small,
  },
  {
    labels: ["Third Scanning (big projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.thirdScanning.big,
  },
  {
    labels: ["Third Scanning (middle projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.thirdScanning.medium,
  },
  {
    labels: ["Third Scanning (small projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: INDEXING_PROJECTS.thirdScanning.small,
  },
  {
    labels: ["Number of indexed files (big projects)"],
    measures: ["numberOfIndexedFiles"],
    projects: INDEXING_PROJECTS.numberOfIndexedFiles.big,
  },
  {
    labels: ["Number of indexed files (middle projects)"],
    measures: ["numberOfIndexedFiles"],
    projects: INDEXING_PROJECTS.numberOfIndexedFiles.medium,
  },
  {
    labels: ["Number of indexed files (small projects)"],
    measures: ["numberOfIndexedFiles"],
    projects: INDEXING_PROJECTS.numberOfIndexedFiles.small,
  },
  {
    labels: ["Indexing (after checkout)", "Scanning (after checkout)", "Number of indexed files (after checkout)"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: INDEXING_PROJECTS.afterCheckout,
  },
  {
    labels: ["Indexing", "Scanning", "Number of indexed files", "Number of indexed files with writing index value", "Number of indexed files with nothing to write"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles", "numberOfIndexedFilesWritingIndexValue", "numberOfIndexedFilesWithNothingToWrite"],
    projects: INDEXING_PROJECTS.indexStorage,
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
