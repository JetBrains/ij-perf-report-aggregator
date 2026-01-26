<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_kotlin_dashboard_devserver"
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
    labels: ["First Code Analysis", "Local Inspections", "File Openings: code loaded", "File Openings: tab shown"],
    measures: ["firstCodeAnalysis", "localInspections", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"],
    projects: ["kotlin/localInspection", "kotlin_coroutines/localInspection"],
  },
  {
    labels: ["Highlighting - remove symbol", "Highlighting - remove symbol warmup", "Highlighting - type symbol", "Highlighting - type symbol warmup"],
    measures: ["typing_EditorBackSpace_duration", "typing_EditorBackSpace_warmup_duration", "typing_}_duration", "typing_}_warmup_duration"],
    projects: ["intellij_commit/editor-kotlin-highlighting"],
  },
  {
    labels: ["Highlighting - remove method"],
    measures: ["replaceTextCodeAnalysis"],
    projects: ["intellij_commit/red-code-kotlin"],
  },
  {
    labels: ["Inspection", "Inspection (Full GC Pause)", "Inspection (JVM GC collection times)"],
    measures: ["globalInspections", "fullGCPause", "JVM.GC.collectionTimesMs"],
    projects: ["kotlin/inspection", "kotlin_coroutines/inspection"],
  },
  {
    labels: ["FindUsages (all usages)", "FindUsages (first usage)", "FindUsages (Full GC Pause)", "FindUsages (JVM GC collection times)"],
    measures: [["findUsages", "fus_find_usages_all"], ["findUsages_firstUsage", "fus_find_usages_first"], ["fullGCPause"], ["JVM.GC.collectionTimesMs"]],
    projects: [
      "community/findUsages/ActionsKt_runReadAction_Before",
      "intellij_commit/findUsages/ActionsKt_runReadAction",
      "community/findUsages/Persistent_absolutePath_Before",
      "intellij_commit/findUsages/Persistent_absolutePath",
    ],
  },
  {
    labels: ["Creating a new file"],
    measures: ["createKotlinFile"],
    projects: ["intellij_commit/createKotlinClass"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
