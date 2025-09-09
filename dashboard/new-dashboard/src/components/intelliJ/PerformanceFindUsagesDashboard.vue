<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_find_usages_dashboard_dev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsWithClientChart
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
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["FindUsages Library#getName (all usages)", "FindUsages Library#getName (first usage)", "JVM Total Time to safepoints", "Full GC Pause", "JVM GC collection times"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
      ["JVM.totalTimeToSafepointsMs"],
      ["fullGCPause"],
      ["JVM.GC.collectionTimesMs"],
    ],
    projects: ["community/findUsages/Library_getName_Before", "intellij_commit/findUsages/Library_getName"],
  },
  {
    labels: ["FindUsages Application#runReadAction (all usages)", "FindUsages Application#runReadAction (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Application_runReadAction_Before", "intellij_commit/findUsages/Application_runReadAction"],
  },
  {
    labels: ["FindUsages ActionsKt#runReadAction (all usages)", "FindUsages ActionsKt#runReadAction (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/ActionsKt_runReadAction_Before", "intellij_commit/findUsages/ActionsKt_runReadAction"],
  },
  {
    labels: ["FindUsages Persistent#absolutePath (all usages)", "FindUsages Persistent#absolutePath (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Persistent_absolutePath_Before", "intellij_commit/findUsages/Persistent_absolutePath"],
  },
  {
    labels: ["FindUsages String#toString (all usages)", "FindUsages String#toString (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/String_toString_Before", "intellij_commit/findUsages/String_toString"],
  },
  {
    labels: ["Find Usages with idea.is.internal=true", "First found usage"],
    measures: ["findUsages", "findUsages_firstUsage"],
    projects: ["intellij_commit/findUsages/Library_getName", "intellij_commit/findUsages/Application_runReadAction", "intellij_commit/findUsages/String_toString"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
