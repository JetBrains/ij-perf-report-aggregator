<template>
  <DashboardPage
    db-name="mlEvaluation"
    table="report"
    persistent-id="aiaDashboard"
    :initial-machine="null"
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
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { getAllProjects } from "./aia"
</script>
