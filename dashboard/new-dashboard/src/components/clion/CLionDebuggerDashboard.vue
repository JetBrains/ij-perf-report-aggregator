<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_debugger_dashboard"
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
        />
      </section>
    </template>
  </DashboardPage>
</template>

<script lang="ts" setup>
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

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

const debugProjects = ["radler/fmtlib/debug/args-test/basic"]

const debugCharts: ChartDef[] = [
  { key: "stepInto", label: "Step Into", measure: "debugStep_into#mean_value" },
  { key: "stepOut", label: "Step Out", measure: "debugStep_out#mean_value" },
  { key: "stepOver", label: "Step Over", measure: "debugStep_over#mean_value" },
]

const allGroups: GroupDef[] = [{ value: "debugActions", title: "Debug Actions", prefix: "Debug", projects: debugProjects, charts: debugCharts }]
</script>
