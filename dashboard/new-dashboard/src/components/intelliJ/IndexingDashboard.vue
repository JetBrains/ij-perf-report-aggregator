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
    labels: ["Indexing (Big projects)"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["community/indexing", "intellij_commit/indexing", "community/indexingWithHighlighting", "indexing-community/indexingWithHighlighting"],
    aliases: ["Community", "IntelliJ Commit", "Community with Highlighting", "Indexing Community with Highlighting"],
  },
  {
    labels: ["Indexing (after checkout)", "Scanning (after checkout)"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses"],
    projects: ["intellij_sources/checkout/243"],
    aliases: ["IntelliJ Sources Checkout"],
  },
  {
    labels: ["Scanning (Big projects)"],
    measures: ["scanningTimeWithoutPauses"],
    projects: [
      "community/indexing",
      "intellij_commit/indexing",
      "community/second-scanning",
      "intellij_commit/second-scanning",
      "community/third-scanning",
      "intellij_commit/third-scanning",
      "intellij_commit/allow-skipping-full-scanning",
    ],
    aliases: [
      "Community First Scan",
      "IntelliJ First Scan",
      "Community Second Scan",
      "IntelliJ Second Scan",
      "Community Third Scan",
      "IntelliJ Third Scan",
      "IntelliJ Skip Full Scan",
    ],
  },
  {
    labels: ["Number of indexed files (Big projects)"],
    measures: ["numberOfIndexedFiles"],
    projects: [
      "community/indexing",
      "intellij_commit/indexing",
      "community/second-scanning",
      "intellij_commit/second-scanning",
      "community/third-scanning",
      "intellij_commit/third-scanning",
      "intellij_commit/allow-skipping-full-scanning",
    ],
    aliases: [
      "Community First Scan",
      "IntelliJ First Scan",
      "Community Second Scan",
      "IntelliJ Second Scan",
      "Community Third Scan",
      "IntelliJ Third Scan",
      "IntelliJ Skip Full Scan",
    ],
  },
  {
    labels: ["Number of indexed files with writing index value (Big projects)", "Number of indexed files with nothing to write (Big projects)"],
    measures: ["numberOfIndexedFilesWritingIndexValue", "numberOfIndexedFilesWithNothingToWrite"],
    projects: ["community/indexing", "intellij_commit/indexing"],
    aliases: ["Community", "IntelliJ Commit"],
  },
  {
    labels: ["Processing Time Java (Big projects)", "Processing Time Kotlin (Big projects)"],
    measures: ["processingTime#JAVA", "processingTime#Kotlin"],
    projects: ["community/indexing", "intellij_commit/indexing"],
    aliases: ["Community", "IntelliJ Commit"],
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
      "keycloak_release_20/indexing",
      "toolbox_enterprise/indexing",
      "train-ticket/indexing",
      "space/indexing",
    ],
    aliases: [
      "Empty Project",
      "Grails",
      "Java",
      "Kotlin",
      "Kotlin Coroutines",
      "Spring Boot",
      "Spring Boot Maven",
      "Kotlin PetClinic",
      "Keycloak 20",
      "Toolbox Enterprise",
      "Train Ticket",
      "Space",
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
      "keycloak_release_20/indexing",
      "toolbox_enterprise/indexing",
      "train-ticket/indexing",
      "space/indexing",
    ],
    aliases: [
      "Empty Project",
      "Grails",
      "Java",
      "Kotlin",
      "Kotlin Coroutines",
      "Spring Boot",
      "Spring Boot Maven",
      "Kotlin PetClinic",
      "Keycloak 20",
      "Toolbox Enterprise",
      "Train Ticket",
      "Space",
    ],
  },
  {
    labels: ["Second Scanning", "Second Number of indexed files"],
    measures: ["scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: [
      "empty_project/second-scanning",
      "grails/second-scanning",
      "java/second-scanning",
      "kotlin/second-scanning",
      "kotlin_coroutines/second-scanning",
      "spring_boot/second-scanning",
      "spring_boot_maven/second-scanning",
      "kotlin_petclinic/second-scanning",
      "keycloak_release_20/second-scanning",
      "toolbox_enterprise/second-scanning",
      "train-ticket/second-scanning",
      "space/second-scanning",
    ],
    aliases: [
      "Empty Project (2nd)",
      "Grails (2nd)",
      "Java (2nd)",
      "Kotlin (2nd)",
      "Kotlin Coroutines (2nd)",
      "Spring Boot (2nd)",
      "Spring Boot Maven (2nd)",
      "Kotlin PetClinic (2nd)",
      "Keycloak 20 (2nd)",
      "Toolbox Enterprise (2nd)",
      "Train Ticket (2nd)",
      "Space (2nd)",
    ],
  },
  {
    labels: ["Third Scanning", "Third Number of indexed files"],
    measures: ["scanningTimeWithoutPauses", "numberOfIndexedFiles"],
    projects: [
      "empty_project/third-scanning",
      "grails/third-scanning",
      "java/third-scanning",
      "kotlin/third-scanning",
      "kotlin_coroutines/third-scanning",
      "spring_boot/third-scanning",
      "spring_boot_maven/third-scanning",
      "kotlin_petclinic/third-scanning",
      "keycloak_release_20/third-scanning",
      "toolbox_enterprise/third-scanning",
      "train-ticket/third-scanning",
      "space/third-scanning",
    ],
    aliases: [
      "Empty Project (3rd)",
      "Grails (3rd)",
      "Java (3rd)",
      "Kotlin (3rd)",
      "Kotlin Coroutines (3rd)",
      "Spring Boot (3rd)",
      "Spring Boot Maven (3rd)",
      "Kotlin PetClinic (3rd)",
      "Keycloak 20 (3rd)",
      "Toolbox Enterprise (3rd)",
      "Train Ticket (3rd)",
      "Space (3rd)",
    ],
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
    aliases: [
      "Default (2 Shards)",
      "Default",
      "Default + Coroutines",
      "Default + MRU Cache",
      "Fake Writer",
      "Memory Mapped (2 Shards)",
      "Memory Mapped (2 Shards + Coroutines)",
      "Memory Mapped",
      "Memory Mapped + Coroutines",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
