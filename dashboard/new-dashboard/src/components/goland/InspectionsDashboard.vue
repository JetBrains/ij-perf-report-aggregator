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
  {
    key: "firstCodeAnalysis",
    label: "File Analysis on Open",
    measure: "firstCodeAnalysis",
    valueUnit: "ms",
    description: "Time to highlight a file the first time it opens (cold daemon pass).",
  },
]

const globalInspectionCharts: ChartDef[] = [
  {
    key: "globalInspections",
    label: "Global Inspections",
    measure: "globalInspections",
    valueUnit: "ms",
    description: "Batch Inspect Code over the whole project — the offline run, not on-the-fly highlighting.",
  },
]

const golangciLintCharts: ChartDef[] = [
  {
    key: "localInspectionsGolangci",
    label: "Local Inspections (golangci-lint)",
    measure: "localInspections",
    valueUnit: "ms",
    description: "On-the-fly daemon analysis time; here it includes the golangci-lint external linter.",
  },
  {
    key: "globalInspectionsGolangci",
    label: "Global Inspections (golangci-lint)",
    measure: "globalInspections",
    valueUnit: "ms",
    description: "Batch Inspect Code with golangci-lint external linter over the whole project.",
  },
]

const singleInspectionCharts: ChartDef[] = [
  {
    key: "singleInspectionsCodeStyle",
    label: "Single Inspections: Code Style",
    measure: "globalInspections",
    valueUnit: "ms",
    description: "Batch time of each individual Go inspection run in isolation.",
  },
]

const allGroups: GroupDef[] = [
  { value: "fileAnalysis", title: "File Analysis on Open", prefix: "File Analysis", projects: fileAnalysisProjects, charts: fileAnalysisCharts },
  { value: "globalInspections", title: "Global Inspections", prefix: "Global", projects: globalInspectionProjects, charts: globalInspectionCharts },
  { value: "golangciLint", title: "golangci-lint", prefix: "golangci-lint", projects: [...golangciLintLocalProjects, ...golangciLintGlobalProjects], charts: golangciLintCharts },
  { value: "singleInspections", title: "Single Inspections: Code Style", prefix: "Code Style", projects: singleInspectionProjects, charts: singleInspectionCharts },
]
</script>
