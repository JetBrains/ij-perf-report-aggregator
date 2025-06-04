<template>
  <DashboardPage
    db-name="bazel"
    table="report"
    persistent-id="bazel_plugin_dashboard"
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
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsDeclaration = [
  ["Calculating external repository mapping"],
  ["Realizing language aspect files from templates"],
  ["resolve.project.time.ms", "Resolve project"],
  ["building.project.with.aspect.time.ms", "Building project with aspect"],
  ["mapping.to.internal.model.time.ms", "Mapping to internal model"],
  ["parsing.aspect.outputs.time.ms", "Parsing aspect outputs"],
  ["discovering.supported.external.rules.time.ms", "Discovering supported external rules"],
  ["select.targets.time.ms", "Select targets"],
  ["build.dependency.tree.time.ms"],
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
]

interface BazelCharts {
  labels: string[]
  measures: string[]
  projects: string[]
}

const chartsDeclaration: BazelCharts[] = metricsDeclaration.map((metricGroup) => {
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
      "Synthetic 5000 project (ij)",
      "Synthetic 10000 project (ij)",
      "Synthetic 20000 project (ij)",
    ],
  }
})
</script>
