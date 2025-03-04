<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_java_dashboard_devserver"
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
    labels: ["Indexing", "Scanning", "Number of indexed files", "Number of indexed files with writing index value", "Number of indexed files with nothing to write"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles", "numberOfIndexedFilesWritingIndexValue", "numberOfIndexedFilesWithNothingToWrite"],
    projects: ["empty_project/indexing", "grails/indexing", "java/indexing", "spring_boot/indexing", "spring_boot_maven/indexing"],
  },
  {
    labels: ["Rebuild"],
    measures: ["build_compilation_duration"],
    projects: ["grails/rebuild", "java/rebuild", "spring_boot/rebuild"],
  },
  {
    labels: ["Inspection"],
    measures: ["globalInspections"],
    projects: ["java/inspection", "grails/inspection", "spring_boot_maven/inspection", "spring_boot/inspection"],
  },
  {
    labels: ["Show Intentions (average awt delay)", "Show Intentions (awt dispatch time)"],
    measures: ["test#average_awt_delay", "AWTEventQueue.dispatchTimeTotal"],
    projects: ["grails/showIntentions/Find cause", "spring_boot/showIntentions"],
  },
  {
    labels: ["Completion", "Completion 90p"],
    measures: ["completion", "fus_completion_duration_90p"],
    projects: ["grails/completion/groovy_file", "grails/completion/java_file", "intellij_commit/completion/java_file"],
  },
  {
    labels: ["Rename method, rename class, change signature, move"],
    measures: [["performInlineRename", "performInlineRename", "changeJavaSignature: add parameter", "moveClassToPackage"]],
    projects: ["hadoop_commit/rename-method", "hadoop_commit/rename-class", "hadoop_commit/change-signature", "hadoop_commit/move-class"],
  },
  {
    labels: ["Inline method"],
    measures: ["inlineJavaMethod"],
    projects: ["hadoop_commit/inline-method"],
  },
  {
    labels: ["Rename package"],
    measures: ["renameDirectoryAsPackage"],
    projects: ["hadoop_commit/rename-package"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
