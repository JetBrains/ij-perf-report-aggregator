<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    :with-installer="false"
    persistent-id="goland_findusages_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section
      v-for="chart in charts"
      :key="chart.key"
    >
      <GroupProjectsChart
        :label="chart.label"
        :measure="chart.measure"
        :projects="findUsagesProjects"
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

const findUsagesProjects = [
  "vault/backend(interface)",
  "vault/list(method)",
  "vault/path(struct)",
  "vault/string(method)",
  "vault/unmarshalJSON(method)",
  "pocketbase/write(method)",
  "pocketbase/start(method)",
  "pocketbase/open(method)",
  "pocketbase/file(struct)",
  "pocketbase/close(method)",
]

const charts: ChartDef[] = [
  { key: "total", label: "Total Execution Time", measure: "findUsages", description: "Time to find and show all usages in the popup." },
  { key: "firstUsage", label: "First Usage Time", measure: "findUsages_firstUsage", description: "Time until the first usage appears in the popup." },
  {
    key: "number",
    label: "Number of Usages",
    measure: "findUsages#number",
    betterDirection: "stable",
    description: "Count of usages found. Stable for a fixed query; a drop often signals a resolve regression.",
  },
]
</script>
