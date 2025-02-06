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
        measure="Precision"
        :projects="getAllProjects('code-generation')"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="Precision"
      />
      <GroupProjectsChart
        measure="FileValidationSuccess"
        :projects="getAllProjects('code-generation')"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="FileValidationSuccess"
      />
      <GroupProjectsChart
        measure="MeanContextSize"
        :projects="getAllProjects('code-generation')"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="MeanContextSize"
      />
      <GroupProjectsChart
        measure="MeanLatency"
        :projects="getAllProjects('code-generation')"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="MeanLatency"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { aiaLanguages, aiaModels, getAllProjects } from "./aia"

const chartsDeclaration: ChartDefinition[] = aiaLanguages.map((project) => {
  return {
    labels: [project],
    measures: ["Precision"],
    projects: aiaModels.map((model) => "code-generation_" + project + "_" + model),
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
