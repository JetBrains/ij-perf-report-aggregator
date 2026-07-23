<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_ssh_performance_dashboard"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :with-installer="false"
    :charts="charts"
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
import { computed } from "vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const projects = [
  "com.intellij.ssh.integration.tests.exec.SshOverIjentPerformanceTest",
]

const charts = computed(() => {
  const chartsDeclaration: ChartDefinition[] = [
    {
      labels: ["Warm Command Exec (ms)"],
      measures: ["ijent.ssh.exec_warm.ms"],
      projects,
    },
    {
      labels: ["Cold Bootstrap (ms)"],
      measures: ["ijent.ssh.bootstrap_cold.ms"],
      projects,
    },
    {
      labels: ["Warm SFTP Operations (ms)"],
      measures: ["ijent.ssh.sftp_warm.ms"],
      projects,
    },
  ]

  return combineCharts(chartsDeclaration)
})
</script>
