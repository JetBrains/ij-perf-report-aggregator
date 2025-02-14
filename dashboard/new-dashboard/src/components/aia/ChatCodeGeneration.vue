<template>
  <DashboardPage
    db-name="mlEvaluation"
    table="report"
    persistent-id="chatCodeGenerationDashboard"
    :initial-machine="null"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        measure="Sessions"
        :projects="getAllProjects('chat-code-generation')"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="Sessions"
      />
      <GroupProjectsChart
        measure="APIRecall"
        :projects="getAllProjects('chat-code-generation')"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="APIRecall"
      />
      <GroupProjectsChart
        measure="NotEmptyResultsRatio"
        :projects="getAllProjects('chat-code-generation')"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="NotEmptyResultsRatio"
      />
      <GroupProjectsChart
        measure="MeanLatency"
        :projects="getAllProjects('chat-code-generation')"
        :machines="['Linux EC2 c5.xlarge (4 vCPU, 8 GB)']"
        label="MeanLatency"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { getAllProjects } from "./aia"
</script>
