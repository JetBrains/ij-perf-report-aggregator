<template>
  <DashboardPage
    db-name="perfintDev"
    table="kmt"
    persistent-id="kmt_performance_dashboard"
    initial-machine="Mac Cidr Performance"
    :initial-mode="MODES"
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
        :legend-formatter="legendFormatter"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"
import { legendFormatter, MODES } from "./KmtMeasurements"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: [["indexingTimeWithoutPauses", "scanningTimeWithoutPauses"]],
    projects: ["Wizard/indexing", "KotlinConf/indexing"],
  },
  {
    labels: ["KMP Setup"],
    measures: [["Progress: Setting up run configurations...", "Progress: Generating Xcode files…"]],
    projects: ["Wizard/setup", "KotlinConf/setup"],
  },
  {
    labels: ["Swift support"],
    measures: ["globalInspections"],
    projects: ["Wizard/inspection", "KotlinConf/inspection"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
