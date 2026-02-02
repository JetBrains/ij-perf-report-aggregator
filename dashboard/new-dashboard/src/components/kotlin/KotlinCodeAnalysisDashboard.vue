<template>
  <DashboardPage
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <K1K2DashboardGroupCharts :definitions="codeAnalysisCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import K1K2DashboardGroupCharts from "./K1K2DashboardGroupCharts.vue"
import { codeAnalysisCharts } from "./projects"
</script>
