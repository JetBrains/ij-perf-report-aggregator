<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_debugger_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Run Debug Configuration"
        measure="debugRunConfiguration"
        :projects="['runDebugConfiguration/debug']"
      />
      <GroupProjectsChart
        label="Debug Step Into"
        measure="debugStep_into"
        :projects="['runDebugConfiguration/debug']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
