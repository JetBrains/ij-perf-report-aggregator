<template>
  <DashboardPage
    db-name="mlEvaluation"
    table="report"
    persistent-id="aiaDashboard"
    :initial-machine="null"
    :charts="charts"
    :with-installer="false"
    :branch="null"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :machines="['Linux EC2 C5ad.xlarge (4 vCPU AMD EPYC 7002, 8 GB)']"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { aiaLanguages, aiaModels } from "./aia"

const chartsDeclaration: ChartDefinition[] = aiaLanguages.map((project) => {
  return {
    labels: [project],
    measures: ["SyntaxErrorsSessionRatio"],
    projects: aiaModels.map((model) => "code-generation_" + project + "_" + model),
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
