<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="ideaDev_dashboard"
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
    labels: ["Indexing", "Scanning", "Dumb mode time", "Updating time"],
    measures: [["indexingWithoutPauses", "indexing"], ["scanningWithoutPauses", "scanning"], "dumbModeTimeWithPauses", "updatingTime"],
    projects: ["intellij_sources/indexing", "intellij_commit/indexing"],
  },
  {
    labels: ["Find Usages Java Application_runReadAction (all usages)", "Find Usages Java Application_runReadAction (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["intellij_commit/findUsages/Application_runReadAction"],
  },
  {
    labels: ["Find Usages Java Library_getName (all usages)", "Find Usages Java Library_getName (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["intellij_commit/findUsages/Library_getName"],
  },
  {
    labels: ["Find Usages Java PsiManager_getInstance (all usages)", "Find Usages Java PsiManager_getInstance (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["intellij_commit/findUsages/PsiManager_getInstance"],
  },
  {
    labels: ["Find Usages Java PropertyMapping_value (all usages)", "Find Usages Java PropertyMapping_value (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["intellij_commit/findUsages/PropertyMapping_value"],
  },
  {
    labels: ["Find Usages Kotlin ActionsKt_runReadAction (all usages)", "Find Usages Kotlin ActionsKt_runReadAction (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["intellij_commit/findUsages/ActionsKt_runReadAction"],
  },
  {
    labels: ["Find Usages Kotlin DynamicPluginListener_TOPIC (all usages)", "Find Usages Kotlin DynamicPluginListener_TOPIC (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["intellij_commit/findUsages/DynamicPluginListener_TOPIC"],
  },
  {
    labels: ["Find Usages Kotlin Path_div (all usages)", "Find Usages Kotlin Path_div (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["intellij_commit/findUsages/Path_div"],
  },
  {
    labels: ["Find Usages Kotlin Persistent_absolutePath (all usages)", "Find Usages Kotlin Persistent_absolutePath (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["intellij_commit/findUsages/Persistent_absolutePath"],
  },
  {
    labels: ["Local Inspection"],
    measures: ["localInspections"],
    projects: [
      "intellij_sources/localInspection/java_file",
      "intellij_sources/localInspection/kotlin_file",
      "intellij_commit/localInspection/java_file",
      "intellij_commit/localInspection/kotlin_file",
    ],
  },
  {
    labels: ["Completion: execution time"],
    measures: ["completion"],
    projects: [
      "intellij_sources/completion/java_file",
      "intellij_sources/completion/kotlin_file",
      "intellij_commit/completion/java_file",
      "intellij_commit/completion/kotlin_file",
    ],
  },
  {
    labels: ["Completion: awt delay"],
    measures: ["test#average_awt_delay"],
    projects: [
      "intellij_sources/completion/java_file",
      "intellij_sources/completion/kotlin_file",
      "intellij_commit/completion/java_file",
      "intellij_commit/completion/kotlin_file",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
