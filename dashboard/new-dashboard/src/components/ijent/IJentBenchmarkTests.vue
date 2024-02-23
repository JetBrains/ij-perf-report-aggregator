<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="ijent_benchmarks_dashboard"
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
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsDeclaration = [
  "ijent.file.exists.events.count",
  "ijent.file.exists.median.ns",
  "ijent.file.exists.standard.deviation.ns",
  "ijent.file.exists.mad.ns",
  "ijent.file.exists.range.ns",
  "ijent.file.exists.95.percentile.ns",
  "ijent.file.exists.99.percentile.ns",
  "ijent.file.exists.min.ns",
  "ijent.file.exists.max.ns",

  "AWTEventQueue.dispatchTimeTotal",
  "gcPause",
  "gcPauseCount",
  "fullGCPause",
  "freedMemoryByGC",
  "totalHeapUsedMax",

  "JVM.GC.collectionTimesMs",
  "JVM.GC.collections",
  "JVM.maxHeapMegabytes",
  "JVM.maxThreadCount",
  "JVM.totalCpuTimeMs",
]

const projects = [
  "com.intellij.platform.ijent.performance.benchmarks.IjentWslNioFsBenchmarkTest.WSL - Files Exists - Provider checkAccess for existing files - WSL",
  "com.intellij.platform.ijent.performance.benchmarks.IjentWslNioFsBenchmarkTest.IJent - Files Exists - Provider checkAccess for existing files - IJENT",
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
