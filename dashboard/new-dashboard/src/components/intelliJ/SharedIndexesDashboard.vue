<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="shared_indexing_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
      />
    </section>
  </DashboardPage>
  >
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing", "Scanning", "Number od Indexed Files", "Number of Files indexed by Shared Indexes"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "numberOfIndexedFiles", "numberOfFilesIndexedByExtensions"],
    projects: ["project-shared-indexes-downloading-community", "project-shared-indexes-downloading-intellij"],
    aliases: ["community", "intellij"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
