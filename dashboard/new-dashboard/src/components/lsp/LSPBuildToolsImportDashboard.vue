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

const TEST_PACKAGE = "com.jetbrains.ls.performanceTest.buildToolsImport"

const tools = [
  {
    tool: "Maven",
    projects: [`${TEST_PACKAGE}.MavenImportPerformanceTest.logback`, `${TEST_PACKAGE}.MavenImportPerformanceTest.mybatisDynamicSql`],
    aliases: ["logback", "mybatis-dynamic-sql"],
  },
  {
    tool: "Gradle",
    projects: [`${TEST_PACKAGE}.GradleImportPerformanceTest.ehcache3`, `${TEST_PACKAGE}.GradleImportPerformanceTest.springAmqp`],
    aliases: ["ehcache3", "spring-amqp"],
  },
]

const phases = [
  { phase: "total", label: "total", description: "initialize, to intellij/ready-for-test" },
  { phase: "import", label: "import", description: "the build tool run, importLog started to succeeded" },
  { phase: "indexing", label: "indexing", description: "importLog succeeded, to the FINISHED import state" },
]

const chartsDeclaration: ChartDefinition[] = tools.map(({ tool, projects, aliases }) => ({
  labels: phases.map(({ label }) => `${tool} import - ${label}`),
  measures: phases.map(({ phase }) => `${phase}.ms`),
  descriptions: phases.map(({ description }) => description),
  projects,
  aliases,
}))

const charts = combineCharts(chartsDeclaration)
</script>
