<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <K1K2DashboardGroupCharts :definitions="scriptCompletionCharts" />
    <K1K2DashboardGroupCharts :definitions="highlightingScriptCharts" />
    <K1K2DashboardGroupCharts :definitions="codeAnalysisScriptCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import K1K2DashboardGroupCharts from "../K1K2DashboardGroupCharts.vue"
import { codeAnalysisScriptCharts, highlightingScriptCharts, scriptCompletionCharts } from "../projects"
</script>
