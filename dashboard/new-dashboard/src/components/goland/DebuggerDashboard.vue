<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    persistent-id="goland_debugger_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <template
      v-for="group in allGroups"
      :key="group.value"
    >
      <Divider :label="group.title" />
      <section>
        <GroupProjectsChart
          v-for="chart in group.charts"
          :key="chart.key"
          :label="`${group.prefix}: ${chart.label}`"
          :measure="chart.measure"
          :projects="group.projects"
          :value-unit="chart.valueUnit"
          :better-direction="chart.betterDirection"
          :description="chart.description"
        />
      </section>
    </template>

    <AdditionalMetrics :projects="allProjects" />
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import AdditionalMetrics from "./AdditionalMetrics.vue"
import type { ValueUnit } from "../common/chart"
import type { BetterDirection } from "../../shared/changeDetector/algorithm"

interface ChartDef {
  key: string
  label: string
  measure: string
  description: string
  valueUnit?: ValueUnit
  betterDirection?: BetterDirection
}

interface GroupDef {
  value: string
  title: string
  prefix: string
  projects: string[]
  charts: ChartDef[]
}

const debugProjects = ["river/debug", "trufflehog/debug"]
const allProjects = debugProjects

const debugCharts: ChartDef[] = [
  {
    key: "launch",
    label: "Launch Debug",
    measure: "debugRunConfiguration",
    valueUnit: "ms",
    description: "Time from starting the debug configuration to the first pause at a breakpoint.",
  },
  { key: "stepInto", label: "Step Into", measure: "debugStep_into", valueUnit: "ms", description: "Time from invoking Step Into until the debugger pauses again." },
  { key: "stepOut", label: "Step Out", measure: "debugStep_out", valueUnit: "ms", description: "Time from invoking Step Out until the debugger pauses again." },
  { key: "stepOver", label: "Step Over", measure: "debugStep_over", valueUnit: "ms", description: "Time from invoking Step Over until the debugger pauses again." },
]

const allGroups: GroupDef[] = [{ value: "debugActions", title: "Debug Actions", prefix: "Debug", projects: debugProjects, charts: debugCharts }]
</script>
