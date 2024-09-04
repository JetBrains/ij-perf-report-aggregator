<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="popups_dashboard_dev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
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
    labels: ["EditorContextMenu"],
    measures: ["popupShown#EditorContextMenu"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["ProjectViewContextMenu"],
    measures: ["popupShown#ProjectViewContextMenu"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["ProjectWidget"],
    measures: ["popupShown#ProjectWidget"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["RunConfigurations"],
    measures: ["popupShown#RunConfigurations"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsLogBranchFilter"],
    measures: ["popupShown#VcsLogBranchFilter"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsLogDateFilter"],
    measures: ["popupShown#VcsLogDateFilter"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsLogPathFilter"],
    measures: ["popupShown#VcsLogPathFilter"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsLogUserFilter"],
    measures: ["popupShown#VcsLogUserFilter"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsWidget"],
    measures: [["popupShown#VcsWidget", "afterShow#GitBranchesTreePopup"]],
    projects: ["popups-performance-test/test-popups"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
