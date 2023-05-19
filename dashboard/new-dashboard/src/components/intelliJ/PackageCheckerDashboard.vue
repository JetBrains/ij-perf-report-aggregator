<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_package_checker_dashboard"
    initial-machine="Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
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
  </DashboardPage>>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: Array<ChartDefinition> = [{
  labels: ["Package Checker execution time", "Total heap max", "Freed memory by GC", "GC pause count", "Full GC pause", "GC pause"],
  measures: ["runServiceInPlugin", "totalHeapUsedMax", "freedMemoryByGC", "gcPauseCount", "fullGCPause", "gcPause"],
  projects: ["package-checker-gradle-500-modules/get_declared_dependencies", "package-checker-gradle-500-modules/get_all_modules/maven",
    "package-checker-gradle-500-modules/get_all_modules/gradle"]
}]

const charts = combineCharts(chartsDeclaration)
</script>