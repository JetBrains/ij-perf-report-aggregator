<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    :with-installer="false"
    persistent-id="goland_dfa_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section
      v-for="chart in charts"
      :key="chart.key"
    >
      <GroupProjectsChart
        :label="chart.label"
        :measure="chart.measure"
        :projects="dfaProjects"
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

const charts: ChartDef[] = [
  { key: "globalInspections", label: "Global Inspections Time", measure: "globalInspections" },
  { key: "totalTime", label: "General Total Time", measure: "go.dfa.general.total.time.ms" },
  { key: "avgTime", label: "General Average Time", measure: "go.dfa.general.avg.time.ms" },
  { key: "avgWithoutSummaryLoad", label: "General Average Time Without Summary Load Time", measure: "go.dfa.general.avg.without.summary.load.time.ms" },
  { key: "gistsCount", label: "General Computed File Gists Count", measure: "go.dfa.general.computed.file.gists.count", betterDirection: "stable" },
  { key: "filesCount", label: "General Files Count", measure: "go.dfa.general.files.count", betterDirection: "stable" },
  { key: "functionsCount", label: "General Functions Count", measure: "go.dfa.general.functions.count", betterDirection: "stable" },
]
</script>
