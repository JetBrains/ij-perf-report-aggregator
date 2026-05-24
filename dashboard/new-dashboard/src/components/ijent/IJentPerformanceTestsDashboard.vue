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
  "intellij-community/indexing/Local",
  "intellij-community/indexing/Docker",
  "intellij-community/indexing/WSL",

  "intellij-community/build/Local",
  "intellij-community/build/Docker",
  "intellij-community/build/WSL",

  "intellij-ultimate/indexing/Local",
  "intellij-ultimate/indexing/Docker",
  "intellij-ultimate/indexing/WSL",

  "php-project/indexing/Local",
  "php-project/indexing/Docker",
  "php-project/indexing/WSL",

  "spring-pet-clinic-maven/indexing/Local",
  "spring-pet-clinic-maven/indexing/Docker",
  "spring-pet-clinic-maven/indexing/WSL",

  "spring-pet-clinic-gradle/indexing/Local",
  "spring-pet-clinic-gradle/indexing/Docker",
  "spring-pet-clinic-gradle/indexing/WSL",

  "jps-1000-modules/import/Local",
  "jps-1000-modules/import/Docker",
  "jps-1000-modules/import/WSL",
]

// Maps each canonical project name to the legacy raw names it should also pull data from.
// The canonical name itself is always queried — list only the legacy aliases here.
const projectAliases: Record<string, string[]> = {
  // "community/indexingLocal" is also one semantically of the aliases, but its data lacks core "indexingTimeWithoutPauses" metrics on Linux and on Windows
  "intellij-community/indexing/Local":  ["ijent-import-intellij-Local"],
  // "community/indexingDocker" is also one semantically of the aliases, but its data lacks core "indexingTimeWithoutPauses" metrics on Linux, and is totally absent on Windows
  "intellij-community/indexing/Docker": ["ijent-import-intellij-Docker"],
  // "community/indexingWSL" is also one semantically of the aliases, but its data lacks core "indexingTimeWithoutPauses" metrics on Windows
  "intellij-community/indexing/WSL":    [],

  // "community/rebuild_*" is based on another community revision and should not be compared to the current metrics
  // builds within "ijent-build-intellij" present on Windows, but never succeeded there and the metrics are absent on buildserver under the same directory, the build
  "intellij-community/build/Local":  [],
  "intellij-community/build/Docker": ["ijent-build-intellij-Docker"],
  "intellij-community/build/WSL":    ["ijent-build-intellij-WSL"],

  "intellij-ultimate/indexing/Local":  ["ijent-import-intellij", "nio_default-import-intellij-Local"],
  "intellij-ultimate/indexing/Docker": [],
  "intellij-ultimate/indexing/WSL":    ["wsl-import-intellij-WSL"],

  "php-project/indexing/Local":  ["php-project/indexing/Local/indexingLocal",   "indexing-php-project/indexingLocal"],
  "php-project/indexing/Docker": ["php-project/indexing/Docker/indexingDocker", "indexing-php-project/indexingDocker"],
  "php-project/indexing/WSL":    ["php-project/indexing/WSL/indexingWSL",       "indexing-php-project/indexingWSL"],

  "spring-pet-clinic-maven/indexing/Local":  ["spring-pet-clinic-maven/indexing/Local/indexingLocal",   "spring-pet-clinic-maven/indexingLocal"],
  "spring-pet-clinic-maven/indexing/Docker": ["spring-pet-clinic-maven/indexing/Docker/indexingDocker", "spring-pet-clinic-maven/indexingDocker"],
  "spring-pet-clinic-maven/indexing/WSL":    ["spring-pet-clinic-maven/indexing/WSL/indexingWSL",       "spring-pet-clinic-maven/indexingWSL"],

  "spring-pet-clinic-gradle/indexing/Local":  ["spring-pet-clinic-gradle/indexing/Local/indexingLocal",   "spring-pet-clinic-gradle/indexingLocal"],
  "spring-pet-clinic-gradle/indexing/Docker": ["spring-pet-clinic-gradle/indexing/Docker/indexingDocker", "spring-pet-clinic-gradle/indexingDocker"],
  "spring-pet-clinic-gradle/indexing/WSL":    ["spring-pet-clinic-gradle/indexing/WSL/indexingWSL",       "spring-pet-clinic-gradle/indexingWSL"],

  "jps-1000-modules/import/Local":  ["ijent-import-jps-1000-modules-Local", "nio_default-import-jps-1000-modules-Local"],
  "jps-1000-modules/import/Docker": ["ijent-import-jps-1000-modules-Docker"],
  "jps-1000-modules/import/WSL":    ["wsl-import-jps-1000-modules-WSL"],
}

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(projects)

const charts = computed(() => {
  const selected = testConfigurator.selected.value ?? []
  const rawProjects: string[] = []
  const displayAliases: string[] = []
  for (const canonical of selected) {
    rawProjects.push(canonical)
    displayAliases.push(canonical)
    for (const legacy of projectAliases[canonical] ?? []) {
      rawProjects.push(legacy)
      displayAliases.push(canonical)
    }
  }

  const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
    return {
      labels: [metric],
      measures: [metric],
      projects: rawProjects,
      aliases: displayAliases,
    }
  })
  return combineCharts(chartsDeclaration)
})
</script>
