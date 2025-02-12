<template>
  <DashboardPage
    db-name="mlEvaluation"
    table="report"
    persistent-id="chatCodeGenerationDashboard"
    :initial-machine="null"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        measure="Sessions"
        :projects="['chat-code-generation_java_']"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="Sessions"
      />
      <GroupProjectsChart
        measure="APIRecall"
        :projects="['chat-code-generation_java_']"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="APIRecall"
      />
      <GroupProjectsChart
        measure="NotEmptyResultsRatio"
        :projects="['chat-code-generation_java_']"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="NotEmptyResultsRatio"
      />
      <GroupProjectsChart
        measure="MeanLatency"
        :projects="['chat-code-generation_java_']"
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

const measures = ["Sessions", "APIRecall", "NotEmptyResultsRatio", "MeanLatency"]
const chartsDeclaration: ChartDefinition[] = measures.map((measure) => {
  return {
    labels: ["Java"],
    measures: [measure],
    projects: ["chat-code-generation_java_"],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
