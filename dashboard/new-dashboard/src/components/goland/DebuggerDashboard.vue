<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_debugger_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Launch Debugger"
        measure="debugRunConfiguration"
        :projects="['river/debug']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Step Into"
        measure="debugStep_into"
        :projects="['river/debug']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
