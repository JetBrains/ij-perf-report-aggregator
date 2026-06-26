<template>
  <DashboardPage
    db-name="perfintDev"
    :with-installer="false"
    table="goland"
    persistent-id="goland_indexing_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <template
      v-for="(row, rowIndex) in chartRows"
      :key="rowIndex"
    >
      <section
        v-if="row.length > 1"
        class="flex gap-x-6"
      >
        <div
          v-for="chart in row"
          :key="chart.key"
          class="flex-1 min-w-0"
        >
          <GroupProjectsChart
            :label="chart.label"
            :measure="chart.measure"
            :projects="chart.projects"
            :better-direction="chart.betterDirection"
            :description="chart.description"
          />
        </div>
      </section>
      <section v-else>
        <GroupProjectsChart
          :label="row[0].label"
          :measure="row[0].measure"
          :projects="row[0].projects"
          :better-direction="row[0].betterDirection"
          :description="row[0].description"
        />
      </section>
    </template>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import type { BetterDirection } from "../../shared/changeDetector/algorithm"

interface ChartDef {
  key: string
  label: string
  measure: string
  projects: string[]
  description?: string
  betterDirection?: BetterDirection
}

const indexingProjects = ["cockroach/indexing", "delve/indexing", "mattermost/indexing", "kubernetes/indexing", "flux/indexing", "istio/indexing"]
const breakdownProjects = ["kubernetes/indexing", "flux/indexing", "istio/indexing", "cockroach/indexing", "delve/indexing", "mattermost/indexing"]

const chartRows: ChartDef[][] = [
  [{ key: "indexingTime", label: "Indexing Time", measure: "indexingTimeWithoutPauses", projects: indexingProjects }],
  [
    { key: "indexedFiles", label: "Indexed Files", measure: "numberOfIndexedFilesWritingIndexValue", projects: indexingProjects, betterDirection: "stable" },
    { key: "indexSize", label: "Index Size", measure: "indexSize", projects: indexingProjects },
  ],
  [
    { key: "processingTime", label: "Processing Time", measure: "processingTime#Go", projects: breakdownProjects },
    { key: "processingSpeed", label: "Processing Speed", measure: "processingSpeedAvg#Go", projects: breakdownProjects, betterDirection: "higher" },
  ],
  [
    { key: "parsingTime", label: "Parsing Time", measure: "parsingTime#go", projects: breakdownProjects },
    { key: "lexingTime", label: "Lexing Time", measure: "lexingTime#go", projects: breakdownProjects },
  ],
  [{ key: "gcPause", label: "GC Pause, ms", measure: "gcPause", projects: indexingProjects }],
  [{ key: "freedMemoryByGC", label: "GC Memory Collected, Mb", measure: "freedMemoryByGC", projects: indexingProjects }],
  [{ key: "scanningTime", label: "Scanning Time", measure: "scanningTimeWithoutPauses", projects: indexingProjects }],
]
</script>
