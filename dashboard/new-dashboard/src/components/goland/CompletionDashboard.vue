<template>
  <DashboardPage
    :with-installer="false"
    db-name="perfintDev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="goland_completion_dashboard"
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

const typingProjects = ["act/typing"]
const basicProjects = ["caddy/completion/variable", "caddy/completion/type", "caddy/completion/return", "caddy/completion/interface", "caddy/completion/import"]
const smartProjects = ["permify/completion/method", "permify/completion/any"]
const allProjects = [...typingProjects, ...basicProjects, ...smartProjects]

const completionCharts: ChartDef[] = [
  { key: "duration90", label: "Duration 90p", measure: "fus_completion_duration_90p" },
  { key: "median", label: "Duration (median)", measure: "completion#median_value" },
  { key: "timeToShow90", label: "Time to show 90p", measure: "fus_time_to_show_90p" },
  { key: "timeToShowCold", label: "Time to show (cold)", measure: "completion_1#firstElementShown" },
  { key: "timeToShowMean", label: "Time to show (mean)", measure: "completion#firstElementShown#mean_value" },
  { key: "resolveCold", label: "Resolve (cold)", measure: "performCompletion_1" },
  { key: "items", label: "Result items", measure: "completion#number#mean_value" },
  { key: "stddev", label: "Standard deviation", measure: "completion#standard_deviation" },
]

const allGroups = [
  {
    value: "typing",
    title: "Typing",
    prefix: "Typing",
    projects: typingProjects,
    charts: [{ key: "totalTime", label: "Total Time", measure: "typing" }],
  },
  { value: "basic", title: "Basic Completion", prefix: "Basic", projects: basicProjects, charts: completionCharts },
  { value: "smart", title: "Smart Completion", prefix: "Smart", projects: smartProjects, charts: completionCharts },
]
</script>
