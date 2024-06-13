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
  "add.bsp.fetched.jdks.ms",
  "apply.changes.on.workspace.model.ms",
  "bsp.max.used.memory.mb",
  "bsp.sync.project.ms",
  "bsp.used.at.exit.mb",
  "calculate.all.unique.jdk.infos.ms",
  "collect.project.details.ms",
  "create.target.id.to.module.entities.map.ms",
  "load.modules.ms",
  "replacebysource.in.apply.on.workspace.model.ms",
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
