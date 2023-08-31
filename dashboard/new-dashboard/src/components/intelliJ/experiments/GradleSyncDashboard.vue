<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="gradle_sync_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
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
    labels: ["Indexing time", "Scanning time", "Paused time", "Dumb mode with pauses", "Number of indexed files", "Memory 95th pctl", "Freed memory by GC"],
    measures: [
      "indexingTimeWithoutPauses",
      "scanningTimeWithoutPauses",
      "pausedTimeInIndexingOrScanning",
      "dumbModeTimeWithPauses",
      "numberOfIndexedFiles",
      "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
      "freedMemoryByGC",
    ],
    projects: [
      "indexing-with-gradle-sync-in-dumb",
      "indexing-with-gradle-sync-in-dumb-big-project",
      "indexing-with-gradle-sync-in-smart",
      "indexing-with-gradle-sync-in-smart-big-project",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
