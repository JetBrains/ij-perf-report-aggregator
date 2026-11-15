<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="goLspDashboard"
    initial-machine="linux-blade-hetzner"
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
        :aliases="chart.aliases"
      />
    </section>
    <!-- gopls is not a JVM, so only the GoLand server reports these. -->
    <Divider label="GoLand server runtime" />
    <section>
      <GroupProjectsChart
        v-for="chart in runtimeCharts"
        :key="chart.measure"
        :label="chart.label"
        :measure="chart.measure"
        :projects="jvmProjects"
        :aliases="jvmAliases"
      />
    </section>
    <ChartAccordion :lazy="true">
      <AccordionPanel value="0">
        <AccordionHeader>Additional GoLand server metrics</AccordionHeader>
        <AccordionContent>
          <section>
            <GroupProjectsChart
              v-for="chart in additionalCharts"
              :key="chart.measure"
              :label="chart.label"
              :measure="chart.measure"
              :projects="jvmProjects"
              :aliases="jvmAliases"
            />
          </section>
        </AccordionContent>
      </AccordionPanel>
    </ChartAccordion>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import ChartAccordion from "../charts/ChartAccordion.vue"
import Divider from "../common/Divider.vue"
import AccordionPanel from "openvue/accordionpanel"
import AccordionHeader from "openvue/accordionheader"
import AccordionContent from "openvue/accordioncontent"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import { GO_LSP_REQUESTS, GO_LSP_SCENARIOS } from "./goLspMetrics"

const projects = ["kubernetes-goland-server/codeSearch", "kubernetes-gopls/codeSearch"]
const aliases = ["GoLand server", "gopls"]

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["ready"],
    measures: ["ready"],
    projects,
    aliases,
  },
  ...GO_LSP_REQUESTS.map((request) => ({
    labels: GO_LSP_SCENARIOS.map((scenario) => `${request} - ${scenario}`),
    measures: GO_LSP_SCENARIOS.map((scenario) => `${request}.${scenario}`),
    projects,
    aliases,
  })),
]

const charts = combineCharts(chartsDeclaration)

const jvmProjects = [projects[0]]
const jvmAliases = [aliases[0]]

// The JVM metrics of the whole server run, under the names of the GoLand IDE dashboards (goland/AdditionalMetrics.vue).
const runtimeCharts = [
  { label: "Allocated", measure: "JVM.totalMegabytesAllocated" },
  { label: "Total CPU", measure: "JVM.totalCpuTimeMs" },
  { label: "GC collection time", measure: "JVM.GC.collectionTimesMs" },
  { label: "GC collections", measure: "JVM.GC.collections" },
  { label: "Used heap, peak", measure: "JVM.usedHeapMegabytes" },
]

const additionalCharts = [
  { label: "Max heap", measure: "JVM.maxHeapMegabytes" },
  { label: "Committed heap", measure: "JVM.committedHeapMegabytes" },
  { label: "Native used", measure: "JVM.usedNativeMegabytes" },
  { label: "Direct byte buffers", measure: "JVM.totalDirectByteBuffersMegabytes" },
  { label: "Max threads", measure: "JVM.maxThreadCount" },
  { label: "New threads", measure: "JVM.newThreadsCount" },
  { label: "Time to safepoints", measure: "JVM.totalTimeToSafepointsMs" },
  { label: "Time at safepoints", measure: "JVM.totalTimeAtSafepointsMs" },
  { label: "Safepoint count", measure: "JVM.totalSafepointCount" },
  { label: "OS load average", measure: "OS.loadAverage" },
]
</script>
