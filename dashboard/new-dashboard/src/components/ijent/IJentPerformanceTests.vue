<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_performance_dashboard"
    initial-machine="windows-azure"
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

// TODO: add future projects here
const projects = ["NO PROJECTS"]

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
