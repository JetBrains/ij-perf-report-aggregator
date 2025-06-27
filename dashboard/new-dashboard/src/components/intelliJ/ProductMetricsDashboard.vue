<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_product_dashboard_dev"
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

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: [["indexingTimeWithoutPauses", "fus_dumb_indexing_time"]],
    projects: ["community/indexing", "intellij_commit/indexing", "kotlin/indexing"],
  },
  {
    labels: ["Scanning"],
    measures: [["scanningTimeWithoutPauses", "fus_scanning_time"]],
    projects: ["community/indexing", "intellij_commit/indexing", "kotlin/indexing"],
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: [["firstCodeAnalysis", "fus_daemon_finished_full_duration_since_started_ms"]],
    projects: [
      "intellij_commit/localInspection/java_file",
      "kotlin/localInspection",
      "kotlin_coroutines/localInspection",
      "intellij_commit/localInspection/java_file/embeddedClient",
      "kotlin/localInspection/embeddedClient",
      "kotlin_coroutines/localInspection/embeddedClient",
    ],
  },
  {
    labels: ["Completion JAVA Duration"],
    measures: [["completion", "fus_completion_duration_90p", "fus_completion_duration_sum"]],
    projects: [
      "intellij_commit/completion/java_file",
      "intellij_commit/completion/java_file/embeddedClient",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "keycloak_release_20/ultimateCase/JpaUserProvider/embeddedClient",
      "train-ticket/ultimateCase/ExecuteServiceImpl",
      "train-ticket/ultimateCase/ExecuteServiceImpl/embeddedClient",
      "grails/completion/java_file",
      "grails/completion/java_file/embeddedClient",
    ],
  },
  {
    labels: ["Completion JAVA Time to Show"],
    measures: [["completion#firstElementShown#mean_value", "fus_time_to_show_90p"]],
    projects: [
      "intellij_commit/completion/java_file",
      "intellij_commit/completion/java_file/embeddedClient",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "keycloak_release_20/ultimateCase/JpaUserProvider/embeddedClient",
      "train-ticket/ultimateCase/ExecuteServiceImpl",
      "train-ticket/ultimateCase/ExecuteServiceImpl/embeddedClient",
      "grails/completion/java_file",
      "grails/completion/java_file/embeddedClient",
    ],
  },
  {
    labels: ["Completion Kotlin Duration"],
    measures: [["completion", "fus_completion_duration_90p", "fus_completion_duration_sum"]],
    projects: ["toolbox_enterprise/ultimateCase/UserController", "toolbox_enterprise/ultimateCase/UserController/embeddedClient"],
  },
  {
    labels: ["Completion Kotlin Time to Show"],
    measures: [["completion#firstElementShown#mean_value", "fus_time_to_show_90p"]],
    projects: ["toolbox_enterprise/ultimateCase/UserController", "toolbox_enterprise/ultimateCase/UserController/embeddedClient"],
  },
  {
    labels: ["Completion Others Duration"],
    measures: [["completion", "fus_completion_duration_90p", "fus_completion_duration_sum"]],
    projects: [
      "keycloak_release_20/completion/CorePomXml",
      "grails/completion/groovy_file",
      "keycloak_release_20/completion/CorePomXml/embeddedClient",
      "grails/completion/groovy_file/embeddedClient",
    ],
  },
  {
    labels: ["Completion Others Time to Show"],
    measures: [["completion#firstElementShown#mean_value", "fus_time_to_show_90p"]],
    projects: [
      "keycloak_release_20/completion/CorePomXml",
      "grails/completion/groovy_file",
      "keycloak_release_20/completion/CorePomXml/embeddedClient",
      "grails/completion/groovy_file/embeddedClient",
    ],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: [
      "community/go-to-all/Editor/typingLetterByLetter",
      "community/go-to-all-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-all/Editor/typingLetterByLetter/embeddedClient",
      "community/go-to-all-with-warmup/Editor/typingLetterByLetter/embeddedClient",
    ],
  },
  {
    labels: ["TypingCodeAnalysis"],
    measures: ["typingCodeAnalyzing"],
    projects: [
      "toolbox_enterprise/ultimateCase/SecurityTests",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "train-ticket/ultimateCase/InsidePaymentServiceImpl",
      "toolbox_enterprise/ultimateCase/SecurityTests/embeddedClient",
      "keycloak_release_20/ultimateCase/JpaUserProvider/embeddedClient",
      "train-ticket/ultimateCase/InsidePaymentServiceImpl/embeddedClient",
    ],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["kotlin_coroutines/inspection", "spring_boot_maven/inspection", "kotlin_coroutines/inspection/embeddedClient", "spring_boot_maven/inspection/embeddedClient"],
  },
  {
    labels: ["Gradle Import"],
    measures: ["ExternalSystemSyncProjectTask"],
    projects: [
      "project-import-android-500-modules/fastInstaller",
      "project-import-gradle-1000-modules/fastInstaller",
      "project-import-gradle-hibernate-orm/fastInstaller",
      "project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/fastInstaller",
    ],
  },
  {
    labels: ["Maven Import"],
    measures: ["maven.import.stats.sync.project.task"],
    projects: [
      "project-import-maven-1000-modules/fastInstaller",
      "project-import-maven-javaee8/fastInstaller",
      "project-reimport-maven-quarkus/fastInstaller",
      "project-import-maven-flink/fastInstaller",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
