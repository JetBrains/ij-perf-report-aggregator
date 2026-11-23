<template>
  <DashboardPage
    :with-installer="false"
    db-name="perfintDev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="goland_code_analyzes_dashboard"
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
      title="File analysis per file"
      :projects="fileAnalysisProjects"
      :charts="fileAnalysisBreakdown"
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

const fileAnalysisProjects = [
  "istio/localInspection/adsc.go",
  "minotaur/localInspection/server.go",
  "caddy/localInspection/replacer.go",
  "cockroach/localInspection/tochar.go",
  "flux/localInspection/sourcesecret.go",
  "kubernetes/localInspection/client_test.go",
]

const globalInspectionProjects = [
  "delve/inspection",
  "kubernetes/inspection",
  "caddy/inspection",
  "cockroach/inspection",
  "k8sDevice/inspection",
  "mattermost-server/inspection",
  "milvus/inspection",
  "rclone/inspection",
  "tempo/inspection",
  "volcano/inspection",
]

const golangciLintLocalProjects = ["volcano-golinter-local-without-linter/localInspection/scheduler.go", "volcano-golinter-local-with-linter/localInspection/scheduler.go"]

const golangciLintGlobalProjects = ["volcano-golinter-global-with-linter/inspection", "volcano-golinter-global-without-linter/inspection"]

const singleInspectionProjects = [
  "milvus/GoCommentLeadingSpace",
  "milvus/GoCommentStart",
  "milvus/GoErrorStringFormat",
  "milvus/GoExportedElementShouldHaveComment",
  "milvus/GoExportedOwnDeclaration",
  "milvus/GoNameStartsWithPackageName",
  "milvus/GoReceiverNames",
  "milvus/GoUnsortedImport",
  "milvus/GoUnitSpecificDurationSuffix",
  "milvus/GoTypeParameterInLowerCase",
  "milvus/GoStructInitializationWithoutFieldNames",
  "milvus/GoSnakeCaseUsage",
  "milvus/GoRedundantTrueInForCondition",
  "milvus/GoRedundantElseInIf",
]

const allProjects = [...fileAnalysisProjects, ...globalInspectionProjects, ...golangciLintLocalProjects, ...golangciLintGlobalProjects, ...singleInspectionProjects]

const fileAnalysisCharts: ChartDef[] = [
  { key: "firstCodeAnalysis", label: "File Analysis on Open", measure: "firstCodeAnalysis" },
  { key: "fusHighlighting", label: "Highlighting Duration (FUS)", measure: "fus_daemon_finished_full_duration_since_started_ms" },
]

// Which highlighting pass moved the file analysis time, and whether the findings changed with it: a faster run
// that also reports fewer warnings lost an inspection rather than got faster.
const fileAnalysisBreakdown: BreakdownChart[] = [
  { key: "passes", label: "Highlighting passes", measures: ["LocalInspectionsPass", "GeneralHighlightingPass", "LineMarkersPass", "InlayHintsPass", "CodeVisionPass"] },
  { key: "findings", label: "Findings", measures: ["localInspections#Errors", "localInspections#Warnings", "localInspections#Weak Warnings"] },
]

const globalInspectionCharts: ChartDef[] = [
  {
    key: "globalInspections",
    label: "Global Inspections",
    measure: "globalInspections",
    description: "Batch Inspect Code over the whole project — the offline run, not on-the-fly highlighting.",
  },
]

const golangciLintCharts: ChartDef[] = [
  {
    key: "localInspectionsGolangci",
    label: "Local Inspections (golangci-lint)",
    measure: "localInspections",
    description: "On-the-fly daemon analysis time; here it includes the golangci-lint external linter.",
  },
  {
    key: "globalInspectionsGolangci",
    label: "Global Inspections (golangci-lint)",
    measure: "globalInspections",
    description: "Batch Inspect Code with golangci-lint external linter over the whole project.",
  },
]

const singleInspectionCharts: ChartDef[] = [
  {
    key: "singleInspectionsCodeStyle",
    label: "Single Inspections: Code Style",
    measure: "globalInspections",
    description: "Batch time of each individual Go inspection run in isolation.",
  },
]

const allGroups: GroupDef[] = [
  {
    value: "fileAnalysis",
    title: "File Analysis on Open",
    prefix: "File Analysis",
    projects: fileAnalysisProjects,
    charts: fileAnalysisCharts,
    jvmCharts: jvmCostCharts("firstCodeAnalysis", "File Analysis on Open", "firstCodeAnalysis"),
  },
  {
    value: "globalInspections",
    title: "Global Inspections",
    prefix: "Global",
    projects: globalInspectionProjects,
    charts: globalInspectionCharts,
    jvmCharts: jvmCostCharts("globalInspections", "Global Inspections", "globalInspections"),
  },
  {
    value: "golangciLint",
    title: "golangci-lint",
    prefix: "golangci-lint",
    projects: [...golangciLintLocalProjects, ...golangciLintGlobalProjects],
    charts: golangciLintCharts,
    jvmCharts: [
      ...jvmCostCharts("localInspectionsGolangci", "Local Inspections (golangci-lint)", "localInspections"),
      ...jvmCostCharts("globalInspectionsGolangci", "Global Inspections (golangci-lint)", "globalInspections"),
    ],
  },
  { value: "singleInspections", title: "Single Inspections: Code Style", prefix: "Code Style", projects: singleInspectionProjects, charts: singleInspectionCharts },
]
</script>
