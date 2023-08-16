<template>
  <DashboardPage
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <DashboardGroupCharts :definitions="localInspectionsCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import DashboardGroupCharts from "./DashboardGroupCharts.vue"
import { localInspectionsCharts } from "./projects"
</script>
