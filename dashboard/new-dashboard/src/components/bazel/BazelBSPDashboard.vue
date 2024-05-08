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
  "build.dependency.tree.memory.mb",
  "build.dependency.tree.time.ms",
  "build.reverse.sources.memory.mb",
  "build.reverse.sources.time.ms",
  "building.project.with.aspect.memory.mb",
  "building.project.with.aspect.time.ms",
  "create.ap.libraries.memory.mb",
  "create.ap.libraries.time.ms",
  "create.kotlin.stdlibs.memory.mb",
  "create.kotlin.stdlibs.time.ms",
  "create.modules.memory.mb",
  "create.modules.time.ms",
  "fetching.all.possible.target.names.memory.mb",
  "fetching.all.possible.target.names.time.ms",
  "libraries.from.jdeps.memory.mb",
  "libraries.from.jdeps.time.ms",
  "libraries.from.targets.and.deps.memory.mb",
  "libraries.from.targets.and.deps.time.ms",
  "mapping.to.internal.model.memory.mb",
  "mapping.to.internal.model.time.ms",
  "merge.all.libraries.memory.mb",
  "merge.libraries.from.deps.memory.mb",
  "parsing.aspect.outputs.memory.mb",
  "parsing.aspect.outputs.time.ms",
  "reading.aspect.output.paths.memory.mb",
  "reading.aspect.output.paths.time.ms",
  "reading.project.view.and.creating.workspace.context.memory.mb",
  "save.invalid.target.labels.memory.mb",
  "save.invalid.target.labels.time.ms",
  "select.targets.memory.mb",
  "select.targets.time.ms",
  "targets.as.libraries.memory.mb",
  "targets.as.libraries.time.ms",
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      "Bazel",
      "Bazel-BSP",
      "Datalore",
      "Synthetic 1 project",
      "Synthetic 1000 project",
      "Synthetic 10000 project",
      "Synthetic 20000 project",
      "Synthetic 5000 project",
      "gRPC-Java",
    ],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
