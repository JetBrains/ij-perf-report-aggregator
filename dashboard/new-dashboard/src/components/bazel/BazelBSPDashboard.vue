<template>
  <DashboardPage
    db-name="bazel"
    table="report"
    persistent-id="basel_bsp_dashboard"
    initial-machine="Linux EC2 M5ad.2xlarge (8 vCPU Xeon, 32 GB)"
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

const metricsDeclaration = [
  "resolve.project.time.ms",
  "max.used.memory.mb",
  "used.at.exit.mb",
  "building.project.with.aspect.time.ms",
  "mapping.to.internal.model.time.ms",
  "parsing.aspect.outputs.time.ms",
  "create.modules.time.ms",
  "reading.aspect.output.paths.time.ms",
  "fetching.all.possible.target.names.time.ms",
  "discovering.supported.external.rules.time.ms",
  "select.targets.time.ms",
  "libraries.from.jdeps.time.ms",
  "libraries.from.targets.and.deps.time.ms",
  "build.dependency.tree.time.ms",
  "build.reverse.sources.time.ms",
  "targets.as.libraries.time.ms",
  "create.ap.libraries.time.ms",
  "create.kotlin.stdlibs.time.ms",
  "save.invalid.target.labels.time.ms",
  "libraries.from.transitive.compile-time.jars.time.ms",
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      "Bazel",
      "Bazel-BSP",
      "Datalore",
      "gRPC-Java",
      "Hirschgarten",
      "Synthetic 1 project",
      "Synthetic 1000 project",
      "Synthetic 10000 project",
      "Synthetic 20000 project",
      "Synthetic 5000 project",
    ],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
