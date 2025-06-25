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
    projects: ["intellij_commit/localInspection/java_file", "kotlin/localInspection", "kotlin_coroutines/localInspection"],
  },
  {
    labels: ["Completion JAVA Duration"],
    measures: [["completion", "fus_completion_duration_90p", "fus_completion_duration_sum"]],
    projects: [
      "intellij_commit/completion/java_file",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "train-ticket/ultimateCase/ExecuteServiceImpl",
      "grails/completion/java_file",
    ],
  },
  {
    labels: ["Completion JAVA Time to Show"],
    measures: [["completion#firstElementShown#mean_value", "fus_time_to_show_90p"]],
    projects: [
      "intellij_commit/completion/java_file",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "train-ticket/ultimateCase/ExecuteServiceImpl",
      "grails/completion/java_file",
    ],
  },
  {
    labels: ["Completion Kotlin Duration"],
    measures: [["completion", "fus_completion_duration_90p", "fus_completion_duration_sum"]],
    projects: ["toolbox_enterprise/ultimateCase/UserController"],
  },
  {
    labels: ["Completion Kotlin Time to Show"],
    measures: [["completion#firstElementShown#mean_value", "fus_time_to_show_90p"]],
    projects: ["toolbox_enterprise/ultimateCase/UserController"],
  },
  {
    labels: ["Completion Others Duration"],
    measures: [["completion", "fus_completion_duration_90p", "fus_completion_duration_sum"]],
    projects: ["keycloak_release_20/completion/CorePomXml", "grails/completion/groovy_file"],
  },
  {
    labels: ["Completion Others Time to Show"],
    measures: [["completion#firstElementShown#mean_value", "fus_time_to_show_90p"]],
    projects: ["keycloak_release_20/completion/CorePomXml", "grails/completion/groovy_file"],
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
    projects: ["toolbox_enterprise/ultimateCase/SecurityTests", "keycloak_release_20/ultimateCase/JpaUserProvider", "train-ticket/ultimateCase/InsidePaymentServiceImpl"],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["kotlin_coroutines/inspection", "spring_boot_maven/inspection"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
