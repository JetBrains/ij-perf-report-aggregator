<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 c5.xlarge (4 vCPU, 8 GB)"
    persistent-id="test-generation"
    :with-installer="false"
  >
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationPassedProjects" />
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationCompilationProjects" />
  </DashboardPage>
</template>
<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import ApiTestGenerationDashboardGroupCharts from "./components/ApiTestsDashboardGroupCharts.vue"
import { testGenerationCompilationProjects, testGenerationPassedProjects } from "./api-tests-projects"
</script>
