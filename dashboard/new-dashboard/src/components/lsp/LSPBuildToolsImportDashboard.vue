<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="lspBuildToolsImportDashboard"
    initial-machine="linux-blade-hetzner"
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
        :aliases="chart.aliases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"

const TEST_CLASS = "com.jetbrains.ls.performanceTest.buildToolsImport.MavenImportPerformanceTest"

const projects = [`${TEST_CLASS}.logback`, `${TEST_CLASS}.mybatisDynamicSql`]
const aliases = ["logback", "mybatis-dynamic-sql"]

const phases = [
  { phase: "total", label: "Maven import - total", description: "initialize, to intellij/ready-for-test" },
  { phase: "import", label: "Maven import - import", description: "the Maven tool run, importLog started to succeeded" },
  { phase: "indexing", label: "Maven import - indexing", description: "importLog succeeded, to the FINISHED import state" },
]

const chartsDeclaration: ChartDefinition[] = phases.map(({ phase, label, description }) => ({
  labels: [label],
  measures: [`${phase}.ms`],
  descriptions: [description],
  projects,
  aliases,
}))

const charts = combineCharts(chartsDeclaration)
</script>
