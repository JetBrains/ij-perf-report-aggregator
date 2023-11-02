<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_maven_dashboard"
    initial-machine="linux-blade-hetzner"
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
import { ChartDefinition, combineCharts } from "../../../charts/DashboardCharts"
import GroupProjectsChart from "../../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../../common/DashboardPage.vue"
import { MAVEN_PROJECTS } from "./maven-projects"

const metricsDeclaration = [
  // Main flow
  "maven.sync.duration",
  "maven.projects.processor.resolving.task",
  "maven.projects.processor.reading.task",
  "maven.import.stats.plugins.resolving.task",
  "maven.import.stats.applying.model.task",
  "maven.import.stats.sync.project.task",
  "maven.import.after.import.configuration",
  "maven.project.importer.base.refreshing.files.task",
  "maven.project.importer.post.importing.task.marker",
  "post_import_tasks_run.total_duration_ms",

  // Workspace commit
  "workspace_commit.attempts",
  "workspace_commit.duration_in_background_ms",
  "workspace_commit.duration_in_write_action_ms",
  "workspace_commit.duration_of_workspace_update_call_ms",

  // Workspace import
  "workspace_import.commit.duration_ms",
  "workspace_import.duration_ms",
  "workspace_import.legacy_importers.duration_ms",
  "workspace_import.legacy_importers.stats.duration_of_bridges_creation_ms",
  "workspace_import.legacy_importers.stats.duration_of_bridges_commit_ms",
  "workspace_import.populate.duration_ms",

  // Legacy import
  "legacy_import.create_modules.duration_ms",
  "legacy_import.delete_obsolete.duration_ms",
  "legacy_import.duration_ms",
  "legacy_import.importers.duration_ms",

  // IDE/Java common metrics
  "AWTEventQueue.dispatchTimeTotal",
  "CPU | Load |Total % 95th pctl",
  "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
  "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
  "gcPause",
  "gcPauseCount",
  "fullGCPause",
  "freedMemoryByGC",
  "totalHeapUsedMax",
  "JVM.GC.collectionTimesMs",
  "JVM.GC.collections",
  "JVM.maxHeapMegabytes",
  "JVM.threadCount",
  "JVM.totalCpuTimeMs",
  "JVM.totalMegabytesAllocated",
  "JVM.usedHeapMegabytes",
  "JVM.usedNativeMegabytes",
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
  return {
    labels: [metric],
    measures: [metric],
    projects: MAVEN_PROJECTS,
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
