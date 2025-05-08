<template>
  <DashboardPage
    db-name="perfintDev"
    table="kmt"
    persistent-id="kmt_performance_dashboard"
    initial-machine="Mac Cidr Performance"
    :initial-mode="modes"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsWithClientChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Startup"],
    measures: ["totalOpeningTime/timeFromAppStartTillAnalysisFinished"],
    projects: ["KMT_Basic/measureStartup"],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["KMT_Basic/inspection/KMT_Basic"],
  },
  {
    labels: ["Indexing"],
    measures: [["indexingTimeWithoutPauses", "scanningTimeWithoutPauses"]],
    projects: ["KMT_Basic/indexingKMT_Basic"],
  },
]

const charts = combineCharts(chartsDeclaration)
const modes = ["intellij-idea", "android-studio"]
</script>
