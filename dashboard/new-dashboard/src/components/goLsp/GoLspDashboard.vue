<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="goLspDashboard"
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

const projects = ["kubernetes-goland-server/codeSearch", "kubernetes-gopls/codeSearch"]
const aliases = ["GoLand server", "gopls"]

const requests = ["definition", "references", "implementation", "workspaceSymbol"]
const scenarios = ["first", "warm", "edit.nearbyBody", "edit.nearbyFunction", "edit.coreLibraryFunction"]

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["ready"],
    measures: ["ready"],
    projects,
    aliases,
  },
  ...requests.map((request) => ({
    labels: scenarios.map((scenario) => `${request} - ${scenario}`),
    measures: scenarios.map((scenario) => `${request}.${scenario}`),
    projects,
    aliases,
  })),
]

const charts = combineCharts(chartsDeclaration)
</script>
