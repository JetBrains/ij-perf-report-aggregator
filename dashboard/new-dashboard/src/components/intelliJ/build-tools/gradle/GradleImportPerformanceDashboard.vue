<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_gradle_dashboard"
    initial-machine="linux-blade-hetzner"
    :charts="charts"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="testConfigurator"
        title="Test"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4 text-gray-500" />
        </template>
      </MeasureSelect>
    </template>
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
import { computed } from "vue"
import { SimpleMeasureConfigurator } from "../../../../configurators/SimpleMeasureConfigurator"
import { ChartDefinition, combineCharts } from "../../../charts/DashboardCharts"
import GroupProjectsChart from "../../../charts/GroupProjectsChart.vue"
import MeasureSelect from "../../../charts/MeasureSelect.vue"
import DashboardPage from "../../../common/DashboardPage.vue"

const metricsDeclaration = [
  // - total sync time
  "ExternalSystemSyncProjectTask",
  // GradleExecution - full time of the sink operation, with all our overhead for preparation
  "gradle.sync.duration",
  // - work inside Gradle connection, operations that are performed inside connection
  "GradleConnection",
  // GradleCall - resolving models from daemon
  "GRADLE_CALL",
  // - processing the data we received from Gradle
  "ExternalSystemSyncResultProcessing",
  // ProjectDataServices - work of data services
  "DATA_SERVICES",
  // GradleProjectResolvers - project resolve
  "PROJECT_RESOLVERS",
  // WorkspaceModelApply - apply ws model
  "WORKSPACE_MODEL_APPLY",
  "fus_gradle.sync",

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

const projects = [
  "grazie-platform-project-import-gradle/measureStartup",
  "project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup",
  "project-import-gradle-micronaut/measureStartup",
  "project-import-gradle-hibernate-orm/measureStartup",
  "project-import-gradle-cas/measureStartup",
  "project-reimport-gradle-cas/measureStartup",
  "project-import-from-cache-gradle-cas/measureStartup",
  "project-import-gradle-1000-modules/measureStartup",
  "project-import-gradle-1000-modules-limited-ram/measureStartup",
  "project-import-gradle-5000-modules/measureStartup",
  "project-import-gradle-android-extra-large/measureStartup",
  "project-import-android-500-modules/measureStartup",
  "project-reimport-space/measureStartup",
  "project-import-space/measureStartup",
  "project-import-open-telemetry/measureStartup",
]

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(projects)

const charts = computed(() => {
  const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
    return {
      labels: [metric],
      measures: [metric],
      projects: testConfigurator.selected.value ?? [],
    }
  })
  return combineCharts(chartsDeclaration)
})
</script>
