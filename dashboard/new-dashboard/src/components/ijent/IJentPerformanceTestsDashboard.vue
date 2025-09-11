<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_performance_dashboard"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :with-installer="false"
    :charts="charts"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="testConfigurator"
        title="Test"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
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
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsDeclaration = [
  "ijent.events.count",
  "indexingTimeWithoutPauses",
  "scanningTimeWithoutPauses",
  "numberOfIndexedFiles",
  "build_compilation_duration",

  "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
  "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",

  "AWTEventQueue.dispatchTimeTotal",
  "gcPause",
  "gcPauseCount",
  "fullGCPause",
  "freedMemoryByGC",
  "totalHeapUsedMax",

  "JVM.maxThreadCount",
  "JVM.totalCpuTimeMs",
  "JVM.GC.collectionTimesMs",
  "JVM.GC.collections",
  "JVM.maxHeapMegabytes",
  "MEM.avgRamMegabytes",
  "MEM.avgFileMappingsRamMegabytes",

  "ijent.directoryStreamClose.sum.ms",
  "ijent.directoryStreamIteratorNext.sum.ms",
  "ijent.fileSystemClose.sum.ms",
  "ijent.providerCheckAccess.sum.ms",
  "ijent.providerCopy.sum.ms",
  "ijent.providerCreateDirectory.sum.ms",
  "ijent.providerDelete.sum.ms",
  "ijent.providerMove.sum.ms",
  "ijent.providerNewByteChannel.sum.ms",
  "ijent.providerNewDirectoryStream.sum.ms",
  "ijent.providerReadAttributes.sum.ms",
  "ijent.seekableByteChannelClose.sum.ms",
  "ijent.seekableByteChannelNewPosition.sum.ms",
  "ijent.seekableByteChannelRead.sum.ms",
  "ijent.seekableByteChannelSize.sum.ms",
  "ijent.seekableByteChannelWrite.sum.ms",

  "project.opening",
  "jps.aggregate.sync.duration",
  "jps.aggregate.counters",

  "jps.app.storage.content.reader.load.component.ms",
  "jps.project.serializers.load.ms",
  "jps.storage.jps.conf.reader.load.component.ms",

  "workspaceModel.loading.total.ms",
  "workspaceModel.moduleBridgeLoader.loading.modules.ms",
  "workspaceModel.moduleManagerBridge.load.module.ms",
]

const projects = [
  "wsl-import-jps-1000-modules-WSL",
  "wsl-import-intellij",
  "nio_default-import-jps-1000-modules-Local",
  "nio_default-import-intellij",
  "ijent-import-jps-1000-modules-Local",
  "ijent-import-jps-1000-modules-Docker",
  "ijent-import-intellij-Local",
  "ijent-import-intellij-Docker",
  "ijent-import-intellij",
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
