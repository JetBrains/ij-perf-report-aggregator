<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <DashboardGroupCharts :definitions="codeAnalysisCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import DashboardGroupCharts from "../DashboardGroupCharts.vue"
import { codeAnalysisCharts } from "../projects"
</script>
