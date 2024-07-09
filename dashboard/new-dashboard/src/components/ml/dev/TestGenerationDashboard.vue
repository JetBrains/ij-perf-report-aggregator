<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 c5.xlarge (4 vCPU, 8 GB)"
    persistent-id="test-generation"
    :with-installer="false"
  >
    <ApiTestsDashboardGroupCharts :definitions="testGenerationPassedProjects" />
    <ApiTestsDashboardGroupCharts :definitions="testGenerationCompilationProjects" />
  </DashboardPage>
</template>
<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import ApiTestsDashboardGroupCharts from "./components/ApiTestsDashboardGroupCharts.vue"
import { testGenerationCompilationProjects, testGenerationPassedProjects } from "./api-tests-projects"
</script>
