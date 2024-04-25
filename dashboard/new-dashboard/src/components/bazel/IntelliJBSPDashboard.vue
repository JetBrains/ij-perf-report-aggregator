<template>
  <DashboardPage
    db-name="bazel"
    table="report"
    persistent-id="intellij_bsp_dashboard"
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
  "add.bsp.fetched.jdks.mb",
  "add.bsp.fetched.jdks.ms",
  "apply.changes.on.workspace.model.mb",
  "apply.changes.on.workspace.model.ms",
  "calculate.all.unique.jdk.infos.mb",
  "calculate.all.unique.jdk.infos.ms",
  "collect.project.details.mb",
  "collect.project.details.ms",
  "compute.non.overlapping.targets.mb",
  "compute.non.overlapping.targets.ms",
  "create.libraries.mb",
  "create.libraries.ms",
  "create.loaded.targets.storage.mb",
  "create.loaded.targets.storage.ms",
  "create.overlapping.targets.graph.mb",
  "create.overlapping.targets.graph.ms",
  "create.target.details.for.document.provider.mb",
  "create.target.details.for.document.provider.ms",
  "create.target.id.to.module.entities.map.mb",
  "create.target.id.to.module.entities.map.ms",
  "initialize.magic.meta.model.mb",
  "initialize.magic.meta.model.ms",
  "load.default.targets.mb",
  "load.default.targets.ms",
  "load.modules.mb",
  "load.modules.ms",
  "replacebysource.in.apply.on.workspace.model.mb",
  "replacebysource.in.apply.on.workspace.model.ms",
  "replaceprojectmodel.in.apply.on.workspace.model.mb",
  "replaceprojectmodel.in.apply.on.workspace.model.ms",
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      "Bazel (ij)",
      "Bazel-BSP (ij)",
      "Datalore (ij)",
      "Synthetic 1 project (ij)",
      "Synthetic 1000 project (ij)",
      "Synthetic 10000 project (ij)",
      "Synthetic 20000 project (ij)",
      "Synthetic 5000 project (ij)",
      "gRPC-Java (ij)",
    ],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
