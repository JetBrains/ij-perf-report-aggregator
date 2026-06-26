<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    :with-installer="false"
    persistent-id="goland_code_analyzes_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section
      v-for="chart in charts"
      :key="chart.key"
    >
      <GroupProjectsChart
        :label="chart.label"
        :measure="chart.measure"
        :projects="chart.projects"
        :better-direction="chart.betterDirection"
        :description="chart.description"
      />
    </section>
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

const charts: ChartDef[] = [
  {
    key: "firstCodeAnalysis",
    label: "File Analysis on Open",
    measure: "firstCodeAnalysis",
    projects: [
      "istio/localInspection/adsc.go",
      "minotaur/localInspection/server.go",
      "caddy/localInspection/replacer.go",
      "cockroach/localInspection/tochar.go",
      "flux/localInspection/sourcesecret.go",
      "kubernetes/localInspection/client_test.go",
    ],
    description: "Time to highlight a file the first time it opens (cold daemon pass).",
  },
  {
    key: "globalInspections",
    label: "Global Inspections",
    measure: "globalInspections",
    projects: [
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
    ],
    description: "Batch Inspect Code over the whole project — the offline run, not on-the-fly highlighting.",
  },
  {
    key: "localInspectionsGolangci",
    label: "Local Inspections (golangci-lint)",
    measure: "localInspections",
    projects: ["volcano-golinter-local-without-linter/localInspection/scheduler.go", "volcano-golinter-local-with-linter/localInspection/scheduler.go"],
    description: "On-the-fly daemon analysis time; here it includes the golangci-lint external linter.",
  },
  {
    key: "globalInspectionsGolangci",
    label: "Global Inspections (golangci-lint)",
    measure: "globalInspections",
    projects: ["volcano-golinter-global-with-linter/inspection", "volcano-golinter-global-without-linter/inspection"],
  },
  {
    key: "singleInspectionsCodeStyle",
    label: "Single Inspections: Code Style",
    measure: "globalInspections",
    projects: [
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
    ],
    description: "Batch time of each individual Go inspection run in isolation.",
  },
]
</script>
