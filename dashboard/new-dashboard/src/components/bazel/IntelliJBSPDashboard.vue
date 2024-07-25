<template>
  <DashboardPage
    db-name="bazel"
    table="report"
    persistent-id="intellij_bsp_dashboard"
    initial-machine="Linux EC2 M5ad.2xlarge (8 vCPU Xeon, 32 GB)"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in chartsDeclaration"
        :key="chart.labels.join('')"
        :label="chart.labels[chart.labels.length - 1]"
        :measure="chart.measures"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsDeclaration = [
  ["bsp.sync.project.ms"],
  ["bsp.used.at.exit.mb", "bsp.used.after.sync.mb"],
  ["bsp.used.after.indexing.mb"],
  ["bsp.max.used.memory.mb"],
  ["collect.project.details.ms"],
  ["apply.changes.on.workspace.model.ms"],
  ["replacebysource.in.apply.on.workspace.model.ms"],
  ["replaceprojectmodel.in.apply.on.workspace.model.ms"],
  ["add.bsp.fetched.jdks.ms"],
  ["create.target.id.to.module.entities.map.ms"],
  ["load.modules.ms"],
  ["calculate.all.unique.jdk.infos.ms"],
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metricGroup) => {
  return {
    labels: metricGroup,
    measures: metricGroup,
    projects: [
      "Bazel (ij)",
      "Bazel-BSP (ij)",
      "Datalore (ij)",
      "gRPC-Java (ij)",
      "Hirschgarten (ij)",
      "Synthetic 1 project (ij)",
      "Synthetic 1000 project (ij)",
      "Synthetic 10000 project (ij)",
      "Synthetic 20000 project (ij)",
      "Synthetic 5000 project (ij)",
    ],
  }
})
</script>
