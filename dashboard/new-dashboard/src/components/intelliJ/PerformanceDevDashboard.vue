<template>
  <DashboardPage
    v-slot="{serverConfigurator, dashboardConfigurators, warnings}"
    db-name="perfintDev"
    table="idea"
    persistent-id="ideaDev_dashboard"
    initial-machine="Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
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
        :server-configurator="serverConfigurator"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: Array<ChartDefinition> = [{
  labels: ["Indexing", "Scanning", "Updating time"],
  measures: ["indexing", "scanning", "updatingTime"],
  projects: ["intellij_sources/indexing", "intellij_commit/indexing"],
}, {
  labels: ["Find Usages Java"],
  measures: ["findUsages"],
  projects: ["intellij_sources/findUsages/Application_runReadAction", "intellij_sources/findUsages/LocalInspectionTool_getID",
    "intellij_sources/findUsages/PsiManager_getInstance", "intellij_sources/findUsages/PropertyMapping_value",
    "intellij_commit/findUsages/Application_runReadAction", "intellij_commit/findUsages/LocalInspectionTool_getID",
    "intellij_commit/findUsages/PsiManager_getInstance", "intellij_commit/findUsages/PropertyMapping_value"],
}, {
  labels: ["Find Usages Kotlin"],
  measures: ["findUsages"],
  projects: ["intellij_sources/findUsages/ActionsKt_runReadAction", "intellij_sources/findUsages/DynamicPluginListener_TOPIC", "intellij_sources/findUsages/Path_div",
    "intellij_sources/findUsages/Persistent_absolutePath", "intellij_sources/findUsages/RelativeTextEdit_rangeTo",
    "intellij_sources/findUsages/TemporaryFolder_invoke", "intellij_sources/findUsages/Project_guessProjectDir",
    "intellij_commit/findUsages/ActionsKt_runReadAction", "intellij_commit/findUsages/DynamicPluginListener_TOPIC", "intellij_commit/findUsages/Path_div",
    "intellij_commit/findUsages/Persistent_absolutePath", "intellij_commit/findUsages/RelativeTextEdit_rangeTo",
    "intellij_commit/findUsages/TemporaryFolder_invoke", "intellij_commit/findUsages/Project_guessProjectDir"],
}, {
  labels: ["Local Inspection"],
  measures: ["localInspections"],
  projects: ["intellij_sources/localInspection/java_file", "intellij_sources/localInspection/kotlin_file",
    "intellij_commit/localInspection/java_file", "intellij_commit/localInspection/kotlin_file"],
}, {
  labels: ["Completion: execution time"],
  measures: ["completion"],
  projects: ["intellij_sources/completion/java_file", "intellij_sources/completion/kotlin_file",
    "intellij_commit/completion/java_file", "intellij_commit/completion/kotlin_file"],
}, {
  labels: ["Completion: awt delay"],
  measures: ["test#average_awt_delay"],
  projects: ["intellij_sources/completion/java_file", "intellij_sources/completion/kotlin_file",
    "intellij_commit/completion/java_file", "intellij_commit/completion/kotlin_file"],
}]

const charts = combineCharts(chartsDeclaration)
</script>