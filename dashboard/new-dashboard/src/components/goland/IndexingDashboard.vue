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
          :label="`${group.prefix}: ${chart.label}`"
          :measure="chart.measure"
          :projects="group.projects"
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

interface ChartDef {
  key: string
  label: string
  measure: string
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
  { key: "indexingTime", label: "Indexing Time", measure: "indexingTimeWithoutPauses" },
  { key: "indexedFiles", label: "Indexed Files", measure: "numberOfIndexedFilesWritingIndexValue" },
  { key: "indexSize", label: "Index Size", measure: "indexSize" },
]

const processingCharts: ChartDef[] = [
  { key: "processingTime", label: "Processing Time", measure: "processingTime#Go" },
  { key: "processingSpeed", label: "Processing Speed", measure: "processingSpeedAvg#Go" },
]

// Lexing/parsing metrics are keyed by Language.getID() ("go", lowercase); processing metrics above
// are keyed by FileType.getName() ("Go"). The casing split is intentional — do not normalize it.
const parsingCharts: ChartDef[] = [
  { key: "parsingTime", label: "Parsing Time", measure: "parsingTime#go" },
  { key: "lexingTime", label: "Lexing Time", measure: "lexingTime#go" },
]

const scanningCharts: ChartDef[] = [{ key: "scanningTime", label: "Scanning Time", measure: "scanningTimeWithoutPauses" }]

const allGroups: GroupDef[] = [
  { value: "total", title: "Total Indexing", prefix: "Total", projects: indexingProjects, charts: indexingCharts },
  { value: "processing", title: "Processing", prefix: "Processing", projects: breakdownProjects, charts: processingCharts },
  { value: "parsing", title: "Parsing & Lexing", prefix: "Parsing", projects: breakdownProjects, charts: parsingCharts },
  { value: "scanning", title: "Scanning", prefix: "Scanning", projects: indexingProjects, charts: scanningCharts },
]

const indexingPipelineCharts: ChartDef[] = [
  { key: "dumbMode", label: "Dumb mode (ms)", measure: "dumbModeTimeWithPauses" },
  { key: "pausedTime", label: "Paused time (ms)", measure: "pausedTimeInIndexingOrScanning" },
  { key: "indexingRuns", label: "Indexing runs", measure: "numberOfRunsOfIndexing" },
  { key: "scanningRuns", label: "Scanning runs", measure: "numberOfRunsOfScannning" },
  { key: "scannedFiles", label: "Scanned files", measure: "numberOfScannedFiles" },
  { key: "indexedFilesTotal", label: "Indexed files (total)", measure: "numberOfIndexedFiles" },
  { key: "indexedNothing", label: "Indexed (nothing to write)", measure: "numberOfIndexedFilesWithNothingToWrite" },
  { key: "filesWithExtensions", label: "Files with extensions", measure: "numberOfFilesIndexedByExtensions" },
  { key: "filesWithoutExtensions", label: "Files without extensions", measure: "numberOfFilesIndexedWithoutExtensions" },
]

const writeActionCharts: ChartDef[] = [
  { key: "waCount", label: "Write action count", measure: "writeAction.count" },
  { key: "waWaitTotal", label: "Write action wait (ms)", measure: "writeAction.wait.ms" },
  { key: "waWaitMax", label: "Write action max wait (ms)", measure: "writeAction.max.wait.ms" },
  { key: "waWaitMedian", label: "Write action median wait (ms)", measure: "writeAction.median.wait.ms" },
]

const awtCharts: ChartDef[] = [{ key: "awtDispatch", label: "AWT dispatch total (ms)", measure: "AWTEventQueue.dispatchTimeTotal" }]
</script>
