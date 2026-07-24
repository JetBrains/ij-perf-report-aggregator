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
        :aliases="chart.aliases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

// SshOverIjentPerformanceTest reports through the Benchmark DSL: each measured operation runs under a distinct launch
// name (`<backend>.<operation>`), so the launch name lands in the `project` column and the timed block's duration lands
// in the generic `attempt.mean.ms` measure. Each chart compares the two backends selected via the `ijent.ssh.backend`
// registry key — `jetbrainsd` vs `sshj` — rendered as two aliased series.
//
// NOTE: the exact `project` string depends on the perf harness. It may be stored either as the bare launch name
// (`sshj.exec_warm`, as in IJentBenchmarskDashboard) or prefixed with the test id
// (`…SshOverIjentPerformanceTest.testSshjBackendPerformance - sshj.exec_warm`, as in the Rust/PhpStorm unit dashboards).
// Confirm against the first CI run's data (perfUnit tab / metrics.performance.json) and adjust the strings below if needed.
const BACKENDS = ["jetbrainsd", "sshj"] as const
const aliases = [...BACKENDS]

function backendProjects(operation: string): string[] {
  return BACKENDS.map((backend) => `${backend}.${operation}`)
}

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Warm Command Exec"],
    measures: ["attempt.mean.ms"],
    projects: backendProjects("exec_warm"),
    aliases,
  },
  {
    labels: ["Cold Bootstrap"],
    measures: ["attempt.mean.ms"],
    projects: backendProjects("bootstrap_cold"),
    aliases,
  },
  {
    labels: ["Warm SFTP Operations"],
    measures: ["attempt.mean.ms"],
    projects: backendProjects("sftp_warm"),
    aliases,
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
