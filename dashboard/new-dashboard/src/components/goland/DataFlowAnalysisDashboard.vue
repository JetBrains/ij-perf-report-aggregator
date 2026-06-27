<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    persistent-id="goland_dfa_dashboard"
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

const dfaProjects = [
  "kubernetes/kubernetes/dfa",
  "cockroach/cockroach/dfa",
  "mattermost-server/mattermost-server/dfa",
  "volcano/volcano/dfa",
  "tempo/tempo/dfa",
  "rclone/rclone/dfa",
  "milvus/milvus/dfa",
  "k8sDevice/k8sDevice/dfa",
]
const allProjects = dfaProjects

const inspectionCharts: ChartDef[] = [
  {
    key: "globalInspections",
    label: "Global Inspections Time",
    measure: "globalInspections",
    valueUnit: "ms",
    description: "Batch Inspect Code over the whole project — baseline cost before DFA analysis.",
  },
]

const generalCharts: ChartDef[] = [
  {
    key: "totalTime",
    label: "General Total Time",
    measure: "go.dfa.general.total.time.ms",
    valueUnit: "ms",
    description: "Total wall-clock time for general DFA analysis across all files.",
  },
  {
    key: "avgTime",
    label: "General Average Time",
    measure: "go.dfa.general.avg.time.ms",
    valueUnit: "ms",
    description: "Average DFA analysis time per file, including summary loading.",
  },
  {
    key: "avgWithoutSummaryLoad",
    label: "General Average Time Without Summary Load",
    measure: "go.dfa.general.avg.without.summary.load.time.ms",
    valueUnit: "ms",
    description: "Average DFA analysis time per file excluding summary load — isolates pure analysis cost.",
  },
  {
    key: "gistsCount",
    label: "General Computed File Gists Count",
    measure: "go.dfa.general.computed.file.gists.count",
    valueUnit: "counter",
    betterDirection: "stable",
    description: "Number of file gists computed. Stable for a fixed project; changes indicate analysis scope shifts.",
  },
  {
    key: "filesCount",
    label: "General Files Count",
    measure: "go.dfa.general.files.count",
    valueUnit: "counter",
    betterDirection: "stable",
    description: "Total files analyzed. Should remain constant for a fixed project.",
  },
  {
    key: "functionsCount",
    label: "General Functions Count",
    measure: "go.dfa.general.functions.count",
    valueUnit: "counter",
    betterDirection: "stable",
    description: "Total functions analyzed. Should remain constant for a fixed project.",
  },
]

const allGroups: GroupDef[] = [
  { value: "inspections", title: "Inspections", prefix: "Inspections", projects: dfaProjects, charts: inspectionCharts },
  { value: "general", title: "General DFA Metrics", prefix: "General", projects: dfaProjects, charts: generalCharts },
]
</script>
