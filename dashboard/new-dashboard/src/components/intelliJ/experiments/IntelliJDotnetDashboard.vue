<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="monorepo_dashboard"
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
    labels: ["Git Status"],
    measures: ["vfs_initial_refresh"],
    projects: ["intellij_sources/vfsRefresh/git-status", "monorepo/vfsRefresh/git-status"],
  },
  {
    labels: ["Indexing time", "Scanning time", "Dumb mode with pauses", "Number of indexed files"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses", "dumbModeTimeWithPauses", "numberOfIndexedFiles"],
    projects: ["intellij_sources/indexing", "monorepo/indexing"],
  },
  {
    labels: ["Compilation time", "Freed memory by GC"],
    measures: ["build_compilation_duration", "freedMemoryByGC"],
    projects: ["intellij_sources/rebuild", "monorepo/rebuild"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
