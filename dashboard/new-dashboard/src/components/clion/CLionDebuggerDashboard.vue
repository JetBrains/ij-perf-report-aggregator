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
          :measure="chart.measures"
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
  measures: string[]
}

interface GroupDef {
  value: string
  title: string
  prefix: string
  projects: string[]
  charts: ChartDef[]
}

const debugProjects = ["radler/fmtlib/debug/args-test/basic", "radler/luau/debug/Analyze.cpp", "radler/opencv/debug/test_arithm.cpp"]

const debugCharts: ChartDef[] = [
  { key: "launch", label: "Launch Debug", measures: ["fus_debug_session_initialized_ms"] },
  { key: "frame", label: "Frame Variables Computed", measures: ["fus_frame_variables_computed_ms"] },
  { key: "stepInto", label: "Step Into", measures: ["debugStep_into"] },
  { key: "stepOut", label: "Step Out", measures: ["debugStep_out"] },
  { key: "stepOver", label: "Step Over", measures: ["debugStep_over"] },
]

const allGroups: GroupDef[] = [{ value: "debugActions", title: "Debug Actions", prefix: "Debug", projects: debugProjects, charts: debugCharts }]
</script>
