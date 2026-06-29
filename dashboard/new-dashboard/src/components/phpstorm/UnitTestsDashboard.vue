<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="phpstorm_unit_tests_dashboard"
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
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import { phpStormUnitTestGroups } from "./PhpStormTestCases"

const chartsDeclaration: ChartDefinition[] = phpStormUnitTestGroups.map((group) => ({
  labels: [group.label],
  measures: ["attempt.mean.ms"],
  projects: group.projects,
}))

const charts = combineCharts(chartsDeclaration)
</script>
