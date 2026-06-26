<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    :with-installer="false"
    persistent-id="goland_debugger_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section
      v-for="chart in charts"
      :key="chart.key"
    >
      <GroupProjectsChart
        :label="chart.label"
        :measure="chart.measure"
        :projects="debugProjects"
        :better-direction="chart.betterDirection"
        :description="chart.description"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import type { BetterDirection } from "../../shared/changeDetector/algorithm"

interface ChartDef {
  key: string
  label: string
  measure: string
  description?: string
  betterDirection?: BetterDirection
}

const debugProjects = ["river/debug", "trufflehog/debug"]

const charts: ChartDef[] = [
  { key: "launch", label: "Launch Debug", measure: "debugRunConfiguration" },
  { key: "stepInto", label: "Step Into", measure: "debugStep_into" },
  { key: "stepOut", label: "Step Out", measure: "debugStep_out" },
  { key: "stepOver", label: "Step Over", measure: "debugStep_over" },
]
</script>
