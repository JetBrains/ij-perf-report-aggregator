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
    labels: ["IsUpToDateCheck duration"],
    measures: ["isUpToDateCheck"],
    projects: ["community/findUsages/PsiManager_getInstance_Before", "community/findUsages/PsiManager_getInstance_After"],
  },
  {
    labels: ["FindUsages PsiManager#getInstance Before and After Compilation (all usages)", "FindUsages PsiManager#getInstance Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/PsiManager_getInstance_Before", "community/findUsages/PsiManager_getInstance_After"],
  },
  {
    labels: [
      "FindUsages Library#getName Before and After Compilation (all usages)",
      "FindUsages Library#getName Before and After Compilation (first usage)",
      "JVM Total Time to safepoints",
      "Full GC Pause",
      "JVM GC collection times",
    ],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
      ["JVM.totalTimeToSafepointsMs"],
      ["fullGCPause"],
      ["JVM.GC.collectionTimesMs"],
    ],
    projects: ["community/findUsages/Library_getName_Before", "community/findUsages/Library_getName_After", "intellij_commit/findUsages/Library_getName"],
  },
  {
    labels: ["FindUsages LocalInspectionTool#getID Before and After Compilation (all usages)", "FindUsages LocalInspectionTool#getID Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/LocalInspectionTool_getID_Before", "community/findUsages/LocalInspectionTool_getID_After"],
  },
  {
    labels: ["FindUsages Application#runReadAction Before and After Compilation (all usages)", "FindUsages Application#runReadAction Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: [
      "community/findUsages/Application_runReadAction_Before",
      "community/findUsages/Application_runReadAction_After",
      "intellij_commit/findUsages/Application_runReadAction",
    ],
  },
  {
    labels: ["FindUsages ActionsKt#runReadAction Before and After Compilation (all usages)", "FindUsages ActionsKt#runReadAction Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/ActionsKt_runReadAction_Before", "community/findUsages/ActionsKt_runReadAction_After"],
  },
  {
    labels: ["FindUsages Persistent#absolutePath Before and After Compilation (all usages)", "FindUsages Persistent#absolutePath Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Persistent_absolutePath_Before", "community/findUsages/Persistent_absolutePath_After"],
  },
  {
    labels: ["FindUsages PropertyMapping#value Before and After Compilation (all usages)", "FindUsages PropertyMapping#value Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/PropertyMapping_value_Before", "community/findUsages/PropertyMapping_value_After"],
  },
  {
    labels: ["FindUsages Object#hashCode Before and After Compilation (all usages)", "FindUsages Object#hashCode Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Object_hashCode_Before", "community/findUsages/Object_hashCode_After"],
  },
  {
    labels: ["FindUsages Objects#hashCode Before and After Compilation (all usages)", "FindUsages Objects#hashCode Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Objects_hashCode_Before", "community/findUsages/Objects_hashCode_After"],
  },
  {
    labels: ["FindUsages Path#div Before and After Compilation (all usages)", "FindUsages Path#div Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Path_div_Before", "community/findUsages/Path_div_After"],
  },
  {
    labels: ["FindUsages String#toString Before and After Compilation (all usages)", "FindUsages String#toString Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/String_toString_Before", "community/findUsages/String_toString_After", "intellij_commit/findUsages/String_toString"],
  },
  {
    labels: ["Find Usages with idea.is.internal=true Before Compilation", "First found usage"],
    measures: ["findUsages", "findUsages_firstUsage"],
    projects: [
      "intellij_commit/findUsages/PsiManager_getInstance_firstCall",
      "intellij_commit/findUsages/PsiManager_getInstance_secondCall",
      "intellij_commit/findUsages/PsiManager_getInstance_thirdCallInternalMode",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
