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
    labels: ["VFS Refresh"],
    measures: ["vfs_initial_refresh"],
    projects: ["intellij_commit/vfsRefresh/default", "intellij_commit/vfsRefresh/with-1-thread(s)", "intellij_commit/vfsRefresh/git-status"],
    aliases: ["default", "1 thread", "git status"],
  },
  {
    labels: ["VFS Refresh after Mass Changes"],
    measures: [["vfsRefreshAfterMassCreate", "vfsRefreshAfterMassModify", "vfsRefreshAfterMassDelete"]],
    projects: ["empty_project/vfs-mass-update-txt", "empty_project/vfs-mass-update-java", "empty_project/vfs-mass-update-kt"],
    aliases: ["txt", "java", "kotlin"],
  },
  {
    labels: ["Indexing (big projects)"],
    measures: [["indexingTimeWithoutPauses", "fus_dumb_indexing_time"]],
    projects: [
      "community/indexing",
      "intellij_commit/indexing",
      "community/indexingWithHighlighting",
      "indexing-community/indexingWithHighlighting",
      "indexing-intellij_commit/indexingWithHighlighting",
      "space/indexing",
      "kotlin_coroutines/indexing",
    ],
  },
  {
    labels: ["Indexing (middle projects)"],
    measures: [["indexingTimeWithoutPauses", "fus_dumb_indexing_time"]],
    projects: ["kotlin/indexing", "toolbox_enterprise/indexing", "spring_boot/indexing", "keycloak_release_20/indexing"],
  },
  {
    labels: ["Indexing (small projects)"],
    measures: [["indexingTimeWithoutPauses", "fus_dumb_indexing_time"]],
    projects: ["empty_project/indexing", "grails/indexing", "java/indexing", "spring_boot_maven/indexing", "kotlin_petclinic/indexing", "train-ticket/indexing"],
  },
  {
    labels: ["First Scanning (big projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: ["community/indexing", "intellij_commit/indexing", "space/indexing", "kotlin_coroutines/indexing"],
  },
  {
    labels: ["First Scanning (middle projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: ["kotlin/indexing", "toolbox_enterprise/indexing", "spring_boot/indexing", "keycloak_release_20/indexing"],
  },
  {
    labels: ["First Scanning (small projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: ["empty_project/indexing", "grails/indexing", "java/indexing", "spring_boot_maven/indexing", "kotlin_petclinic/indexing", "train-ticket/indexing"],
  },
  {
    labels: ["Second Scanning (big projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: ["community/second-scanning", "intellij_commit/second-scanning", "space/second-scanning", "kotlin_coroutines/second-scanning"],
  },
  {
    labels: ["Second Scanning (middle projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: ["kotlin/second-scanning", "toolbox_enterprise/second-scanning", "spring_boot/second-scanning", "keycloak_release_20/second-scanning"],
  },
  {
    labels: ["Second Scanning (small projects)"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: [
      "empty_project/second-scanning",
      "grails/second-scanning",
      "java/second-scanning",
      "spring_boot_maven/second-scanning",
      "kotlin_petclinic/second-scanning",
      "train-ticket/second-scanning",
    ],
  },
  {
    labels: ["Number of indexed files (big projects)"],
    measures: ["numberOfIndexedFiles"],
    projects: [
      "community/indexing",
      "intellij_commit/indexing",
      "space/indexing",
      "kotlin_coroutines/indexing",
      "community/second-scanning",
      "intellij_commit/second-scanning",
      "space/second-scanning",
      "kotlin_coroutines/second-scanning",
    ],
  },
  {
    labels: ["Number of indexed files (middle projects)"],
    measures: ["numberOfIndexedFiles"],
    projects: [
      "kotlin/indexing",
      "toolbox_enterprise/indexing",
      "spring_boot/indexing",
      "keycloak_release_20/indexing",
      "kotlin/second-scanning",
      "toolbox_enterprise/second-scanning",
      "spring_boot/second-scanning",
      "keycloak_release_20/second-scanning",
    ],
  },
  {
    labels: ["Number of indexed files (small projects)"],
    measures: ["numberOfIndexedFiles"],
    projects: [
      "empty_project/indexing",
      "grails/indexing",
      "java/indexing",
      "spring_boot_maven/indexing",
      "kotlin_petclinic/indexing",
      "train-ticket/indexing",
      "empty_project/second-scanning",
      "grails/second-scanning",
      "java/second-scanning",
      "spring_boot_maven/second-scanning",
      "kotlin_petclinic/second-scanning",
      "train-ticket/second-scanning",
    ],
  },
  {
    labels: ["Indexing (after checkout)", "Scanning (after checkout)", "Number of indexed files (after checkout)"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: ["intellij_sources/checkout/243", "intellij_commit/checkout/243"],
  },
  {
    labels: ["Indexing", "Scanning", "Number of indexed files", "Number of indexed files with writing index value", "Number of indexed files with nothing to write"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles", "numberOfIndexedFilesWritingIndexValue", "numberOfIndexedFilesWithNothingToWrite"],
    projects: [
      "index-storage/default-2shards-intellij_commit/indexing",
      "index-storage/default-intellij_commit/indexing",
      "index-storage/default-w-coroutines-intellij_commit/indexing",
      "index-storage/default-w-mru-cache-intellij_commit/indexing",
      "index-storage/fake-writer-intellij_commit/indexing",
      "index-storage/mmapped-2shards-intellij_commit/indexing",
      "index-storage/mmapped-2shards-w-coroutines-intellij_commit/indexing",
      "index-storage/mmapped-intellij_commit/indexing",
      "index-storage/mmapped-w-coroutines-intellij_commit/indexing",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
