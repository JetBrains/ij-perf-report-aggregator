<template>
  <DashboardPage
    :with-installer="false"
    db-name="perfintDev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="goland_debugger_dashboard"
    table="goland"
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
        />
      </section>
    </template>

    <AdditionalMetrics :projects="allProjects" />
  </DashboardPage>
</template>

<script lang="ts" setup>
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import AdditionalMetrics from "./AdditionalMetrics.vue"

interface ChartDef {
  key: string
  label: string
  measure: string
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
  { key: "launch", label: "Launch Debug", measure: "debugRunConfiguration" },
  { key: "stepInto", label: "Step Into", measure: "debugStep_into" },
  { key: "stepOut", label: "Step Out", measure: "debugStep_out" },
  { key: "stepOver", label: "Step Over", measure: "debugStep_over" },
]

const allGroups: GroupDef[] = [{ value: "debugActions", title: "Debug Actions", prefix: "Debug", projects: debugProjects, charts: debugCharts }]
</script>
