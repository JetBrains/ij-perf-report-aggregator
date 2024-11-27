<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_debugger_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Launch Debug"
        measure="debugRunConfiguration"
        :projects="['river/debug', 'trufflehog/debug']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Step Into"
        measure="debugStep_into"
        :projects="['river/debug', 'trufflehog/debug']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Step Out"
        measure="debugStep_into"
        :projects="['river/debug', 'trufflehog/debug']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Step Over"
        measure="debugStep_into"
        :projects="['river/debug', 'trufflehog/debug']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
