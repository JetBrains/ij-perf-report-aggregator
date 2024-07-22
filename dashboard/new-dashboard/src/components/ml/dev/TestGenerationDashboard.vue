<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 c5.xlarge (4 vCPU, 8 GB)"
    persistent-id="test-generation"
    :with-installer="false"
  >
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationPsiValidProjects" />
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationCompilationProjects" />
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationPassedProjects" />
  </DashboardPage>
</template>
<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import ApiTestGenerationDashboardGroupCharts from "./components/ApiTestGenerationDashboardGroupCharts.vue"
import { testGenerationCompilationProjects, testGenerationPassedProjects, testGenerationPsiValidProjects } from "./api-tests-projects"
</script>
