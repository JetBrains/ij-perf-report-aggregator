<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    persistent-id="goland_completion_dashboard"
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

const typingProjects = ["act/typing"]
const basicProjects = ["caddy/completion/variable", "caddy/completion/type", "caddy/completion/return", "caddy/completion/interface", "caddy/completion/import"]
const smartProjects = ["permify/completion/method", "permify/completion/any"]
const allProjects = [...typingProjects, ...basicProjects, ...smartProjects]

const completionCharts: ChartDef[] = [
  { key: "duration90", label: "Duration 90p", measure: "fus_completion_duration_90p", description: "90th-percentile completion time (FUS) — the slow tail." },
  { key: "median", label: "Duration (median)", measure: "completion#median_value", description: "Median completion time — typical warm latency." },
  { key: "timeToShow90", label: "Time to show 90p", measure: "fus_time_to_show_90p", description: "90th-percentile time to first suggestion (FUS) — perceived lag." },
  {
    key: "timeToShowCold",
    label: "Time to show (cold)",
    measure: "completion_1#firstElementShown",
    valueUnit: "ms",
    description: "Time to first suggestion on the cold first run (span instrument).",
  },
  {
    key: "timeToShowMean",
    label: "Time to show (mean)",
    measure: "completion#firstElementShown#mean_value",
    description: "Mean time to first suggestion (span); includes the cold run.",
  },
  {
    key: "resolveCold",
    label: "Resolve (cold)",
    measure: "performCompletion_1",
    description: "Cold-run cost of computing suggestions: resolve, type inference, stub-index lookups. A sensitive resolve-regression signal.",
  },
  {
    key: "items",
    label: "Result items",
    measure: "completion#number#mean_value",
    valueUnit: "counter",
    betterDirection: "stable",
    description: "Mean suggestion count. Any change is suspicious — a drop loses suggestions, a rise adds noise.",
  },
  {
    key: "stddev",
    label: "Standard deviation",
    measure: "completion#standard_deviation",
    description: "Spread of completion time across invocations. Rising values mean inconsistent latency.",
  },
]

const allGroups = [
  {
    value: "typing",
    title: "Typing",
    prefix: "Typing",
    projects: typingProjects,
    charts: [{ key: "totalTime", label: "Total Time", measure: "typing", description: "Total time to type the sample (act)." }],
  },
  { value: "basic", title: "Basic Completion", prefix: "Basic", projects: basicProjects, charts: completionCharts },
  { value: "smart", title: "Smart Completion", prefix: "Smart", projects: smartProjects, charts: completionCharts },
]
</script>
