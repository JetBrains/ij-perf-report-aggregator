<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_find_usages_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
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
    labels: ["FindUsages PsiManager#getInstance Before and After Compilation"],
    measures: ["findUsages"],
    projects: ["community/findUsages/PsiManager_getInstance_Before", "community/findUsages/PsiManager_getInstance_After"],
  },
  {
    labels: ["FindUsages Library#getName (all usages)", "FindUsages Library#getName (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Library_getName_Before", "community/findUsages/Library_getName_After"],
  },
  {
    labels: ["FindUsages LocalInspectionTool#getID Before and After Compilation"],
    measures: ["findUsages"],
    projects: ["community/findUsages/LocalInspectionTool_getID_Before", "community/findUsages/LocalInspectionTool_getID_After"],
  },
  {
    labels: ["FindUsages ActionsKt#runReadAction and Application#runReadAction Before and After Compilation"],
    measures: ["findUsages"],
    projects: [
      "community/findUsages/ActionsKt_runReadAction_Before",
      "community/findUsages/ActionsKt_runReadAction_After",
      "community/findUsages/Application_runReadAction_Before",
      "community/findUsages/Application_runReadAction_After",
    ],
  },
  {
    labels: ["FindUsages Persistent#absolutePath and PropertyMapping#value Before and After Compilation"],
    measures: ["findUsages"],
    projects: [
      "community/findUsages/Persistent_absolutePath_Before",
      "community/findUsages/Persistent_absolutePath_After",
      "community/findUsages/PropertyMapping_value_Before",
      "community/findUsages/PropertyMapping_value_After",
    ],
  },
  {
    labels: ["FindUsages Object#hashCode and Path#toString Before and After Compilation"],
    measures: ["findUsages"],
    projects: [
      "community/findUsages/Object_hashCode_Before",
      "community/findUsages/Object_hashCode_After",
      "community/findUsages/Path_toString_Before",
      "community/findUsages/Path_toString_After",
    ],
  },
  {
    labels: ["FindUsages Objects#hashCode Before and After Compilation", "FindUsages Objects#hashCode Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Objects_hashCode_Before", "community/findUsages/Objects_hashCode_After"],
  },
  {
    labels: ["FindUsages Path#div Before and After Compilation"],
    measures: ["findUsages"],
    projects: ["community/findUsages/Path_div_Before", "community/findUsages/Path_div_After"],
  },
  {
    labels: ["Find Usages with idea.is.internal=true Before Compilation"],
    measures: ["findUsages"],
    projects: [
      "intellij_sources/findUsages/PsiManager_getInstance_firstCall",
      "intellij_sources/findUsages/PsiManager_getInstance_secondCall",
      "intellij_sources/findUsages/PsiManager_getInstance_thirdCallInternalMode",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
