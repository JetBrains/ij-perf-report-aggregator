<template>
  <DashboardPage
    :with-installer="false"
    db-name="perfintDev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="goland_indexing_dashboard"
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
          :better-direction="chart.betterDirection"
          :description="chart.description"
          :label="`${group.prefix}: ${chart.label}`"
          :measure="chart.measure"
          :projects="group.projects"
          :value-unit="chart.valueUnit"
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
import type { ValueUnit } from "../common/chart"
import type { BetterDirection } from "../../shared/changeDetector/algorithm"

interface ChartDef {
  key: string
  label: string
  measure: string
  description: string
  valueUnit?: ValueUnit
  betterDirection?: BetterDirection
}

interface GroupDef {
  value: string
  title: string
  prefix: string
  projects: string[]
  charts: ChartDef[]
}

const indexingProjects = ["cockroach/indexing", "delve/indexing", "mattermost/indexing", "kubernetes/indexing", "flux/indexing", "istio/indexing"]
const breakdownProjects = ["kubernetes/indexing", "flux/indexing", "istio/indexing", "cockroach/indexing", "delve/indexing", "mattermost/indexing"]
const allProjects = indexingProjects

const indexingCharts: ChartDef[] = [
  { key: "indexingTime", label: "Indexing Time", measure: "indexingTimeWithoutPauses", valueUnit: "ms", description: "Time to build indexes, excluding paused intervals." },
  {
    key: "indexedFiles",
    label: "Indexed Files",
    measure: "numberOfIndexedFilesWritingIndexValue",
    valueUnit: "counter",
    betterDirection: "stable",
    description: "Files that produced index data this run. Should stay flat for a fixed project; a change hints at an indexing-scope shift.",
  },
  { key: "indexSize", label: "Index Size", measure: "indexSize", description: "Total size of indexes written on disk after indexing." },
]

const processingCharts: ChartDef[] = [
  { key: "processingTime", label: "Processing Time", measure: "processingTime#Go", valueUnit: "ms", description: "CPU time spent indexing Go files." },
  {
    key: "processingSpeed",
    label: "Processing Speed",
    measure: "processingSpeedAvg#Go",
    betterDirection: "higher",
    description: "Average indexing throughput for Go files (kB/s); higher is better.",
  },
]

const parsingCharts: ChartDef[] = [
  { key: "parsingTime", label: "Parsing Time", measure: "parsingTime#go", valueUnit: "ms", description: "Time the parser spends building PSI for Go files during indexing." },
  { key: "lexingTime", label: "Lexing Time", measure: "lexingTime#go", valueUnit: "ms", description: "Time the lexer spends tokenizing Go files during indexing." },
]

const scanningCharts: ChartDef[] = [
  {
    key: "scanningTime",
    label: "Scanning Time",
    measure: "scanningTimeWithoutPauses",
    valueUnit: "ms",
    description: "Time to scan files for changes before indexing, excluding pauses.",
  },
]

const allGroups: GroupDef[] = [
  { value: "total", title: "Total Indexing", prefix: "Total", projects: indexingProjects, charts: indexingCharts },
  { value: "processing", title: "Processing", prefix: "Processing", projects: breakdownProjects, charts: processingCharts },
  { value: "parsing", title: "Parsing & Lexing", prefix: "Parsing", projects: breakdownProjects, charts: parsingCharts },
  { value: "scanning", title: "Scanning", prefix: "Scanning", projects: indexingProjects, charts: scanningCharts },
]
</script>
