<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_startup_dashboard_dev"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Startup FUS total duration"],
    measures: ["startup/fusTotalDuration"],
    projects: ["idea/measureStartup", "idea/measureStartup/embeddedClient"],
  },
  {
    labels: ["Code Analysis Daemon FUS execution time"],
    measures: ["startup/fusTotalDuration"],
    projects: ["idea/measureStartup", "idea/measureStartup/embeddedClient"],
  },
  {
    labels: ["Reopen Project FUS code visible in editor duration"],
    measures: ["reopenProjectPerformance/fusCodeVisibleInEditorDurationMs"],
    projects: ["idea/measureStartup", "idea/measureStartup/embeddedClient"],
  },
  {
    labels: ["Reopen Project FUS first UI shown"],
    measures: ["reopenProjectPerformance/fusFirstUIShowsMs"],
    projects: ["idea/measureStartup", "idea/measureStartup/embeddedClient"],
  },
  {
    labels: ["Reopen Project FUS frame became interactive"],
    measures: ["reopenProjectPerformance/fusFrameBecameInteractiveMs"],
    projects: ["idea/measureStartup", "idea/measureStartup/embeddedClient"],
  },
  {
    labels: ["Reopen Project FUS frame became visible"],
    measures: ["reopenProjectPerformance/fusFrameBecameVisibleMs"],
    projects: ["idea/measureStartup", "idea/measureStartup/embeddedClient"],
  },
  {
    labels: ["Total Opening Time"],
    measures: ["totalOpeningTime/timeFromAppStartTillAnalysisFinished"],
    projects: ["idea/measureStartup", "idea/measureStartup/embeddedClient"],
  },
  {
    labels: ["First file opening time"],
    measures: ["openFile"],
    projects: ["idea/warmup/embeddedClient", "idea/warmup"],
  },
  {
    labels: ["Reopening file time after ide restart"],
    measures: ["reopenFileAfterIdeRestart"],
    projects: ["idea/measureStartup/embeddedClient", "idea/measureStartup"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
