<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <KotlinDashboardGroupCharts :definitions="scriptCompletionCharts" />
    <KotlinDashboardGroupCharts :definitions="codeAnalysisScriptCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import KotlinDashboardGroupCharts from "../KotlinDashboardGroupCharts.vue"
import { codeAnalysisScriptCharts, scriptCompletionCharts } from "../projects"
</script>
