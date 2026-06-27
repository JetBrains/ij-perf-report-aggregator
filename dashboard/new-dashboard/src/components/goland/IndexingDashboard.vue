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

    <ChartAccordion :lazy="true">
      <AccordionPanel value="0">
        <AccordionHeader>Indexing details</AccordionHeader>
        <AccordionContent>
          <Divider label="Indexing pipeline" />
          <section>
            <GroupProjectsChart
              v-for="chart in indexingPipelineCharts"
              :key="chart.key"
              :label="chart.label"
              :measure="chart.measure"
              :projects="allProjects"
              :value-unit="chart.valueUnit"
              :description="chart.description"
            />
          </section>

          <Divider label="Write actions" />
          <section>
            <GroupProjectsChart
              v-for="chart in writeActionCharts"
              :key="chart.key"
              :label="chart.label"
              :measure="chart.measure"
              :projects="allProjects"
              :value-unit="chart.valueUnit"
              :description="chart.description"
            />
          </section>

          <Divider label="AWT" />
          <section>
            <GroupProjectsChart
              v-for="chart in awtCharts"
              :key="chart.key"
              :label="chart.label"
              :measure="chart.measure"
              :projects="allProjects"
              :value-unit="chart.valueUnit"
              :description="chart.description"
            />
          </section>
        </AccordionContent>
      </AccordionPanel>
    </ChartAccordion>

    <AdditionalMetrics :projects="allProjects" />
  </DashboardPage>
</template>

<script lang="ts" setup>
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import ChartAccordion from "../charts/ChartAccordion.vue"
import AdditionalMetrics from "./AdditionalMetrics.vue"
import AccordionPanel from "primevue/accordionpanel"
import AccordionHeader from "primevue/accordionheader"
import AccordionContent from "primevue/accordioncontent"
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

const indexingPipelineCharts: ChartDef[] = [
  { key: "dumbMode", label: "Dumb mode (ms)", measure: "dumbModeTimeWithPauses", valueUnit: "ms", description: "Total time in dumb mode (indexes not ready), including pauses." },
  { key: "pausedTime", label: "Paused time (ms)", measure: "pausedTimeInIndexingOrScanning", valueUnit: "ms", description: "Time indexing/scanning was paused (e.g., for GC or UI)." },
  { key: "indexingRuns", label: "Indexing runs", measure: "numberOfRunsOfIndexing", valueUnit: "counter", description: "Number of indexing passes. Higher = more incremental re-indexing." },
  { key: "scanningRuns", label: "Scanning runs", measure: "numberOfRunsOfScannning", valueUnit: "counter", description: "Number of file system scanning passes." },
  { key: "scannedFiles", label: "Scanned files", measure: "numberOfScannedFiles", valueUnit: "counter", description: "Total files scanned. Versus indexed files shows filter efficiency." },
  { key: "indexedFilesTotal", label: "Indexed files (total)", measure: "numberOfIndexedFiles", valueUnit: "counter", description: "Total files that went through indexing." },
  { key: "indexedNothing", label: "Indexed (nothing to write)", measure: "numberOfIndexedFilesWithNothingToWrite", valueUnit: "counter", description: "Files indexed but producing no index data. High = wasted indexing work." },
  { key: "filesWithExtensions", label: "Files with extensions", measure: "numberOfFilesIndexedByExtensions", valueUnit: "counter", description: "Files recognized by extension-based file type detection." },
  { key: "filesWithoutExtensions", label: "Files without extensions", measure: "numberOfFilesIndexedWithoutExtensions", valueUnit: "counter", description: "Files indexed despite having no recognized extension." },
]

const writeActionCharts: ChartDef[] = [
  { key: "waCount", label: "Write action count", measure: "writeAction.count", valueUnit: "counter", description: "Number of write actions executed during indexing." },
  { key: "waWaitTotal", label: "Write action wait (ms)", measure: "writeAction.wait.ms", valueUnit: "ms", description: "Total time spent waiting for write actions." },
  { key: "waWaitMax", label: "Write action max wait (ms)", measure: "writeAction.max.wait.ms", valueUnit: "ms", description: "Longest single write action wait. Spikes indicate blocking." },
  { key: "waWaitMedian", label: "Write action median wait (ms)", measure: "writeAction.median.wait.ms", valueUnit: "ms", description: "Median write action wait time. Typical contention level." },
]

const awtCharts: ChartDef[] = [
  { key: "awtDispatch", label: "AWT dispatch total (ms)", measure: "AWTEventQueue.dispatchTimeTotal", valueUnit: "ms", description: "Total AWT event queue dispatch time. High values indicate UI thread contention." },
]
</script>
