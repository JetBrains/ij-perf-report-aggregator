<template>
  <DashboardPage
    :with-installer="false"
    db-name="perfintDev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="goland_dfa_dashboard"
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
          :description="chart.description"
          :label="`${group.prefix}: ${chart.label}`"
          :measure="chart.measure"
          :projects="group.projects"
        />
      </section>
      <ChartAccordion
        v-if="group.jvmCharts"
        :lazy="true"
      >
        <AccordionPanel value="0">
          <AccordionHeader>{{ group.title }}: JVM cost</AccordionHeader>
          <AccordionContent>
            <section>
              <GroupProjectsChart
                v-for="chart in group.jvmCharts"
                :key="chart.key"
                :label="`${group.prefix}: ${chart.label}`"
                :measure="chart.measure"
                :projects="group.projects"
              />
            </section>
          </AccordionContent>
        </AccordionPanel>
      </ChartAccordion>
    </template>

    <PerProjectBreakdown
      title="DFA stages per project"
      :projects="dfaProjects"
      :charts="stageBreakdown"
    />

    <AdditionalMetrics :projects="allProjects" />
  </DashboardPage>
</template>

<script lang="ts" setup>
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import ChartAccordion from "../charts/ChartAccordion.vue"
import AccordionPanel from "openvue/accordionpanel"
import AccordionHeader from "openvue/accordionheader"
import AccordionContent from "openvue/accordioncontent"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import AdditionalMetrics from "./AdditionalMetrics.vue"
import PerProjectBreakdown, { BreakdownChart } from "./PerProjectBreakdown.vue"
import { jvmCostCharts } from "./jvmCostCharts"

interface ChartDef {
  key: string
  label: string
  measure: string
  // Only set when a chart needs a description that differs from the central metricsDescription entry.
  description?: string
}

interface GroupDef {
  value: string
  title: string
  prefix: string
  projects: string[]
  charts: ChartDef[]
  jvmCharts?: ChartDef[]
}

const dfaProjects = [
  "kubernetes/kubernetes/dfa",
  "cockroach/cockroach/dfa",
  "mattermost-server/mattermost-server/dfa",
  "volcano/volcano/dfa",
  "tempo/tempo/dfa",
  "rclone/rclone/dfa",
  "milvus/milvus/dfa",
  "k8sDevice/k8sDevice/dfa",
]
const allProjects = dfaProjects

const inspectionCharts: ChartDef[] = [
  {
    key: "globalInspections",
    label: "Global Inspections Time",
    measure: "globalInspections",
    description: "Batch Inspect Code over the whole project — baseline cost before DFA analysis.",
  },
]

const generalCharts: ChartDef[] = [
  { key: "totalTime", label: "General Total Time", measure: "go.dfa.general.total.time.ms" },
  { key: "avgTime", label: "General Average Time", measure: "go.dfa.general.avg.time.ms" },
  { key: "avgWithoutSummaryLoad", label: "General Average Time Without Summary Load", measure: "go.dfa.general.avg.without.summary.load.time.ms" },
  { key: "gistsCount", label: "General Computed File Gists Count", measure: "go.dfa.general.computed.file.gists.count" },
  { key: "filesCount", label: "General Files Count", measure: "go.dfa.general.files.count" },
  { key: "functionsCount", label: "General Functions Count", measure: "go.dfa.general.functions.count" },
  { key: "resolveIssues", label: "Resolve Issues", measure: "go.dfa.general.resolve.issue.count" },
]

const stageBreakdown: BreakdownChart[] = [
  {
    key: "stages",
    label: "Stage time",
    measures: ["go.dfa.granular.total.ir.time.ms", "go.dfa.solver.total.time.ms", "go.dfa.granular.interpretation.time.ms", "go.dfa.granular.resolve.time.ms"],
  },
]

const allGroups: GroupDef[] = [
  {
    value: "inspections",
    title: "Inspections",
    prefix: "Inspections",
    projects: dfaProjects,
    charts: inspectionCharts,
    jvmCharts: jvmCostCharts("globalInspections", "Global Inspections", "globalInspections"),
  },
  { value: "general", title: "General DFA Metrics", prefix: "General", projects: dfaProjects, charts: generalCharts },
]
</script>
