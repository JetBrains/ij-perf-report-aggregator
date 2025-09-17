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
    projects: [
      "toolbox_enterprise/ultimateCase/SecurityTests",
      "toolbox_enterprise/ultimateCase/ToolController",
      "toolbox_enterprise/ultimateCase/ToolService",
      "toolbox_enterprise/ultimateCase/UserController",
      "toolbox_enterprise/ultimateCase/UserRepository",
      "kotlin/localInspection",
      "kotlin_coroutines/localInspection",
    ],
  },
  {
    labels: ["Typing (typingCodeAnalyzing)", "Typing (average_awt_delay)", "Typing (max_awt_delay)"],
    measures: ["typingCodeAnalyzing", "test#average_awt_delay", "test#max_awt_delay"],
    projects: [
      "toolbox_enterprise/ultimateCase/SecurityTests",
      "toolbox_enterprise/ultimateCase/ToolController",
      "toolbox_enterprise/ultimateCase/ToolService",
      "toolbox_enterprise/ultimateCase/UserController",
      "toolbox_enterprise/ultimateCase/UserRepository",
    ],
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
    labels: ["Completion", "Completion 90p", "Completion time to show 90p"],
    measures: ["completion", "fus_completion_duration_90p", "fus_time_to_show_90p"],
    projects: [
      "toolbox_enterprise/ultimateCase/SecurityTests",
      "toolbox_enterprise/ultimateCase/ToolController",
      "toolbox_enterprise/ultimateCase/ToolService",
      "toolbox_enterprise/ultimateCase/UserController",
      "toolbox_enterprise/ultimateCase/UserRepository",
    ],
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
    labels: ["Show Intentions (average awt delay)", "Show Intentions (awt dispatch time)"],
    measures: ["test#average_awt_delay", "AWTEventQueue.dispatchTimeTotal"],
    projects: ["kotlin/showIntention/Import"],
  },
  {
    labels: ["Creating a new file"],
    measures: ["createKotlinFile"],
    projects: ["intellij_commit/createKotlinClass"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
