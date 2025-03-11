<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="fus_startup_dashboard"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :charts="charts"
  >
    <section>
      <div>
        <GroupProjectsChart
          v-for="chart in charts"
          :key="chart.definition.label"
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../../charts/DashboardCharts"
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Code analysis execution time", "Startup total duration"],
    measures: ["fusExecutionTime", "fusTotalDuration", "fusCodeVisibleInEditorDurationMs"],
    projects: ["idea/measureStartup"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
