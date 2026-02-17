<template>
  <DashboardPage
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <KotlinDashboardGroupCharts :definitions="codeAnalysisCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import KotlinDashboardGroupCharts from "./KotlinDashboardGroupCharts.vue"
import { codeAnalysisCharts } from "./projects"
</script>
