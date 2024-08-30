<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    :initial-machine="null"
    persistent-id="test-generation"
    :with-installer="false"
  >
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationPsiValidProjects" />
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationCompilationProjects" />
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationPassedProjects" />
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationCreditAvgProjects" />
    <ApiTestGenerationDashboardGroupCharts :definitions="testGenerationQuotaLeftProjects" />
  </DashboardPage>
</template>
<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import ApiTestGenerationDashboardGroupCharts from "./components/ApiTestGenerationDashboardGroupCharts.vue"
import {
  testGenerationCompilationProjects,
  testGenerationCreditAvgProjects,
  testGenerationPassedProjects,
  testGenerationPsiValidProjects,
  testGenerationQuotaLeftProjects,
} from "./api-tests-projects"
</script>
