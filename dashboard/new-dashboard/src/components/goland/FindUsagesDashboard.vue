<template>
  <DashboardPage
    :with-installer="false"
    db-name="perfintDev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="goland_findusages_dashboard"
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
          :better-direction="chart.betterDirection"
          :description="chart.description"
          :label="`${group.prefix}: ${chart.label}`"
          :measure="chart.measure"
          :projects="group.projects"
          :value-unit="chart.valueUnit"
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
const allProjects = findUsagesProjects

const findUsagesCharts: ChartDef[] = [
  { key: "total", label: "Total Execution Time", measure: "findUsages", valueUnit: "ms", description: "Time to find and show all usages in the popup." },
  { key: "firstUsage", label: "First Usage Time", measure: "findUsages_firstUsage", valueUnit: "ms", description: "Time until the first usage appears in the popup." },
  {
    key: "number",
    label: "Number of Usages",
    measure: "findUsages#number",
    valueUnit: "counter",
    betterDirection: "stable",
    description: "Count of usages found. Stable for a fixed query; a drop often signals a resolve regression.",
  },
]

const allGroups: GroupDef[] = [{ value: "findUsages", title: "Find Usages", prefix: "Find Usages", projects: findUsagesProjects, charts: findUsagesCharts }]
</script>
