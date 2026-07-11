<template>
  <DashboardPage
    :charts="charts"
    :with-installer="false"
    db-name="perfintDev"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="goland_product_dashboard"
    table="goland"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :description="chart.definition.description"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>

    <AdditionalMetrics :projects="allProjects" />
  </DashboardPage>
</template>

<script lang="ts" setup>
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import AdditionalMetrics from "./AdditionalMetrics.vue"

const allProjects = [
  "cockroach/indexing",
  "delve/indexing",
  "mattermost/indexing",
  "kubernetes/indexing",
  "flux/indexing",
  "istio/indexing",
  "istio/localInspection/adsc.go",
  "minotaur/localInspection/server.go",
  "caddy/completion/variable",
  "caddy/completion/type",
  "caddy/completion/return",
  "caddy/completion/interface",
  "caddy/completion/import",
  "localAi/go-to-all-with-warmup/version.go/typingLetterByLetter",
  "act/typing",
  "delve/inspection",
  "kubernetes/inspection",
]

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["cockroach/indexing", "delve/indexing", "mattermost/indexing", "kubernetes/indexing", "flux/indexing", "istio/indexing"],
  },
  {
    labels: ["FirstCodeAnalysis"],
    measures: ["firstCodeAnalysis"],
    projects: ["istio/localInspection/adsc.go", "minotaur/localInspection/server.go"],
  },
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: ["caddy/completion/variable", "caddy/completion/type", "caddy/completion/return", "caddy/completion/interface", "caddy/completion/import"],
  },
  {
    labels: ["SearchEverywhere"],
    measures: ["searchEverywhere"],
    projects: ["localAi/go-to-all-with-warmup/version.go/typingLetterByLetter"],
  },
  {
    labels: ["TypingCodeAnalysis"],
    measures: ["typingCodeAnalyzing"],
    projects: ["act/typing"],
  },
  {
    labels: ["UILagging"],
    measures: ["ui.lagging#max_value"],
    projects: ["act/typing"],
  },
  {
    labels: ["Inspections"],
    measures: ["globalInspections"],
    projects: ["delve/inspection", "kubernetes/inspection"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
