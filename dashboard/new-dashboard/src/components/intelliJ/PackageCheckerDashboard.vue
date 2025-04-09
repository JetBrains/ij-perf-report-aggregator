<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_package_checker_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
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
  >
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Package Checker execution time", "Total heap max", "Freed memory by GC", "GC pause count", "Full GC pause", "GC pause"],
    measures: ["runServiceInPlugin", "totalHeapUsedMax", "freedMemoryByGC", "gcPauseCount", "fullGCPause", "gcPause"],
    projects: [
      "package-checker-gradle-500-modules/get_declared_dependencies",
      "package-checker-gradle-500-modules/get_imported_dependencies",
      "package-checker-gradle-500-modules/get_all_modules/maven",
      "package-checker-gradle-500-modules/get_all_modules/gradle",
      "package-checker-gradle-500-modules/get_all_modules/maven/5",
      "package-checker-gradle-500-modules/get_all_modules/gradle/5",
      "package-checker-gradle-500-modules/get_all_modules/maven/15",
      "package-checker-gradle-500-modules/get_all_modules/gradle/15",
      "package-checker-npm/edit-package-json/transitives/true",
      "package-checker-npm/edit-package-json/transitives/false",
      "gradlebuilddepchecker/show-warnings/gradle",
      "gradlebuildktsdepchecker/show-warnings/gradle-kts",
      "vulnerable-path-gradle-test/show_warnings/kotlin",
      "kotlin_petclinic/run_inspections/vulnerable_api_inspection",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
