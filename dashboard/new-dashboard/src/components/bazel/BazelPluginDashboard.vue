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
  ["resolve.project.time.ms"],
  ["max.used.memory.mb"],
  ["used.at.exit.mb"],
  ["building.project.with.aspect.time.ms"],
  ["mapping.to.internal.model.time.ms"],
  ["parsing.aspect.outputs.time.ms"],
  ["create.modules.time.ms"],
  ["reading.aspect.output.paths.time.ms"],
  ["fetching.all.possible.target.names.time.ms"],
  ["discovering.supported.external.rules.time.ms"],
  ["select.targets.time.ms"],
  ["libraries.from.jdeps.time.ms"],
  ["libraries.from.targets.and.deps.time.ms"],
  ["build.dependency.tree.time.ms"],
  ["build.reverse.sources.time.ms"],
  ["targets.as.libraries.time.ms"],
  ["create.ap.libraries.time.ms"],
  ["create.kotlin.stdlibs.time.ms"],
  ["save.invalid.target.labels.time.ms"],
  ["libraries.from.transitive.compile-time.jars.time.ms"],
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
      "Bazel",
      "Bazel-BSP",
      "Datalore",
      "gRPC-Java",
      "Hirschgarten",
      "Synthetic 1 project",
      "Synthetic 1000 project",
      "Synthetic 5000 project",
      "Synthetic 10000 project",
      "Synthetic 20000 project",
    ],
  }
})
</script>
